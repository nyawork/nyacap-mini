package captcha

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"nyacap-mini/config"
	"nyacap-mini/consts"
	g "nyacap-mini/global"
	"nyacap-mini/security"
	"nyacap-mini/types"
	"nyacap-mini/utils"
	"time"
)

type RequestResponse struct {
	Key             string `json:"k"`
	BigImageBase64  string `json:"b"`
	ThumbnailBase64 string `json:"t"`
	ExpiresAt       int64  `json:"e"`
}

func Request(c echo.Context) error {
	ip := c.RealIP()
	isCoolingdown, err := security.CheckIPCD(ip, consts.IPCD_POOL_REQUEST)
	if err != nil {
		g.Logger.Error("检查 IP 请求冷却状态失败", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	} else if isCoolingdown {
		g.Logger.Debug("IP 处于冷却池中", zap.String("ip", ip), zap.String("pool", consts.IPCD_POOL_REQUEST))
		return echo.NewHTTPError(http.StatusTooManyRequests)
	} else {
		g.Logger.Debug("IP 没有问题，继续请求", zap.String("ip", ip))
		security.CooldownIP(ip, consts.IPCD_POOL_REQUEST, config.Config.Security.CaptchaRequestCooldown)
	}

	siteKey := c.Param("site_key")

	// 验证站点公钥和 Origin 是否匹配
	siteInfo := findSiteBySiteKey(siteKey)
	origin := c.Request().Header.Get("Origin")
	if siteInfo == nil || !utils.SliceExist(siteInfo.AllowedOrigins, origin) {
		// 站点公钥不匹配或 Origin 未被允许， ban IP
		security.CooldownIP(ip, consts.IPCD_POOL_BAN, config.Config.Security.IPBanPeriod)
		return echo.NewHTTPError(http.StatusForbidden)
	}

	// 生成验证码
	dots, b64, tb64, key, err := g.Captcha.Generate()
	if err != nil {
		g.Logger.Error("验证码创建失败", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	pendingStateBytes, err := json.Marshal(types.CaptchaPending{
		Origin:    origin,
		IP:        ip,
		UserAgent: c.Request().UserAgent(),
		Dots:      dots,
	})
	if err != nil {
		g.Logger.Error("无法格式化验证码信息", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	expireAt := time.Now().Add(config.Config.Captcha.PendingValidFor)
	err = g.Redis.Set(
		context.Background(),
		fmt.Sprintf(consts.REDIS_KEY_REQUEST_SESSION, key),
		pendingStateBytes,
		config.Config.Captcha.PendingValidFor,
	).Err()
	if err != nil {
		g.Logger.Error("无法保存验证码会话", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, &RequestResponse{
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
