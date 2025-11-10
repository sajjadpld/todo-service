package dto

import "microservice/internal/core/domain"

type ListQryRequest struct {
	Page   int    `form:"page" binding:"omitempty,numeric" json:"page"`   // integer value
	Limit  int    `form:"limit" binding:"omitempty,numeric" json:"limit"` // integer value
	Sort   string `form:"sort" binding:"omitempty,ascii" json:"sort"`
	Order  string `form:"order" binding:"omitempty,ascii" json:"order"` // "asc" or "desc"
	Search string `form:"search" binding:"omitempty,alphanum" json:"search"`
}

func (r *ListQryRequest) EvalBaseQry() domain.ReqBaseQryParam {
	d := new(domain.ReqBaseQryParam)

	if r.Page != 0 {
		d.SetPage(&r.Page)
	}

	if r.Limit != 0 {
		d.SetLimit(&r.Limit)
	}

	if len(r.Sort) > 0 {
		d.SetSort(&r.Sort)
	}

	if len(r.Order) > 0 {
		d.SetOrder(&r.Order)
	}

	if len(r.Search) > 0 {
		d.SetSearch(&r.Search)
	}

	return *d
}
