package captcha

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/wenlng/go-captcha/captcha"
	"go.uber.org/zap"
	"net/http"
	"nyacap-mini/config"
	"nyacap-mini/consts"
	g "nyacap-mini/global"
	"nyacap-mini/security"
	"nyacap-mini/types"
	"time"
)

type CaptchaSubmitRequest struct {
	Key  string `json:"k" form:"k"`
	Dots []struct {
		X int64 `json:"x" form:"x"`
		Y int64 `json:"y" form:"y"`
	} `json:"d" form:"d"`
}

type CaptchaSubmitResponse struct {
	Success bool `json:"s"`
}

func Submit(c echo.Context) error {
	ip := c.RealIP()
	isCoolingdown, err := security.CheckIPCD(ip, consts.IPCD_POOL_SUBMIT)
	if err != nil {
		g.Logger.Error("检查 IP 请求冷却状态失败", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	} else if isCoolingdown {
		g.Logger.Debug("IP 处于冷却池中", zap.String("ip", ip), zap.String("pool", consts.IPCD_POOL_REQUEST))
		return echo.NewHTTPError(http.StatusInternalServerError)
	} else {
		g.Logger.Debug("IP 没有问题，继续请求", zap.String("ip", ip))
		security.CooldownIP(ip, consts.IPCD_POOL_SUBMIT, config.Config.Security.CaptchaSubmitCooldown)
	}

	var req CaptchaSubmitRequest
	err = c.Bind(&req)
	if err != nil {
		g.Logger.Error("请求数据格式化失败", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	// 使用指定的 key 请求 redis
	captchaPendingStateByte, err := g.Redis.GetDel(
		context.Background(),
		fmt.Sprintf(consts.REDIS_KEY_REQUEST_SESSION, req.Key),
	).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return echo.NewHTTPError(http.StatusNotFound)
		} else {
			g.Logger.Error("验证码信息拉取失败", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	// 格式化
	var captchaPendingState types.CaptchaPending
	err = json.Unmarshal([]byte(captchaPendingStateByte), &captchaPendingState)
	if err != nil {
		g.Logger.Error("无法解码存储的验证码信息", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// 比较请求来源
	if captchaPendingState.Origin != c.Request().Header.Get("Origin") ||
		//captchaPendingState.IP != ip || // 多出口时候 IP 确实有可能变化，暂时先不根据这个屏蔽
		captchaPendingState.UserAgent != c.Request().UserAgent() {
		// 请求来源变化了
		security.CooldownIP(ip, consts.IPCD_POOL_BAN, config.Config.Security.IPBanPeriod)
		return echo.NewHTTPError(http.StatusForbidden)
	}

	// 初判长度
	if len(captchaPendingState.Dots) != len(req.Dots) {
		// 长度不一致
		g.Logger.Debug("验证结果长度不匹配", zap.Int("期望", len(captchaPendingState.Dots)), zap.Int("得到", len(req.Dots)))
		return c.JSON(http.StatusOK, CaptchaSubmitResponse{
			Success: false,
		})
	}

	// 检测每个点的位置是否对应的上
	for index, dot := range captchaPendingState.Dots {
		if !captcha.CheckPointDistWithPadding(
			req.Dots[index].X, req.Dots[index].Y,
			int64(dot.Dx), int64(dot.Dy),
			int64(dot.Width), int64(dot.Height),
			config.Config.Captcha.Padding,
		) {
			// 点不对应
			g.Logger.Debug("验证结果位置不对应")
			return c.JSON(http.StatusOK, CaptchaSubmitResponse{
				Success: false,
			})
		}
	}

	// 通过校验，记录结果
	timeStamp := time.Now().Format(time.RFC3339)
	resolvedStateBytes, err := json.Marshal(types.CaptchaResolved{
		Success:     true,
		IP:          &captchaPendingState.IP,
		Origin:      &captchaPendingState.Origin,
		ChallengeTS: &timeStamp,
		Hostname:    &c.Request().Host,
	})
	if err != nil {
		g.Logger.Error("无法格式化验证码信息", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = g.Redis.Set(
		context.Background(),
		fmt.Sprintf(consts.REDIS_KEY_RESOLVED_SESSION, req.Key),
		resolvedStateBytes,
		config.Config.Captcha.SubmitValidFor,
	).Err()
	if err != nil {
		g.Logger.Error("无法保存验证码会话", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, CaptchaSubmitResponse{
		Success: true,
	})

}
