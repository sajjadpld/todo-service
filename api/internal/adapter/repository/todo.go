package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"microservice/internal/adapter/locale"
	"microservice/internal/adapter/logger"
	"microservice/internal/adapter/orm"
	"microservice/internal/adapter/orm/model"
	"microservice/internal/core/domain"
	"microservice/internal/core/port"
	"microservice/internal/server/http/status"
	"microservice/pkg/meta"
)

type TodoRepository struct {
	lgr logger.ILogger
	l   locale.ILocale
	db  orm.ISql
}

func NewTodo(l locale.ILocale, lgr logger.ILogger, db orm.ISql) port.ITodoRepository {
	return &TodoRepository{l: l, lgr: lgr, db: db}
}

func (tr *TodoRepository) Tx(db orm.ISql) { tr.db = db }

//

func (tr *TodoRepository) Create(ctx context.Context, ent *domain.Todo) (res *domain.Todo, err error) {
	tx := tr.db.C().WithContext(ctx).Model(model.Todos{})

	m := ent.ToDB()
	if txErr := tx.Omit("uuid", "deleted_at").Clauses(clause.Returning{}).Create(&m).Error; txErr != nil {
		tr.lgr.Error("todo.repo.create", zap.Error(txErr))

		if errors.Is(txErr, gorm.ErrDuplicatedKey) {
			err = meta.ServiceErr(status.ItemExist)
			return
		}

		err = meta.ServiceErr(status.Failed, txErr)
		return
	}

	res = domain.NewTodo().FromDB(m)
	return
}

func (tr *TodoRepository) GetByUUID(ctx context.Context, id *uuid.UUID) (res *domain.Todo, err error) {
	m := model.NewTodo()
	tx := tr.db.C().WithContext(ctx).Model(&model.Todos{})

	tx.First(&m, "uuid = ?", id)
	if tx.Error != nil {
		tr.lgr.Error("burrow.repo.detail", zap.Error(tx.Error))

		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			err = meta.ServiceErr(status.NotFound)
			return
		}

		err = meta.ServiceErr(status.Failed)
		return
	}

	if tx.RowsAffected == 0 {
		err = meta.ServiceErr(status.NotFound)
		return
	}

	res = domain.NewTodo().FromDB(m)
	return
}

func (tr *TodoRepository) GetList(ctx context.Context, qp *domain.TodoListReqQryParam) (res *domain.TodoList, err error) {
	list := domain.NewTodoList()

	var (
		offset = (qp.Page() - 1) * qp.Limit()
		sort   = fmt.Sprintf("%s %s", qp.Sort(), qp.Order())
		models []*model.Todos
		total  int64
	)

	tx := tr.db.C().WithContext(ctx).Model(&model.Todos{})

	if len(qp.Search()) > 0 {
		searchVal := fmt.Sprintf("%%%s%%", qp.Search()) // this returns %search_value%
		tx.Where("description LIKE ?", searchVal)
	}

	//

	count := tx.Count(&total)
	if txErr := count.Error; txErr != nil {
		tr.lgr.Error("todo.repo.list.count.total", zap.Error(txErr))
		err = meta.ServiceErr(status.Failed)
		return
	}

	items := tx.Order(sort).Offset(offset).Limit(qp.Limit()).Find(&models)

	if err = items.Error; err != nil {
		tr.lgr.Error("todo.repo.list", zap.Error(err))
		return
	}

	list.ListFromDB(models)
	list.SetTotal(total)
	res = list
	return
}
