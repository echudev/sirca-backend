package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

// Configuración de la conexión a la base de datos.
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Configura y retorna una conexión a la base de datos
func ConnectDB() (*sql.DB, error) {
	dbConfig := DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode,
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

// Función auxiliar para obtener variables de entorno y convertirlas a enteros
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if valueStr == "" {
		return defaultVal
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultVal
	}
	return value
}
