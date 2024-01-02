package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID     string `gorm:"primaryKey"`
	UserID string
}

func (t *Session) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.NewString()
	return
}
