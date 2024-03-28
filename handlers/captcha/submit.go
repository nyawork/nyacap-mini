package captcha

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/wenlng/go-captcha/captcha"
	"net/http"
	"nya-captcha/config"
	"nya-captcha/consts"
	"nya-captcha/global"
	"nya-captcha/security"
	"nya-captcha/types"
	"time"
)

type CaptchaSubmitRequest struct {
	Key  string `json:"k"`
	Dots []struct {
		X int64 `json:"x"`
		Y int64 `json:"y"`
	} `json:"d"`
}

type CaptchaSubmitResponse struct {
	Success bool `json:"s"`
}

func Submit(ctx *gin.Context) {
	isCoolingdown, err := security.CheckIPCD(ctx.ClientIP(), consts.IPCD_POOL_SUBMIT)
	if err != nil {
		global.Logger.Errorf("检查 IP 提交冷却状态失败: %v", err)
		ctx.Status(http.StatusInternalServerError)
	} else if isCoolingdown {
		ctx.Status(http.StatusTooManyRequests)
		return
	} else {
		security.CooldownIP(ctx.ClientIP(), consts.IPCD_POOL_SUBMIT, config.Config.Security.CaptchaSubmitCooldown)
	}

	var req CaptchaSubmitRequest
	err = ctx.BindJSON(&req)
	if err != nil {
		global.Logger.Errorf("请求数据格式化失败: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	// 使用指定的 key 请求 redis
	captchaPendingStateByte, err := global.Redis.GetDel(
		context.Background(),
		fmt.Sprintf(consts.REDIS_KEY_REQUEST_SESSION, req.Key),
	).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			ctx.Status(http.StatusNotFound)
		} else {
			global.Logger.Errorf("验证码信息拉取失败: %v", err)
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}

	// 格式化
	var captchaPendingState types.CaptchaPending
	err = json.Unmarshal([]byte(captchaPendingStateByte), &captchaPendingState)
	if err != nil {
		global.Logger.Errorf("无法解码存储的验证码信息: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// 比较请求来源
	if captchaPendingState.Origin != ctx.GetHeader("Origin") ||
		//captchaPendingState.IP != ctx.ClientIP() || // 多出口时候 IP 确实有可能变化，暂时先不根据这个屏蔽
		captchaPendingState.UserAgent != ctx.Request.UserAgent() {
		// 请求来源变化了
		ctx.Status(http.StatusForbidden)
		security.CooldownIP(ctx.ClientIP(), consts.IPCD_POOL_BAN, config.Config.Security.IPBanPeriod)
		return
	}

	// 初判长度
	if len(captchaPendingState.Dots) != len(req.Dots) {
		// 长度不一致
		ctx.JSON(http.StatusOK, CaptchaSubmitResponse{
			Success: false,
		})
		return
	}

	// 检测每个点的位置是否对应的上
	for index, dot := range captchaPendingState.Dots {
		if !captcha.CheckPointDistWithPadding(
			req.Dots[index].X, req.Dots[index].Y,
			int64(dot.Dx), int64(dot.Dy),
			int64(dot.Width), int64(dot.Height),
			1,
		) {
			// 点不对应
			ctx.JSON(http.StatusOK, CaptchaSubmitResponse{
				Success: false,
			})
			return
		}
	}

	// 通过校验，记录结果
	timeStamp := time.Now().Format(time.RFC3339)
	resolvedStateBytes, err := json.Marshal(types.CaptchaResolved{
		Success:     true,
		IP:          &captchaPendingState.IP,
		Origin:      &captchaPendingState.Origin,
		ChallengeTS: &timeStamp,
		Hostname:    &ctx.Request.Host,
	})
	if err != nil {
		global.Logger.Errorf("无法格式化验证码信息: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	err = global.Redis.Set(
		context.Background(),
		fmt.Sprintf(consts.REDIS_KEY_RESOLVED_SESSION, req.Key),
		resolvedStateBytes,
		config.Config.Captcha.SubmitValidFor,
	).Err()
	if err != nil {
		global.Logger.Errorf("无法保存验证码会话: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, CaptchaSubmitResponse{
		Success: true,
	})

}
