package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
	"microservice/config"
	"microservice/docs" // the custom path of the generated swagger files
	"microservice/internal/server/http/middlewares"
)

func SwaggerRoute(r *gin.RouterGroup, conf *config.Swagger) {
	if conf.Enable == false {
		return
	}

	docs.SwaggerInfo.Title = conf.Title
	docs.SwaggerInfo.Description = conf.Description
	docs.SwaggerInfo.Version = conf.Version
	docs.SwaggerInfo.Schemes = []string{conf.Schemes}
	docs.SwaggerInfo.Host = conf.Host

	r.GET("public/swagger/*any",
		middlewares.SwaggerAuth(conf.Username, conf.Password),
		swagger.WrapHandler(swaggerFiles.Handler, swagger.DefaultModelsExpandDepth(-1)),
	)
}
