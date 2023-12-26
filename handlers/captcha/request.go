package captcha

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"nya-captcha/config"
	"nya-captcha/consts"
	"nya-captcha/global"
	"nya-captcha/security"
	"nya-captcha/types"
	"time"
)

type CaptchaRequestRequest struct {
	SiteKey string `json:"site_key"`
}

type CaptchaRequestResponse struct {
	Key             string `json:"k"`
	BigImageBase64  string `json:"b"`
	ThumbnailBase64 string `json:"t"`
	ExpiresAt       int64  `json:"e"`
}

func Request(ctx *gin.Context) {
	isCoolingdown, err := security.CheckIPCD(ctx.ClientIP(), consts.IPCD_POOL_REQUEST)
	if err != nil {
		global.Logger.Errorf("检查 IP 请求冷却状态失败: %v", err)
		ctx.Status(http.StatusInternalServerError)
	} else if isCoolingdown {
		ctx.Status(http.StatusTooManyRequests)
		return
	} else {
		security.CooldownIP(ctx.ClientIP(), consts.IPCD_POOL_REQUEST, config.Config.Security.CaptchaSubmitCooldown)
	}

	var req CaptchaRequestRequest
	err = ctx.BindJSON(&req)
	if err != nil {
		global.Logger.Errorf("请求数据格式化失败: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	// 验证站点公钥和 Origin 是否匹配
	siteInfo, ok := config.Config.Sites[ctx.GetHeader("Origin")]
	if !ok || siteInfo.SiteKey != req.SiteKey {
		// 站点公钥不匹配， ban IP
		ctx.Status(http.StatusForbidden)
		security.CooldownIP(ctx.ClientIP(), consts.IPCD_POOL_BAN, config.Config.Security.IPBanPeriod)
		return
	}

	// 生成验证码
	dots, b64, tb64, key, err := global.Captcha.Generate()
	if err != nil {
		global.Logger.Errorf("验证码创建失败: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	pendingStateBytes, err := json.Marshal(types.CaptchaPending{
		Site:      ctx.GetHeader("Origin"),
		IP:        ctx.ClientIP(),
		UserAgent: ctx.Request.UserAgent(),
		Dots:      dots,
	})
	if err != nil {
		global.Logger.Errorf("无法格式化验证码信息: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	expireAt := time.Now().Add(config.Config.Captcha.PendingValidFor)
	err = global.Redis.Set(
		context.Background(),
		fmt.Sprintf(consts.REDIS_KEY_REQUEST_SESSION, key),
		pendingStateBytes,
		config.Config.Captcha.PendingValidFor,
	).Err()
	if err != nil {
		global.Logger.Errorf("无法保存验证码会话: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, &CaptchaRequestResponse{
		Key:             key,
		BigImageBase64:  b64,
		ThumbnailBase64: tb64,
		ExpiresAt:       expireAt.Unix(),
	})

}
