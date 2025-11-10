package delivery

import (
	"microservice/internal/adapter/locale"
	"microservice/internal/adapter/logger"
	"microservice/internal/core/domain"
	"microservice/internal/core/port"
	"microservice/internal/driver/dto"
	"microservice/internal/server/http/status"
	"microservice/pkg/meta"

	"github.com/gin-gonic/gin"
)

type (
	ITodoHandler interface {
		Create(ctx *gin.Context)
		GetDetails(ctx *gin.Context)
		GetList(ctx *gin.Context)
	}

	TodoHandler struct {
		lgr    logger.ILogger
		l      locale.ILocale
		todoUC port.ITodoUsecase
	}
)

func NewTodo(lgr logger.ILogger, l locale.ILocale, todoUC port.ITodoUsecase) ITodoHandler {
	return &TodoHandler{lgr: lgr, l: l, todoUC: todoUC}
}

// Create godoc
// @Summary Create New Todo
// @Tags Todo
// @Accept json
// @Produce json
// @Param Request body dto.CreateRequest true "necessary fields for request"
// @Success 201 {object} meta.Response{data=dto.CreateResponse, error=nil} "success response"
// @Failure	400 {object} meta.Response{data=nil} "process failure"
// @Failure	404 {object} meta.Response{data=nil} "not found"
// @Failure	409 {object} meta.Response{data=nil} "already exists"
// @Failure	422 {object} meta.Response{data=nil} "unprocessable"
// @Router /api/v1/todo/create [post]
func (h *TodoHandler) Create(ctx *gin.Context) {
	req, err := meta.ReqBodyToDomain[*dto.CreateRequest, domain.Todo](ctx)
	if err != nil {
		meta.Resp(ctx, h.l).Status(status.Validate).Err(err).Json()
		return
	}

	res, ucErr := h.todoUC.Create(ctx, req)
	if ucErr != nil {
		meta.Resp(ctx, h.l).ServiceErr(ucErr).Json()
		return
	}

	meta.Resp(ctx, h.l).Data(dto.CreateResp(res)).Status(status.Created).Json()
	return
}

// GetDetails godoc
// @Summary Get Todo Details
// @Tags Todo
// @Accept json
// @Produce json
// @Param uuid path string true "Todo UUID" example(f81eee2d-2cca-4169-8062-7404a78d5c3b)
// @Success 200 {object}  meta.Response{data=dto.DetailResponse, error=nil} "success response"
// @Failure	400 {object} meta.Response{data=nil} "process failure"
// @Failure	404 {object} meta.Response{data=nil} "not found"
// @Failure	422 {object} meta.Response{data=nil} "invalid data types"
// @Router /api/v1/todo/{uuid} [get]
func (h *TodoHandler) GetDetails(ctx *gin.Context) {
	req, err := meta.ReqRouteParamsToDomain[*dto.DetailUriRequest, domain.Todo](ctx)
	if err != nil {
		meta.Resp(ctx, h.l).Status(status.Validate).Err(err).Json()
		return
	}

	id := req.UUID()
	res, ucErr := h.todoUC.Detail(ctx, &id)
	if ucErr != nil {
		meta.Resp(ctx, h.l).ServiceErr(ucErr).Json()
		return
	}

	meta.Resp(ctx, h.l).Data(dto.DetailResp(res)).Json()
	return
}

// GetList godoc
// @Summary Get Todos List
// @Tags Todo
// @Accept json
// @Produce json
// @Param page query int false "Page Number"
// @Param limit query int false "Page Limit"
// @Param sort query string false "`id` `description` `created_at` `updated_at`"
// @Param order query string false "`asc` or `desc`"
// @Param search query string false "Search the Description"
// @Success 200 {object}  meta.Response{data=dto.TodoListResponse, error=nil} "success response"
// @Failure	400 {object} meta.Response{data=nil} "process failure"
// @Failure	422 {object} meta.Response{data=nil} "database error while retrieving"
// @Router /api/v1/todo/list [get]
func (h *TodoHandler) GetList(ctx *gin.Context) {
	qp, err := meta.ReqQryParamToDomain[*dto.TodoListQryRequest, domain.TodoListReqQryParam](ctx)
	if err != nil {
		meta.Resp(ctx, h.l).Status(status.Validate).Err(err).Json()
		return
	}

	res, ucErr := h.todoUC.GetList(ctx, qp)
	if ucErr != nil {
		meta.Resp(ctx, h.l).ServiceErr(ucErr).Json()
		return
	}

	meta.Resp(ctx, h.l).Data(dto.TodoListResp(qp, res)).Json()
	return
}
