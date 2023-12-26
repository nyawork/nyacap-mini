package routers

import (
	"github.com/gin-gonic/gin"
	"nya-captcha/handlers/public"
)

func Public(rg *gin.RouterGroup) {
	rg.GET("/healthcheck", public.HealthCheck)
}
