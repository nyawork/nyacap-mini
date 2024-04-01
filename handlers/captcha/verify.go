package captcha

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"nyacap-mini/config"
	"nyacap-mini/consts"
	g "nyacap-mini/global"
	"nyacap-mini/security"
	"nyacap-mini/types"
	"nyacap-mini/utils"
)

type CaptchaVerifyRequest struct {
	SiteSecret string `form:"secret" binding:"required"`   // 站点密钥
	Key        string `form:"response" binding:"required"` // 会话 key
}

func Verify(c echo.Context) error {
	var req CaptchaVerifyRequest
	err := c.Bind(&req)
	if err != nil || req.SiteSecret == "" || req.Key == "" {
		g.Logger.Error("请求数据格式化失败", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	// 寻找一个私钥对得上的 site
	siteInfo, found := findSiteBySecret(req.SiteSecret)
	if !found {
		// 查无此站
		security.CooldownIP(c.RealIP(), consts.IPCD_POOL_BAN, config.Config.Security.IPBanPeriod)
		return echo.NewHTTPError(http.StatusForbidden)
	}

	// 使用指定的 key 请求 redis
	captchaResolvedStateByte, err := g.Redis.GetDel(
		context.Background(),
		fmt.Sprintf(consts.REDIS_KEY_RESOLVED_SESSION, req.Key),
	).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return c.JSON(http.StatusOK, types.CaptchaResolved{
				Success:    false,
				ErrorCodes: []string{"timeout-or-duplicate"},
			})
		} else {
			g.Logger.Error("验证码结果拉取失败", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	// 格式化
	var captchaResolvedState types.CaptchaResolved
	err = json.Unmarshal([]byte(captchaResolvedStateByte), &captchaResolvedState)
	if err != nil {
		g.Logger.Error("无法解码存储的验证码结果", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// 比较 site
	if !utils.SliceExist(siteInfo.AllowedOrigins, *captchaResolvedState.Origin) {
		// ban IP
		security.CooldownIP(*captchaResolvedState.IP, consts.IPCD_POOL_BAN, config.Config.Security.IPBanPeriod)
		return c.JSON(http.StatusOK, types.CaptchaResolved{
			Success:    false,
			ErrorCodes: []string{"bad-request"},
		})
	}

	// 验证完成，返回结果
	return c.JSON(http.StatusOK, captchaResolvedState)

}

func findSiteBySecret(siteSecret string) (*types.SiteInfo, bool) {
	for _, siteInfo := range config.Config.Sites {
		if siteSecret == siteInfo.SiteSecret {
			return &siteInfo, true
		}
	}
	return nil, false
}
