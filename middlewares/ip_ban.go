package middlewares

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"nyacap-mini/consts"
	g "nyacap-mini/global"
	"nyacap-mini/security"
)

func IPBan() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isBanned, err := security.CheckIPCD(c.RealIP(), consts.IPCD_POOL_BAN)
			if err != nil {
				g.Logger.Error("检查 IP 是否被 ban 状态失败", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError)
			} else if isBanned {
				return echo.NewHTTPError(http.StatusTooManyRequests)
			}

			// 继续处理
			return next(c)
		}
	}
}
