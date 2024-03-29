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
	"nya-captcha/utils"
	"time"
)

type RequestResponse struct {
	Key             string `json:"k"`
	BigImageBase64  string `json:"b"`
	ThumbnailBase64 string `json:"t"`
	ExpiresAt       int64  `json:"e"`
}

func Request(ctx *gin.Context) {
	ip := ctx.ClientIP()
	isCoolingdown, err := security.CheckIPCD(ip, consts.IPCD_POOL_REQUEST)
	if err != nil {
		global.Logger.Errorf("检查 IP 请求冷却状态失败: %v", err)
		ctx.Status(http.StatusInternalServerError)
	} else if isCoolingdown {
		global.Logger.Debugf("IP (%s) 处于 %s 冷却池中", ip, consts.IPCD_POOL_REQUEST)
		ctx.Status(http.StatusTooManyRequests)
		return
	} else {
		global.Logger.Debugf("IP (%s) 没有问题，继续请求", ip)
		security.CooldownIP(ip, consts.IPCD_POOL_REQUEST, config.Config.Security.CaptchaRequestCooldown)
	}

	siteKey := ctx.Param("site_key")

	// 验证站点公钥和 Origin 是否匹配
	siteInfo := findSiteBySiteKey(siteKey)
	if siteInfo == nil || !utils.SliceExist(siteInfo.AllowedOrigins, ctx.GetHeader("Origin")) {
		// 站点公钥不匹配或 Origin 未被允许， ban IP
		ctx.Status(http.StatusForbidden)
		security.CooldownIP(ip, consts.IPCD_POOL_BAN, config.Config.Security.IPBanPeriod)
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
		Origin:    ctx.GetHeader("Origin"),
		IP:        ip,
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

	ctx.JSON(http.StatusOK, &RequestResponse{
		Key:             key,
		BigImageBase64:  b64,
		ThumbnailBase64: tb64,
		ExpiresAt:       expireAt.Unix(),
	})

}

func findSiteBySiteKey(siteKey string) *types.SiteInfo {
	for _, siteInfo := range config.Config.Sites {
		if siteKey == siteInfo.SiteKey {
			return &siteInfo
		}
	}
	return nil
}
