package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ShardulNalegave/todos/go/database"
	"github.com/ShardulNalegave/todos/go/utils"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func TodosRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", CreateTodo)
	router.Get("/{id}", GetTodo)
	router.Delete("/{id}", DeleteTodo)
	router.Put("/{id}", UpdateTodo)
	return router
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)

	var todo database.TodoModel
	res := db.First(&todo, "id = ?", id)

	if res.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}

	data, err := json.Marshal(todo)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	var data CreateTodoBody
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Data Provided"))
	}

	panic("Unimplemented")
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(utils.DatabaseKey).(*gorm.DB)
	var data UpdateTodoBody
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Data Provided"))
	}

	panic("Unimplemented")
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	db := r.Context().Value(utils.DatabaseKey).(*gorm.DB)

	res := db.Delete(&database.TodoModel{}, "id = ?", id)

	if res.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
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
