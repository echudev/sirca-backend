package main

import (
	"context"
	"echudev/sirca-backend/internal/db"
	"echudev/sirca-backend/internal/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}

func run() error {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using environment variables.")
	}

	// Connect to the database
	connPool, err := db.ConnectDB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer connPool.Close()

	// Create queries with the connection pool
	queries := db.New(connPool)

	// Create and configure the server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      setupRoutes(queries, connPool),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server running on http://localhost%s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("ListenAndServe error: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	return gracefulShutdown(srv)
}

func setupRoutes(queries *db.Queries, pool *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()

	// Define routes with HTTP verbs
	mux.HandleFunc("GET /items", handlers.GetItems(queries))
	mux.HandleFunc("GET /analyzers", handlers.GetAnalyzers(queries))
	mux.HandleFunc("POST /analyzers", handlers.CreateAnalyzer(queries, pool))
	mux.HandleFunc("GET /stations", handlers.GetStations(queries))
	return mux
}

func gracefulShutdown(srv *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exited properly")
	return nil
}
