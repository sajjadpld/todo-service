package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseSql Shared Fields of all model have to be added here
type BaseSql struct {
	gorm.Model
	Uuid uuid.UUID `json:"uuid"`
}
