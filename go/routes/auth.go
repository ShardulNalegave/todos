package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ShardulNalegave/todos/go/auth"
	"github.com/ShardulNalegave/todos/go/database"
	"github.com/ShardulNalegave/todos/go/utils"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// Create a new router with all auth-related routes registered
func AuthRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/create", createUser)
	router.Post("/login", login)
	router.Post("/logout", logout)
	router.Get("/user", currentUser)
	return router
}

// GET - Gives info about currently logged in user
func currentUser(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	state := r.Context().Value(utils.AuthKey).(auth.AuthState)

	if !state.IsAuth {
		ResponseData{
			Message: "No user logged in",
		}.Write(&w, http.StatusUnauthorized)
		return
	}

	var user database.User
	db.First(&user, "id = ?", state.Session.UserID)

	WriteJSON(http.StatusOK, &w, struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

// POST - Creates a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	var data auth.CreateUserInp
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		ResponseData{
			Message: "Invalid Data Provided",
		}.Write(&w, http.StatusBadRequest)
		return
	}

	sessionID, userID, err := auth.CreateUser(db, data)
	if err != nil {
		ResponseData{
			Message: "Could not create user",
		}.Write(&w, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.AuthCookie,
		Value:    sessionID,
		HttpOnly: true,
		Path:     "/",
	})
	WriteJSON(http.StatusOK, &w, struct {
		UserID string `json:"user_id"`
	}{
		UserID: userID,
	})
}

// POST - Logs in with given credentials
func login(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	var data auth.LoginUserInp
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		ResponseData{
			Message: "Invalid Data Provided",
		}.Write(&w, http.StatusBadRequest)
		return
	}

	sessionID, _, err := auth.LoginUser(db, data)
	if err != nil {
		ResponseData{
			Message: "Could not login",
		}.Write(&w, http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.AuthCookie,
		Value:    sessionID,
		HttpOnly: true,
		Path:     "/",
	})
	ResponseData{
		Message: "Done",
	}.Write(&w, http.StatusOK)
}

// POST - Logs out
func logout(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)

	cookie, err := r.Cookie(auth.AuthCookie)
	if err != nil {
		ResponseData{
			Message: "No user logged in",
		}.Write(&w, http.StatusBadRequest)
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
	ResponseData{
		Message: "Done",
	}.Write(&w, http.StatusOK)
}
