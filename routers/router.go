package routers

import (
	"github.com/labstack/echo/v4"
	"nyacap-mini/middlewares"
)

func R(e *echo.Echo) {
	// Public
	publicApi := e.Group("")
	Public(publicApi)

	// Captcha
	captchaApi := e.Group("/captcha")
	captchaApi.Use(middlewares.CORS())
	captchaApi.Use(middlewares.IPBan())
	Captcha(captchaApi)
}
