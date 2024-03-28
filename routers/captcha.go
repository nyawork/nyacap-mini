package routers

import (
	"github.com/gin-gonic/gin"
	"nya-captcha/handlers/captcha"
)

func Captcha(rg *gin.RouterGroup) {
	rg.GET("/request/:site_key", captcha.Request)
	rg.POST("/submit", captcha.Submit)
	rg.POST("/verify", captcha.Verify)
}
