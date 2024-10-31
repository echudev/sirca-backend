package handlers

import (
	"echudev/sirca-backend/internal/db"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler para Crear un Analyzer (POST /analyzers)
func CreateAnalyzer(queries *db.Queries, pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type RequestParams struct {
			Item         *db.CreateItemParams     `json:"item"`
			Analyzer     *db.CreateAnalyzerParams `json:"analyzer"`
			ItemTypeName string                   `json:"item_type_name"`
			BrandName    string                   `json:"brand_name"`
			ModelName    string                   `json:"model_name"`
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
		if req.BrandName == "" {
			http.Error(w, "Brand Name is required", http.StatusBadRequest)
			return
		}
		if req.ModelName == "" {
			http.Error(w, "Model Name is required", http.StatusBadRequest)
			return
		}
		if req.ItemTypeName == "" {
			http.Error(w, "Item Type Name is required", http.StatusBadRequest)
			return
		}

		// Genero código de inventario

		// Iniciar la transacción
		tx, err := pool.Begin(r.Context())
		if err != nil {
			http.Error(w, "Error iniciando transacción", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(r.Context()) // Rollback si algo falla

		// Crear queries con la transacción
		qtx := queries.WithTx(tx)

		// Busco id de marca en bd, con el nombre de la marca que envió el cliente
		brandID, err := qtx.GetBrandId(r.Context(), req.BrandName)
		if err != nil {
			http.Error(w, "Error getting brand ID", http.StatusInternalServerError)
			return
		}
		// Busco id de modelo en bd, con el nombre de la marca y el modelo que envió el cliente
		modelID, err := qtx.GetModelId(r.Context(), db.GetModelIdParams{BrandID: brandID, ModelName: req.ModelName})
		if err != nil {
			http.Error(w, "Error getting model ID", http.StatusInternalServerError)
			return
		}
		// Busco tipo de item en bd, con el nombre del tipo de item que envió el cliente
		itemTypeID, err := qtx.GetItemTypeId(r.Context(), req.ItemTypeName)
		if err != nil {
			http.Error(w, "Error getting item type ID", http.StatusInternalServerError)
			return
		}

		//Asignar los ID de tipo de item al campo de estructura de datos
		req.Item.ItemTypeID = itemTypeID

		// Asignar los ID de marca y modelo a los campos de la estructura de datos
		req.Analyzer.BrandID = brandID
		req.Analyzer.ModelID = modelID

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
