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
)

type CaptchaVerifyRequest struct {
	// 模仿 recaptcha 的格式，方便后端迁移使用
	SiteSecret string  `json:"secret"`   // 站点密钥
	Key        string  `json:"response"` // 会话 key
	IP         *string `json:"remoteip,omitempty"`

	// 扩展
	Site *string `json:"site"` // 站点 Origin
}

func Verify(ctx *gin.Context) {
	var req CaptchaVerifyRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		global.Logger.Errorf("请求数据格式化失败: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	// 检查是否有 Site
	var site string
	if req.Site != nil {
		// 尝试定位到对应的 site
		siteInfo, ok := config.Config.Sites[*req.Site]
		if !ok || siteInfo.SiteSecret != req.SiteSecret {
			// 站点密钥不匹配， ban IP
			ctx.Status(http.StatusForbidden)
			security.CooldownIP(ctx.ClientIP(), consts.IPCD_POOL_BAN, config.Config.Security.IPBanPeriod)
			return
		}
		site = *req.Site
	} else {
		// 寻找一个私钥对得上的 site
		var found bool
		site, found = findSiteBySecret(req.SiteSecret)
		if !found {
			// 查无此站
			ctx.Status(http.StatusForbidden)
			security.CooldownIP(ctx.ClientIP(), consts.IPCD_POOL_BAN, config.Config.Security.IPBanPeriod)
			return
		}
	}

	// 使用指定的 key 请求 redis
	captchaResolvedStateByte, err := global.Redis.GetDel(
		context.Background(),
		fmt.Sprintf(consts.REDIS_KEY_RESOLVED_SESSION, req.Key),
	).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			ctx.Status(http.StatusNotFound)
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

	// 比较 IP
	if req.IP != nil {
		if captchaResolvedState.IP != *req.IP {
			ctx.JSON(http.StatusOK, types.CaptchaResolved{
				Success: false,
			})
		}
	}

	// 比较 site
	if captchaResolvedState.Site != site {
		// ban IP
		security.CooldownIP(captchaResolvedState.IP, consts.IPCD_POOL_BAN, config.Config.Security.IPBanPeriod)
		ctx.JSON(http.StatusOK, types.CaptchaResolved{
			Success: false,
		})
	}

	// 验证完成，返回结果
	ctx.JSON(http.StatusOK, captchaResolvedState)

}

func findSiteBySecret(siteSecret string) (string, bool) {
	for site, siteInfo := range config.Config.Sites {
		if siteSecret == siteInfo.SiteSecret {
			return site, true
		}
	}
	return "", false
}
