package handlers

import (
	"echudev/sirca-backend/internal/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler que devuelve todos los Analyzers(GET /analyzers)
func GetAnalyzers(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		analyzers, err := queries.GetAnalyzers(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(analyzers)
	}
}

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
		if !req.Item.ItemAdquisitionDate.Valid {
			http.Error(w, "Date is required", http.StatusBadRequest)
			return
		}
		if req.Item.ItemAdquisitionDate.Time.IsZero() {
			http.Error(w, "Invalid date", http.StatusBadRequest)
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

		//Asigno el ID de tipo de item al campo de estructura de datos
		req.Item.ItemTypeID = itemTypeID
		// Asigno los ID de marca y modelo a los campos de la estructura de datos
		req.Analyzer.BrandID = brandID
		req.Analyzer.ModelID = modelID

		// Creo el item
		itemID, err := qtx.CreateItem(r.Context(), *req.Item)
		if err != nil {
			http.Error(w, "Error creating item", http.StatusInternalServerError)
			return
		}

		// Genero el código de inventario
		itemCode, err := GenerateInventaryCode(req.ItemTypeName, req.BrandName, req.ModelName, req.Item.ItemAdquisitionDate, itemID)
		if err != nil {
			http.Error(w, "Error generating item code", http.StatusInternalServerError)
			return
		}

		// Actualizar el registro con el código
		itemCode, err = qtx.UpdateInventaryCode(r.Context(), db.UpdateInventaryCodeParams{ItemID: itemID, ItemCode: itemCode})
		if err != nil {
			http.Error(w, "Error updating inventary code", http.StatusInternalServerError)
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
			"code":     itemCode,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}

func GenerateInventaryCode(reqTypeName string, reqBrandName string, reqModelName string, reqItemAdquisitonDate pgtype.Date, reqItemId int32) (string, error) {
	// Validar que los strings tengan al menos 3 caracteres
	if len(reqTypeName) < 3 {
		return "", fmt.Errorf("type name must be at least 3 characters long")
	}
	if len(reqModelName) < 3 {
		return "", fmt.Errorf("model name must be at least 3 characters long")
	}
	if len(reqBrandName) < 3 {
		return "", fmt.Errorf("brand name must be at least 3 characters long")
	}

	// Validar fecha
	if !reqItemAdquisitonDate.Valid {
		return "", fmt.Errorf("invalid date")
	}

	// Validar ID
	if reqItemId <= 0 {
		return "", fmt.Errorf("invalid item ID")
	}

	type_code := reqTypeName[:3]
	brand_code := reqBrandName[:3]
	model_code := reqModelName[:3]
	year_code := strconv.Itoa(reqItemAdquisitonDate.Time.Year())
	id_code := strconv.Itoa(int(reqItemId))

	item_code := strings.ToUpper(type_code + "-" + brand_code + "-" + model_code + "-" + year_code + "-" + id_code)

	return item_code, nil
}
