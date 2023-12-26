package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nya-captcha/consts"
	"nya-captcha/global"
	"nya-captcha/security"
)

func IPBan() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isBanned, err := security.CheckIPCD(ctx.ClientIP(), consts.IPCD_POOL_BAN)
		if err != nil {
			global.Logger.Errorf("检查 IP 是否被 ban 状态失败: %v", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		} else if isBanned {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
		}
	}
}
