package http

import (
	"context"
	"github.com/gin-gonic/gin"
)

type (
	IHttpServer interface {
		Init()
		SetRoutes()
		Engine() *gin.Engine
		Start()
		Stop(ctx context.Context)
	}
)
