package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("todos.sqlite"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}