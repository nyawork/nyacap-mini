package routers

import (
	"github.com/labstack/echo/v4"
	"nyacap-mini/handlers/captcha"
	"nyacap-mini/middlewares"
)

func Captcha(rg *echo.Group) {
	// 来自服务器的验证请求不需要 CORS
	rg.POST("/verify", captcha.Verify)

	// 前端请求需要
	rg.Use(middlewares.CORS())
	rg.GET("/request/:site_key", captcha.Request)
	rg.POST("/submit", captcha.Submit)
}
