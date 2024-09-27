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
	// Cargar .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	// Conectar a la base de datos
	conn, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Pasar el pool de conexiones al código generado por sqlc
	queries := db.New(conn)

	// Crear el mux
	mux := http.NewServeMux()

	// Definir rutas con el verbo HTTP
	mux.HandleFunc("GET /items", handlers.GetItems(queries))
	mux.HandleFunc("POST /items", handlers.CreateItem(queries))
	mux.HandleFunc("GET /items/{id}", handlers.GetItem(queries))
	mux.HandleFunc("PUT /items/{id}", handlers.UpdateItem(queries))
	mux.HandleFunc("DELETE /items/{id}", handlers.DeleteItem(queries))

	// Iniciar el servidor
	log.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
