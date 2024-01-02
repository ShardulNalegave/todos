package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoModel struct {
	ID        string `gorm:"primaryKey" json:"id"`
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
	CreatedBy string `json:"created_by"`
}

func (t *TodoModel) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.NewString()
	return
}
