package captcha

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"nya-captcha/config"
	"nya-captcha/consts"
	"nya-captcha/global"
	"nya-captcha/security"
	"nya-captcha/types"
	"nya-captcha/utils"
)

type CaptchaVerifyRequest struct {
	SiteSecret string `form:"secret" binding:"required"`   // 站点密钥
	Key        string `form:"response" binding:"required"` // 会话 key
}

func Verify(ctx *gin.Context) {
	var req CaptchaVerifyRequest
	err := ctx.Bind(&req)
	if err != nil || req.SiteSecret == "" || req.Key == "" {
		global.Logger.Errorf("请求数据格式化失败: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	// 寻找一个私钥对得上的 site
	siteInfo, found := findSiteBySecret(req.SiteSecret)
	if !found {
		// 查无此站
		ctx.Status(http.StatusForbidden)
		security.CooldownIP(ctx.ClientIP(), consts.IPCD_POOL_BAN, config.Config.Security.IPBanPeriod)
		return
	}

	// 使用指定的 key 请求 redis
	captchaResolvedStateByte, err := global.Redis.GetDel(
		context.Background(),
		fmt.Sprintf(consts.REDIS_KEY_RESOLVED_SESSION, req.Key),
	).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			ctx.JSON(http.StatusOK, types.CaptchaResolved{
				Success:    false,
				ErrorCodes: []string{"timeout-or-duplicate"},
			})
		} else {
			global.Logger.Errorf("验证码结果拉取失败: %v", err)
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}

	// 格式化
	var captchaResolvedState types.CaptchaResolved
	err = json.Unmarshal([]byte(captchaResolvedStateByte), &captchaResolvedState)
	if err != nil {
		global.Logger.Errorf("无法解码存储的验证码结果: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// 比较 site
	if !utils.SliceExist(siteInfo.AllowedOrigins, *captchaResolvedState.Origin) {
		// ban IP
		security.CooldownIP(*captchaResolvedState.IP, consts.IPCD_POOL_BAN, config.Config.Security.IPBanPeriod)
		ctx.JSON(http.StatusOK, types.CaptchaResolved{
			Success:    false,
			ErrorCodes: []string{"bad-request"},
		})
	}

	// 验证完成，返回结果
	ctx.JSON(http.StatusOK, captchaResolvedState)

}

func findSiteBySecret(siteSecret string) (*types.SiteInfo, bool) {
	for _, siteInfo := range config.Config.Sites {
		if siteSecret == siteInfo.SiteSecret {
			return &siteInfo, true
		}
	}
	return nil, false
}
