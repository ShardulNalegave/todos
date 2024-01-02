package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./todos.sqlite"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Session{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Todo{})

	return db, nil
}
