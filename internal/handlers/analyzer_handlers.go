package handlers

import (
	"echudev/sirca-backend/internal/db"
	"echudev/sirca-backend/internal/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

func DeleteAnalyzer(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extraer el ID del path en lugar de usar el body
		pathID := r.URL.Path[len("/analyzers/"):]
		analyzerID, err := strconv.Atoi(pathID)
		if err != nil {
			http.Error(w, "Invalid analyzer ID", http.StatusBadRequest)
			return
		}

		// Convertir a int32
		id32 := int32(analyzerID)

		// Verificar si el Analyzer con el ID especificado existe en la base de datos
		exists, err := queries.AnalyzerExists(r.Context(), id32)
		if err != nil {
			log.Printf("Error checking if analyzer exists: %v", err)
			http.Error(w, "Failed to verify analyzer existence", http.StatusInternalServerError)
			return
		}

		if !exists {
			http.Error(w, "Analyzer not found", http.StatusNotFound)
			return
		}

		// Ejecutar la consulta de eliminación en la base de datos
		err = queries.DeleteAnalyzer(r.Context(), id32)
		if err != nil {
			log.Printf("Error deleting analyzer: %v", err)
			http.Error(w, "Failed to delete analyzer", http.StatusInternalServerError)
			return
		}

		// Responder con éxito
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Analyzer deleted successfully")
	}
}

// Handler para actualizar un Analyzer (PATCH /analyzers/:id)
func UpdateAnalyzer(queries *db.Queries, pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pathId := r.URL.Path[len("/analyzers/"):]
		analyzerID, err := strconv.Atoi(pathId)
		if err != nil {
			http.Error(w, "Invalid analyzer ID", http.StatusBadRequest)
			return
		}

		// Crear un mapa para almacenar los campos a actualizar
		var req map[string]any

		defer r.Body.Close() // Cerrar el cuerpo de la solicitud al finalizar

		// Decodificar el JSON en el mapa `req`
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decodificando JSON: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Construir la consulta de actualización dinámica
		setClauses := []string{}
		args := []any{analyzerID} // El primer argumento es el `analyzerID`
		argPosition := 2          // Comenzamos en 2 porque el primer argumento es el `analyzerID`

		for field, value := range req {
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argPosition))
			args = append(args, value)
			argPosition++
		}

		// Si no hay campos para actualizar, devolver un error
		if len(setClauses) == 0 {
			http.Error(w, "No fields to update", http.StatusBadRequest)
			return
		}

		// Construir la consulta SQL
		query := fmt.Sprintf("UPDATE analyzers SET %s WHERE analyzer_id = $1", strings.Join(setClauses, ", "))

		// Ejecutar la consulta en la base de datos
		result, err := pool.Exec(r.Context(), query, args...)
		if err != nil {
			log.Printf("Error updating analyzer: %v", err)
			http.Error(w, "Failed to update analyzer", http.StatusInternalServerError)
			return
		}

		// Verificar si se afectaron filas
		rowsAffected := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Analyzer not found", http.StatusNotFound)
			return
		}

		// Responder con éxito
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Analyzer updated successfully")
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
		itemCode, err := services.GenerateInventaryCode(req.ItemTypeName, req.BrandName, req.ModelName, req.Item.ItemAdquisitionDate, itemID)
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
		response := map[string]any{
			"item_id":     itemID,
			"analyzer_id": analyzerID,
			"item_code":   itemCode,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
