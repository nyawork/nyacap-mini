package routers

import (
	"github.com/labstack/echo/v4"
	"nyacap-mini/handlers/captcha"
)

func Captcha(rg *echo.Group) {
	rg.GET("/request/:site_key", captcha.Request)
	rg.POST("/submit", captcha.Submit)
	rg.POST("/verify", captcha.Verify)
}
