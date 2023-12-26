package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nya-captcha/config"
	"nya-captcha/utils"
)

func CORS() gin.HandlerFunc {
	var validOrigins []string
	for origin := range config.Config.Sites {
		validOrigins = append(validOrigins, origin)
	}

	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		if origin != "" {
			if config.Config.System.Debug || utils.SliceExist(validOrigins, origin) {
				ctx.Header("Access-Control-Allow-Origin", "*")
				ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
				ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
				ctx.Header("Access-Control-Allow-Credentials", "true")

				if ctx.Request.Method == "OPTIONS" {
					ctx.AbortWithStatus(http.StatusNoContent)
				}
			} else {
				// Otherwise block requests for safety concern
				ctx.AbortWithStatus(http.StatusForbidden)
			}
		}

	}
}
