package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandler(c *gin.Context, caughtErr any) {
	if err, ok := caughtErr.(error); ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error)
		return
	}

	c.AbortWithStatusJSON(
		http.StatusInternalServerError,
		caughtErr,
	)
	return
}
