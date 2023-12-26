package routers

import (
	"github.com/gin-gonic/gin"
	"nya-captcha/middlewares"
)

func R(e *gin.Engine) {
	e.Use(middlewares.CORS())
	e.Use(middlewares.IPBan())

	// Public
	publicApi := e.Group("")
	Public(publicApi)

	// Captcha
	captchaApi := e.Group("/captcha")
	Captcha(captchaApi)
}
