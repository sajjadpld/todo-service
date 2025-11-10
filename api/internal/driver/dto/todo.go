package dto

import (
	"fmt"
	"github.com/google/uuid"
	"math"
	"microservice/internal/core/domain"
	"time"
)

type CreateRequest struct {
	Description string `json:"description" validate:"required,ascii" example:"Create new todo item"`
	DueDate     string `json:"dueDate" validate:"required" example:"2025-08-07 10:11:12"`
}

func (dto *CreateRequest) ToDomain() *domain.Todo {
	d := domain.NewTodo()
	d.SetDescription(&dto.Description)

	dateTime, _ := time.Parse(time.DateTime, dto.DueDate)
	d.SetDueDate(&dateTime)

	return d
}

type CreateResponse struct {
	Uuid        string `json:"uuid" example:"e48c48a3-cb72-4d64-b035-5c30fc900ef6"`
	Description string `json:"description" example:"Create new todo item"`
	DueDate     string `json:"dueDate" example:"2025-08-07 10:11:12"`
}

func CreateResp(src *domain.Todo) *CreateResponse {
	return &CreateResponse{
		Uuid: func() string {
			if src.UUID() == uuid.Nil {
				return ""
			}

			return src.UUID().String()
		}(),
		Description: *src.Description(),
		DueDate:     src.DueDate().Format(time.RFC3339),
	}
}

//

type DetailUriRequest struct {
	Uuid string `param:"uuid" validate:"required,uuid" example:"bf56c6b6-dd02-47ba-8dc4-bd7d2843a77a"`
}

func (dto *DetailUriRequest) ToDomain() *domain.Todo {
	id := uuid.MustParse(dto.Uuid)

	d := domain.NewTodo()
	d.SetUUID(&id)
	return d
}

type DetailResponse struct {
	Uuid        string `json:"uuid" example:"e48c48a3-cb72-4d64-b035-5c30fc900ef6"`
	Description string `json:"description" example:"Create new todo item"`
	DueDate     string `json:"dueDate" example:"2025-08-07 10:11:12"`
}

func DetailResp(src *domain.Todo) *DetailResponse {
	return &DetailResponse{
		Uuid: func() string {
			if src.UUID() == uuid.Nil {
				return ""
			}

			return src.UUID().String()
		}(),
		Description: *src.Description(),
		DueDate:     src.DueDate().Format(time.RFC3339),
	}
}

//

type TodoListQryRequest struct {
	ListQryRequest
}

func (r *TodoListQryRequest) ToDomain() *domain.TodoListReqQryParam {
	qry := domain.NewTodoListReqQryParam()
	qry.ReqBaseQryParam = r.EvalBaseQry()

	return qry
}

type (
	TodoListItemDetail struct {
		Uuid        string `json:"id" example:"02bda2f0-61e5-483c-a2d8-15eafb00b945"`
		Description string `json:"name" example:"Create new todo..."`
		DueDate     string `json:"dueDate" validate:"required,ascii" example:"2025-08-07 10:11:12"`
	}

	TodoListResponse struct {
		Page  int                   `json:"page" example:"1"`
		Limit int                   `json:"limit" example:"10"`
		Pages int                   `json:"pages" example:"3"`
		Total int64                 `json:"total" example:"27"`
		Todos []*TodoListItemDetail `json:"todos"`
	}
)

func TodoListResp(qry *domain.TodoListReqQryParam, src *domain.TodoList) *TodoListResponse {
	list := new(TodoListResponse)
	list.Page = qry.Page()
	list.Limit = qry.Limit()
	list.Pages = int(math.Ceil(float64(src.Total()) / float64(qry.Limit())))
	list.Total = src.Total()
	list.Todos = make([]*TodoListItemDetail, 0)

	if len(src.List()) > 0 {
		for _, todo := range src.List() {
			desc := *todo.Description()

			if len(desc) > 20 {
				desc = fmt.Sprintf("%s...", desc[:20])
			}

			list.Todos = append(list.Todos, &TodoListItemDetail{
				Uuid:        todo.UUID().String(),
				Description: desc,
				DueDate:     todo.DueDate().Format(time.RFC3339),
			})
		}
	}

	return list
}
