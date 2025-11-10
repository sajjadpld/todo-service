package http

import (
	"microservice/internal/server/http/routes"
)

func (s *Server) SetRoutes() {
	service := s.Engine()
	router := service.Group("/")

	// general

	{
		router.GET("handshake", routes.Handshake)
		routes.SwaggerRoute(router, &s.swagger)
	}

	// routes groups

	api := router.Group("api")
	{
		v1 := api.Group("/v1")
		{
			routes.TodoRoutes(v1, s.handlers.TodoHandler, s.l)
			// NOTE: set other routes as above
		}
	}
}
