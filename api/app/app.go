package app

import (
	"microservice/config"
	"microservice/internal/adapter/locale"
	"microservice/internal/adapter/logger"
	"microservice/internal/adapter/orm"
	"microservice/internal/adapter/registry"
)

// App Dependency Injection
type App struct {
	config       *config.Service
	swagger      *config.Swagger
	registry     registry.IRegistry
	logger       logger.ILogger
	locale       locale.ILocale
	database     orm.ISql
	repo         *Repositories
	port         *Ports
	httpHandlers *HttpHandlers
}

func New() *App {
	return &App{}
}

func (c *App) Init() {
	//Note: do not change the init priorities
	c.InitClients()
	c.InitRepositories()
	c.InitPorts()
	c.InitHandlers()
}
