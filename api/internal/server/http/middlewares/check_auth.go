package middlewares

import (
	"github.com/gin-gonic/gin"
	"microservice/internal/adapter/locale"
	"net/http"
)

func CheckAuth(l locale.ILocale) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error

		_ = ctx.GetHeader("jwt-token")

		//todo: validate the JWT token here

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{
				"error": l.Get("unauthorized"),
				"data":  nil,
			})
			return
		}

		ctx.Next()
	}
}
