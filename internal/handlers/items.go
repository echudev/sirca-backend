package handlers

import (
	"echudev/sirca-backend/internal/db"
	"encoding/json"
	"net/http"
)

// Handler para obtener la lista de items (GET /items)
func GetItems(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := queries.GetItems(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

// Handler para obtener estaciones (GET /stations)
func GetStations(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stations, err := queries.GetStations(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stations)
	}
}
