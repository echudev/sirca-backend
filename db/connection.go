package db

import (
	"database/sql"
	"fmt"
	"log"

	"os"

	_ "github.com/lib/pq"
)

// Configura y retorna una conexión a la base de datos
func ConnectDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Verifica que la conexión esté activa
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	log.Println("Conexión a la base de datos exitosa")
	return conn, nil
}
