package middlewares

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"nya-captcha/config"
	"nya-captcha/utils"
)

func CORS() echo.MiddlewareFunc {
	var validOrigins []string
	for _, siteInfo := range config.Config.Sites {
		for _, origin := range siteInfo.AllowedOrigins {
			if !utils.SliceExist(validOrigins, origin) {
				validOrigins = append(validOrigins, origin)
			}
		}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			origin := c.Request().Header.Get("Origin")
			if origin != "" {
				if config.Config.System.Debug || utils.SliceExist(validOrigins, origin) {
					h := c.Response().Header()
					h.Set("Access-Control-Allow-Origin", "*")
					h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
					h.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
					h.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
					h.Set("Access-Control-Allow-Credentials", "true")

					if c.Request().Method == "OPTIONS" {
						return c.NoContent(http.StatusNoContent)
					}

					// 继续处理
					return next(c)
				}
			}

			// 否则 (无 origin 或不匹配) 出于安全考虑拒绝请求
			return echo.NewHTTPError(http.StatusForbidden)
		}
	}

}
