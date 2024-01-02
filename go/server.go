package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ShardulNalegave/todos/go/database"
	"github.com/go-chi/chi/v5"
)

const defaultPort = "5000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Could not connect to db")
	}

	router := chi.NewRouter()
	router.Use(database.DatabaseMiddleware(db))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	log.Printf("Listening at :%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
