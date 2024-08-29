package main

import (
	"echudev/sirca-backend/db"
	"echudev/sirca-backend/handlers"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar las variables del archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	// Conectar a la base de datos
	conn, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(conn)

	http.HandleFunc("/items", handlers.HandleItems(queries))
	http.HandleFunc("/items/", handlers.HandleItem(queries))

	log.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
