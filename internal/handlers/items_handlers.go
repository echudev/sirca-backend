package handlers

import (
	"echudev/sirca-backend/internal/db"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
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

// Handler para crear un item (POST /items)
func CreateItem(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newItem db.CreateItemParams
		// Cerrar el cuerpo de la solicitud al finalizar
		defer r.Body.Close()

		// Decodificar el JSON de la solicitud
		if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Validaciones
		if newItem.ItemName == "" {
			http.Error(w, "Nombre is required", http.StatusBadRequest)
			return
		}

		// Crear el item en la base de datos
		newItemId, err := queries.CreateItem(r.Context(), newItem)
		if err != nil {
			http.Error(w, "Error creating item in the database", http.StatusInternalServerError)
			return
		}

		// Configurar la cabecera de respuesta y devolver el ID en JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":                "success",
			"id":                    newItemId,
			"item_name":             newItem.ItemName,
			"item_type_id":          newItem.ItemTypeID,
			"item_description":      newItem.ItemDescription,
			"item_code":             newItem.ItemCode,
			"item_adquisition_date": newItem.ItemAdquisitionDate,
		})
	}
}

// Handler para Crear un Analyzer (POST /analyzers)
func CreateAnalyzer(queries *db.Queries, pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type RequestParams struct {
			Item      *db.CreateItemParams     `json:"item"`
			Analyzer  *db.CreateAnalyzerParams `json:"analyzer"`
			BrandName string                   `json:"brand"`
			ModelName string                   `json:"model"`
		}
		var req RequestParams

		defer r.Body.Close() // Cierra el cuerpo de la solicitud al finalizar

		// Decodifica el JSON de la solicitud
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decodificando JSON: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Validaciones
		if req.Item.ItemName == "" {
			http.Error(w, "Item Name is required", http.StatusBadRequest)
			return
		}
		if req.Item.ItemAdquisitionDate.Time.Format("2006-01-02") == "" {
			http.Error(w, "Item Adquisition Date is required", http.StatusBadRequest)
			return
		}
		if req.Analyzer.AnalyzerSerialnumber == "" {
			http.Error(w, "Analyzer Serial number is required", http.StatusBadRequest)
			return
		}
		if req.Analyzer.AnalyzerPollutant == "" {
			http.Error(w, "Analyzer Pollutant is required", http.StatusBadRequest)
			return
		}
		if req.Analyzer.AnalyzerStateID == 0 {
			http.Error(w, "Analyzer State is required", http.StatusBadRequest)
			return
		}

		// Iniciar la transacción
		tx, err := pool.Begin(r.Context())
		if err != nil {
			http.Error(w, "Error iniciando transacción", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(r.Context()) // Rollback si algo falla

		// Crear queries con la transacción
		qtx := queries.WithTx(tx)

		// Busco id de marca y modelo con el nombre de la marca y modelo que envió el cliente
		brandAndModelID, err := qtx.GetBrandAndModelId(r.Context(), db.GetBrandAndModelIdParams{
			BrandName: req.BrandName,
			ModelName: req.ModelName,
		})
		if err != nil {
			http.Error(w, "Error getting brand and model ID", http.StatusInternalServerError)
			return
		}
		// Asignar los ID de marca y modelo a los campos de la estructura de datos
		req.Analyzer.BrandID = brandAndModelID.BrandID
		req.Analyzer.ModelID = brandAndModelID.ModelID

		// Creo el item
		itemID, err := qtx.CreateItem(r.Context(), *req.Item)
		if err != nil {
			http.Error(w, "Error creating item", http.StatusInternalServerError)
			return
		}

		// Usar el ID del item creado para el analyzer
		req.Analyzer.ItemID = itemID

		// Crear el analyzer
		analyzerID, err := qtx.CreateAnalyzer(r.Context(), *req.Analyzer)
		if err != nil {
			http.Error(w, "Error creating analyzer", http.StatusInternalServerError)
			return
		}

		// Commit de la transacción
		if err := tx.Commit(r.Context()); err != nil {
			http.Error(w, "Error committing transaction", http.StatusInternalServerError)
			return
		}

		// Preparar y enviar la respuesta
		response := map[string]interface{}{
			"item":     itemID,
			"analyzer": analyzerID,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
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
