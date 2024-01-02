package auth

import (
	"github.com/ShardulNalegave/todos/go/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginUser(db *gorm.DB, inp LoginUserInp) (string, string, error) {
	var user database.User
	res := db.First(&user)
	if res.Error != nil {
		return "", "", res.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(inp.Password))
	if err != nil {
		return "", "", err
	}

	session := database.Session{
		UserID: user.ID,
	}
	res = db.Create(&session)
	if res.Error != nil {
		return "", "", res.Error
	}

	return session.ID, user.ID, nil
}

type LoginUserInp struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
