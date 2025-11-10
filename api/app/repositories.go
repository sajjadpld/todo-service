package app

import (
	"microservice/internal/adapter/repository"
	"microservice/internal/core/port"
)

type Repositories struct {
	TodoRepo port.ITodoRepository
}

func (c *App) InitRepositories() {
	c.repo = new(Repositories)
	c.repo.TodoRepo = repository.NewTodo(c.locale, c.logger, c.database)
}

func (c *App) Repositories() *Repositories {
	return c.repo
}
