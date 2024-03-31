package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	g "nya-captcha/global"
)

func LogValues(_ echo.Context, v middleware.RequestLoggerValues) error {
	if v.Error == nil {
		g.Logger.Info("request",
			zap.String("URI", v.URI),
			zap.Int("status", v.Status),
		)
	} else {
		g.Logger.Error("request",
			zap.String("URI", v.URI),
			zap.Int("status", v.Status),
			zap.Error(v.Error),
		)
	}

	return nil
}
