package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionModel struct {
	ID     string `gorm:"primaryKey"`
	UserID string
}

func (t *SessionModel) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.NewString()
	return
}
