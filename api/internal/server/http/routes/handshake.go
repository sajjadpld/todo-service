package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type response struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// Handshake godoc
// @Summary Service Handshake
// @Description Checks the Service Availability
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} response
// @SetID GET-handshake
// @Router /handshake [get]
func Handshake(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response{
		Status:    http.StatusText(http.StatusOK),
		Message:   "connection established",
		Timestamp: time.Now().Format(time.DateTime),
	})
	return
}
