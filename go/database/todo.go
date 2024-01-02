package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Todo struct {
	ID        string `gorm:"primaryKey" json:"id"`
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
	CreatedBy string `json:"created_by"`
}

func (t *Todo) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.NewString()
	return
}
