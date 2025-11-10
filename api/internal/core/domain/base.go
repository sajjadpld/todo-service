package domain

import (
	"github.com/google/uuid"
	"time"
)

type Base struct {
	id        *uint
	uuid      *uuid.UUID
	createdAt *time.Time
	updatedAt *time.Time
	deletedAt *time.Time
}

func NewBase() *Base {
	return &Base{}
}

func (bd *Base) SetID(id *uint) {
	bd.id = id
}

func (bd *Base) ID() uint {
	if bd.id != nil {
		return *bd.id
	}

	return 0
}

func (bd *Base) SetUUID(id *uuid.UUID) {
	bd.uuid = id
}

func (bd *Base) UUID() uuid.UUID {
	if bd.uuid != nil {
		return *bd.uuid
	}

	return uuid.Nil
}

func (bd *Base) SetCreatedAt(t *time.Time) {
	bd.createdAt = t
}

func (bd *Base) CreatedAt() time.Time {
	if bd.createdAt != nil {
		return *bd.createdAt
	}

	return time.Time{} // zero-time
}

func (bd *Base) SetUpdatedAt(t *time.Time) {
	bd.updatedAt = t
}

func (bd *Base) UpdatedAt() time.Time {
	if bd.updatedAt != nil {
		return *bd.updatedAt
	}

	return time.Time{} // zero-time
}

func (bd *Base) SetDeletedAt(t *time.Time) {
	bd.deletedAt = t
}

func (bd *Base) DeletedAt() time.Time {
	if bd.deletedAt != nil {
		return *bd.deletedAt
	}

	return time.Time{} // zero-time
}

// collection default query params

type ReqBaseQryParam struct {
	page   *int
	limit  *int
	order  *string
	sort   *string
	search *string
}

func (bc *ReqBaseQryParam) SetPage(page *int) {
	bc.page = page
}

func (bc *ReqBaseQryParam) Page() int {
	if bc.page != nil {
		return *bc.page
	}

	return 1
}

func (bc *ReqBaseQryParam) SetLimit(limit *int) {
	bc.limit = limit
}

func (bc *ReqBaseQryParam) Limit() int {
	if bc.limit != nil {
		if *bc.limit > 50 {
			*bc.limit = 50
		}

		return *bc.limit
	}

	return 10 // default items count per page
}

func (bc *ReqBaseQryParam) SetOrder(order *string) {
	bc.order = order
}

// Order default order: desc
func (bc *ReqBaseQryParam) Order() string {
	if bc.order != nil {
		return *bc.order
	}

	return "desc"
}

func (bc *ReqBaseQryParam) SetSort(sort *string) { bc.sort = sort }

// Sort related field to be sorted
func (bc *ReqBaseQryParam) Sort() string {
	if bc.sort != nil {
		return *bc.sort
	}

	return "created_at"
}

func (bc *ReqBaseQryParam) SetSearch(search *string) { bc.search = search }

func (bc *ReqBaseQryParam) Search() string {
	if bc.search != nil {
		return *bc.search
	}

	return ""
}
