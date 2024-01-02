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

func TodosRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", allTodos)
	router.Post("/", createTodo)
	router.Get("/{id}", getTodo)
	router.Delete("/{id}", deleteTodo)
	router.Put("/{id}", updateTodo)
	return router
}

func allTodos(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	state := r.Context().Value(utils.AuthKey).(auth.AuthState)

	if !state.IsAuth {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User unauthenticated"))
		return
	}

	var todos []database.Todo
	if err := db.Find(&todos, "created_by = ?", state.Session.UserID).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not read data"))
		return
	}

	todosJSON, err := json.Marshal(todos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse data"))
		return
	}

	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(http.StatusOK)
	w.Write(todosJSON)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	state := r.Context().Value(utils.AuthKey).(auth.AuthState)

	if !state.IsAuth {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User unauthenticated"))
		return
	}

	var todo database.Todo
	res := db.First(&todo, "id = ? AND created_by = ?", id, state.Session.UserID)

	if res.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
		return
	}

	data, err := json.Marshal(todo)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
		return
	}

	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	state := r.Context().Value(utils.AuthKey).(auth.AuthState)

	var data CreateTodoBody
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Data Provided"))
		return
	}

	if !state.IsAuth {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User unauthenticated"))
		return
	}

	todo := database.Todo{
		Content:   data.Content,
		Completed: data.Completed,
		CreatedBy: state.Session.UserID,
	}
	db.Create(&todo)

	todoJSON, err := json.Marshal(todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse todo"))
		return
	}

	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(http.StatusOK)
	w.Write(todoJSON)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	state := r.Context().Value(utils.AuthKey).(auth.AuthState)

	var data UpdateTodoBody
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Data Provided"))
		return
	}

	var todo database.Todo
	if err := db.First(&todo).Error; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Todo doesn't exist"))
		return
	}

	if todo.CreatedBy != state.Session.UserID {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Cannot edit other's todos"))
		return
	}

	todo.Completed = data.Completed
	todo.Content = data.Content

	db.Save(&todo)

	todoJSON, err := json.Marshal(todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse todo"))
		return
	}

	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(http.StatusOK)
	w.Write(todoJSON)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	state := r.Context().Value(utils.DatabaseKey).(auth.AuthState)

	var todo database.Todo
	db.First(&todo, "id = ?", id)
	if todo.CreatedBy != state.Session.UserID {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Cannot delete other's todos"))
		return
	}

	res := db.Delete(&database.Todo{}, "id = ?", id)

	if res.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done"))
}

type CreateTodoBody struct {
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
}

type UpdateTodoBody struct {
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
}
