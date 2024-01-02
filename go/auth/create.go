package auth

import (
	"github.com/ShardulNalegave/todos/go/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, inp CreateUserInp) (string, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(inp.Password), 10)
	if err != nil {
		return "", "", err
	}

	user := database.User{
		Name:         inp.Name,
		Email:        inp.Email,
		PasswordHash: string(hash),
	}
	res := db.Create(&user)
	if res.Error != nil {
		return "", "", res.Error
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

type CreateUserInp struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
