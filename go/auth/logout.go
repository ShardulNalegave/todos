package auth

import (
	"github.com/ShardulNalegave/todos/go/database"
	"gorm.io/gorm"
)

func Logout(db *gorm.DB, id string) {
	db.Delete(&database.Session{}, "id = ?", id)
}
