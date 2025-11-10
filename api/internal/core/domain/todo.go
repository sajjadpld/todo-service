package domain

import (
	"microservice/internal/adapter/orm/model"
	"time"
)

type (
	Todo struct {
		Base
		description *string
		dueDate     *time.Time
	}

	TodoList struct {
		total int64
		list  []*Todo
	}
)

func NewTodo() *Todo {
	return &Todo{}
}

func (d *Todo) Domain() *Todo {
	return d
}

func (d *Todo) Description() *string {
	return d.description
}

func (d *Todo) SetDescription(description *string) {
	d.description = description
}

func (d *Todo) DueDate() *time.Time {
	return d.dueDate
}

func (d *Todo) SetDueDate(dueDate *time.Time) {
	d.dueDate = dueDate
}

//

func (d *Todo) FromDB(src *model.Todos) *Todo {
	if src == nil {
		return nil
	}

	// base
	d.SetID(&src.ID)
	d.SetUUID(&src.Uuid)
	d.SetCreatedAt(&src.CreatedAt)
	d.SetUpdatedAt(&src.UpdatedAt)
	d.SetDeletedAt(&src.DeletedAt.Time)
	//fields
	d.SetDescription(&src.Description)
	d.SetDueDate(&src.DueDate)
	return d
}

func (d *Todo) ToDB() *model.Todos {
	return &model.Todos{
		BaseSql: model.BaseSql{
			Uuid: d.UUID(),
		},
		Description: *d.Description(),
		DueDate:     *d.DueDate(),
	}
}

//

func NewTodoList() *TodoList { return &TodoList{} }

func (ul *TodoList) SetTotal(total int64) { ul.total = total }

func (ul *TodoList) Total() int64 { return ul.total }

func (ul *TodoList) SetList(list []*Todo) { ul.list = list }

func (ul *TodoList) List() []*Todo { return ul.list }

//

func (ul *TodoList) ListFromDB(src []*model.Todos) []*Todo {
	ul.list = make([]*Todo, 0)

	if len(src) == 0 {
		return ul.list
	}

	for _, u := range src {
		ul.list = append(ul.list, NewTodo().FromDB(u))
	}

	return ul.list
}

// Query Params

type TodoListReqQryParam struct {
	ReqBaseQryParam
}

func NewTodoListReqQryParam() *TodoListReqQryParam {
	return &TodoListReqQryParam{}
}
