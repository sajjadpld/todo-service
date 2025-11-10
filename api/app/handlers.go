package app

import "microservice/internal/driver/delivery"

type HttpHandlers struct {
	TodoHandler delivery.ITodoHandler
}

func (c *App) InitHandlers() {
	c.httpHandlers = new(HttpHandlers)
	c.httpHandlers.TodoHandler = delivery.NewTodo(c.logger, c.locale, c.port.TodoUC)
}

func (c *App) HttpHandlers() *HttpHandlers {
	return c.httpHandlers
}
