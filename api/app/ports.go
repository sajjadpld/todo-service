package app

import (
	"microservice/internal/core/port"
	"microservice/internal/core/usecase"
)

type Ports struct {
	TodoUC port.ITodoUsecase
}

func (c *App) InitPorts() {
	c.port = new(Ports)
	c.port.TodoUC = usecase.NewTodo(c.logger, c.locale, c.repo.TodoRepo)
}
