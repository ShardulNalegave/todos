package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           string `gorm:"primaryKey" json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string
}

func (t *User) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.NewString()
	return
}
