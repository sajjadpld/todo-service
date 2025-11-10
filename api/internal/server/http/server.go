package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"microservice/app"
	"microservice/config"
	"microservice/internal/adapter/locale"
	"microservice/internal/adapter/registry"
	"microservice/internal/server/http/middlewares"
	_ "microservice/pkg/validator"
	"net/http"
	"os"
)

type Server struct {
	l            locale.ILocale
	service      config.Service
	swagger      config.Swagger
	config       config.Http
	repositories *app.Repositories
	handlers     *app.HttpHandlers
	engine       *gin.Engine
	server       *http.Server
}

func New(
	registry registry.IRegistry,
	locale locale.ILocale,
	repositories *app.Repositories,
	handlers *app.HttpHandlers,
) IHttpServer {
	server := new(Server)
	registry.Parse(&server.service)
	registry.Parse(&server.swagger)
	registry.Parse(&server.config)

	if server.service.Debug == false {
		server.swagger.Host = os.Getenv("SWAGGER_HOST")
	}

	server.l = locale
	server.handlers = handlers
	server.repositories = repositories
	server.engine = gin.Default()

	return server
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}

func (s *Server) Init() {
	err := s.engine.SetTrustedProxies([]string{trustedProxy})
	if err != nil {
		log.Panicf("[http] set trusted proxy failed: %s\n", err.Error())
	}

	if s.service.Debug == false {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	} else {
		gin.SetMode(gin.DebugMode)
	}

	s.engine.Use(
		//gin.Logger(),
		gin.Recovery(), middlewares.Cors(),
		gin.CustomRecovery(middlewares.ErrorHandler),
	)
}

func (s *Server) Start() {
	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.config.Host, s.config.Port),
		Handler: s.engine,
	}

	fmt.Printf("\n[http] started on port %s\n", s.config.Port)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panicf("[http] failed to start usecase: %s\n", err.Error())
		}
	}()
}

func (s *Server) Stop(ctx context.Context) {
	if s.service.Debug == true {
		_ = s.server.Close()
		log.Printf("[http] server stopped succesfully")
		return
	}

	shutdownErr := make(chan error)
	go func() { shutdownErr <- s.server.Shutdown(ctx) }()

	select {
	case err := <-shutdownErr:
		if err != nil {
			log.Printf("[http] server shutdown err: %s", err)
			_ = s.server.Close() // force to close usecase
			log.Printf("[http] server closed")
		}

		log.Printf("[http] server stopped succesfully")
	case <-ctx.Done():
		_ = s.server.Close()
		log.Printf("[http] context timeout - server closed")
	}
}
