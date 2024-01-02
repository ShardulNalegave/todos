package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ShardulNalegave/todos/go/auth"
	"github.com/ShardulNalegave/todos/go/utils"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func AuthRoutes(cors func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors)
	router.Post("/create", createUser)
	router.Post("/login", login)
	router.Post("/logout", logout)
	return router
}

func createUser(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	var data auth.CreateUserInp
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Data Provided"))
		return
	}

	sessionID, _, err := auth.CreateUser(db, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not create user"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.AuthCookie,
		Value:    sessionID,
		HttpOnly: true,
		Path:     "/",
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done"))
}

func login(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	var data auth.LoginUserInp
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Data Provided"))
		return
	}

	sessionID, _, err := auth.LoginUser(db, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not login"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.AuthCookie,
		Value:    sessionID,
		HttpOnly: true,
		Path:     "/",
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done"))
}

func logout(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)

	cookie, err := r.Cookie(auth.AuthCookie)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not logged in"))
		return
	}

	auth.Logout(db, cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:     auth.AuthCookie,
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done"))
}
