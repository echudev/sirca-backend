package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"

	"echudev/sirca-backend/db"
)

var queries *db.Queries

func main() {
	connStr := "postgresql://user:password@localhost/itemsdb?sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	queries = db.New(conn)

	http.HandleFunc("/items", handleItems)
	http.HandleFunc("/items/", handleItem)

	log.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := queries.ListItems(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(items)
	case http.MethodPost:
		var item db.CreateItemParams
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		createdItem, err := queries.CreateItem(r.Context(), item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(createdItem)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/items/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		item, err := queries.GetItem(r.Context(), int32(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(item)
	case http.MethodPut:
		var item db.UpdateItemParams
		item.ID = int32(id)
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		updatedItem, err := queries.UpdateItem(r.Context(), item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(updatedItem)
	case http.MethodDelete:
		if err := queries.DeleteItem(r.Context(), int32(id)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
