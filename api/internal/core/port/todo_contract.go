package port

import (
	"context"
	"github.com/google/uuid"
	"microservice/internal/core/domain"
)

//go:generate mockgen -source=./todo_contract.go -destination=./mocks/todo_repository_mock.go -package=todo_repository_mock
type ITodoRepository interface {
	IRepository
	Create(ctx context.Context, ent *domain.Todo) (*domain.Todo, error)
	GetByUUID(ctx context.Context, id *uuid.UUID) (*domain.Todo, error)
	GetList(ctx context.Context, qp *domain.TodoListReqQryParam) (*domain.TodoList, error)
}

type ITodoUsecase interface {
	Create(ctx context.Context, ent *domain.Todo) (*domain.Todo, error)
	Detail(ctx context.Context, id *uuid.UUID) (*domain.Todo, error)
	GetList(ctx context.Context, qp *domain.TodoListReqQryParam) (*domain.TodoList, error)
}
