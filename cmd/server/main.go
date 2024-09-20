package main

import (
	"echudev/sirca-backend/db"
	"echudev/sirca-backend/internal/database"
	"echudev/sirca-backend/internal/handlers"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	// Connect database
	conn, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(conn)

	// Item routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /items", handlers.GetItems(queries))
	mux.HandleFunc("POST /items", handlers.CreateItem(queries))
	mux.HandleFunc("GET /items/{id}", handlers.GetItem(queries))
	mux.HandleFunc("PUT /items/{id}", handlers.UpdateItem(queries))
	mux.HandleFunc("DELETE /items/{id}", handlers.DeleteItem(queries))

	// server init
	log.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
