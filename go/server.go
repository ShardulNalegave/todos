package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ShardulNalegave/todos/go/auth"
	"github.com/ShardulNalegave/todos/go/database"
	"github.com/ShardulNalegave/todos/go/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

	corsHandler := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	})

	router := chi.NewRouter()
	router.Use(corsHandler)

	router.Use(database.DatabaseMiddleware(db))
	router.Use(auth.AuthMiddleware())

	router.Mount("/auth", routes.AuthRoutes())
	router.Mount("/todos", routes.TodosRoutes())

	log.Printf("Listening at :%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
