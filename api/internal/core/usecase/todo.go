package usecase

import (
	"context"
	"microservice/internal/adapter/locale"
	"microservice/internal/adapter/logger"
	"microservice/internal/core/domain"
	"microservice/internal/core/port"

	"github.com/google/uuid"
)

type TodoUsecase struct {
	lgr      logger.ILogger
	l        locale.ILocale
	todoRepo port.ITodoRepository
}

func NewTodo(lgr logger.ILogger, l locale.ILocale, todoRepo port.ITodoRepository) port.ITodoUsecase {
	return &TodoUsecase{l: l, lgr: lgr, todoRepo: todoRepo}
}

func (uc *TodoUsecase) Create(ctx context.Context, ent *domain.Todo) (res *domain.Todo, err error) {
	item, txErr := uc.todoRepo.Create(ctx, ent)
	if txErr != nil {
		err = txErr
		return
	}

	res = item
	return
}

func (uc *TodoUsecase) Detail(ctx context.Context, id *uuid.UUID) (res *domain.Todo, err error) {
	item, txErr := uc.todoRepo.GetByUUID(ctx, id)
	if txErr != nil {
		err = txErr
		return
	}

	res = item
	return
}

func (uc *TodoUsecase) GetList(ctx context.Context, qp *domain.TodoListReqQryParam) (res *domain.TodoList, err error) {
	items, txErr := uc.todoRepo.GetList(ctx, qp)
	if txErr != nil {
		err = txErr
		return
	}

	// NOTE: it returns the empty slice for not-found result
	res = items
	return
}
