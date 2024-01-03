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

// Creates a new router with all todo routes registered
func TodosRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", allTodos)
	router.Post("/", createTodo)
	router.Get("/{id}", getTodo)
	router.Delete("/{id}", deleteTodo)
	router.Put("/{id}", updateTodo)
	return router
}

// GET - Returns all todos of the currently logged in user
func allTodos(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	state := r.Context().Value(utils.AuthKey).(auth.AuthState)

	if !state.IsAuth {
		ResponseData{
			Message: "User not logged in",
		}.Write(&w, http.StatusUnauthorized)
		return
	}

	var todos []database.Todo
	if err := db.Find(&todos, "created_by = ?", state.Session.UserID).Error; err != nil {
		ResponseData{
			Message: "Could not read data",
		}.Write(&w, http.StatusInternalServerError)
		return
	}

	WriteJSON(http.StatusOK, &w, todos)
}

// GET - Returns Todo with given ID if created by currently logged in user
func getTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	state := r.Context().Value(utils.AuthKey).(auth.AuthState)

	if !state.IsAuth {
		ResponseData{
			Message: "User not logged in",
		}.Write(&w, http.StatusUnauthorized)
		return
	}

	var todo database.Todo
	res := db.First(&todo, "id = ? AND created_by = ?", id, state.Session.UserID)

	if res.Error != nil {
		ResponseData{
			Message: "Todo not found",
		}.Write(&w, http.StatusNotFound)
		return
	}

	WriteJSON(http.StatusOK, &w, todo)
}

// POST - Creates a new Todo
func createTodo(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	state := r.Context().Value(utils.AuthKey).(auth.AuthState)

	var data CreateTodoBody
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		ResponseData{
			Message: "Invalid Data Provided",
		}.Write(&w, http.StatusBadRequest)
		return
	}

	if !state.IsAuth {
		ResponseData{
			Message: "User not logged in",
		}.Write(&w, http.StatusUnauthorized)
		return
	}

	todo := database.Todo{
		Content:   data.Content,
		Completed: false,
		CreatedBy: state.Session.UserID,
	}
	db.Create(&todo)

	WriteJSON(http.StatusOK, &w, todo)
}

// PUT - Updates Todo with given ID if created by the current user
func updateTodo(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	state := r.Context().Value(utils.AuthKey).(auth.AuthState)

	var data UpdateTodoBody
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		ResponseData{
			Message: "Invalid Data Provided",
		}.Write(&w, http.StatusBadRequest)
		return
	}

	var todo database.Todo
	if err := db.First(&todo).Error; err != nil {
		ResponseData{
			Message: "Todo doesn't exist",
		}.Write(&w, http.StatusBadRequest)
		return
	}

	if todo.CreatedBy != state.Session.UserID {
		ResponseData{
			Message: "Cannot edit other's todos",
		}.Write(&w, http.StatusUnauthorized)
		return
	}

	todo.Completed = data.Completed
	todo.Content = data.Content
	db.Save(&todo)

	WriteJSON(http.StatusOK, &w, todo)
}

// DELETE - Deletes todo with given ID if created by current user
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	state := r.Context().Value(utils.AuthKey).(auth.AuthState)

	var todo database.Todo
	db.First(&todo, "id = ?", id)
	if todo.CreatedBy != state.Session.UserID {
		ResponseData{
			Message: "Cannot delete other's todos",
		}.Write(&w, http.StatusUnauthorized)
		return
	}

	res := db.Delete(&database.Todo{}, "id = ?", id)

	if res.Error != nil {
		ResponseData{
			Message: "Todo not found",
		}.Write(&w, http.StatusNotFound)
		return
	}

	ResponseData{
		Message: "Done",
	}.Write(&w, http.StatusOK)
}

type CreateTodoBody struct {
	Content string `json:"content"`
}

type UpdateTodoBody struct {
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
}
