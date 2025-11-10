package model

import "time"

type Todos struct {
	BaseSql
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
}

func NewTodo() *Todos { return &Todos{} }

func (m *Todos) TableName() string { return "todos" }
