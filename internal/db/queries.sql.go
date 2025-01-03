// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const analyzerExists = `-- name: AnalyzerExists :one
SELECT EXISTS(SELECT 1 FROM analyzers WHERE analyzer_id = $1)
`

func (q *Queries) AnalyzerExists(ctx context.Context, analyzerID int32) (bool, error) {
	row := q.db.QueryRow(ctx, analyzerExists, analyzerID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createAnalyzer = `-- name: CreateAnalyzer :one
INSERT INTO analyzers (
    item_id,
    brand_id,
    model_id,
    analyzer_state_id,
    analyzer_serialnumber,
    analyzer_pollutant
) VALUES (
    $1,           -- item_id obtenido del primer INSERT en items
    $2,           -- brand_id
    $3,           -- model_id
    $4,           -- analyzer_state_id
    $5,           -- analyzer_serialnumber
    $6           -- analyzer_pollutant
) RETURNING analyzer_id
`

type CreateAnalyzerParams struct {
	ItemID               int32  `json:"item_id"`
	BrandID              int32  `json:"brand_id"`
	ModelID              int32  `json:"model_id"`
	AnalyzerStateID      int32  `json:"analyzer_state_id"`
	AnalyzerSerialnumber string `json:"analyzer_serialnumber"`
	AnalyzerPollutant    string `json:"analyzer_pollutant"`
}

func (q *Queries) CreateAnalyzer(ctx context.Context, arg CreateAnalyzerParams) (int32, error) {
	row := q.db.QueryRow(ctx, createAnalyzer,
		arg.ItemID,
		arg.BrandID,
		arg.ModelID,
		arg.AnalyzerStateID,
		arg.AnalyzerSerialnumber,
		arg.AnalyzerPollutant,
	)
	var analyzer_id int32
	err := row.Scan(&analyzer_id)
	return analyzer_id, err
}

const createItem = `-- name: CreateItem :one
INSERT INTO items (
    item_type_id,
    item_code,
    item_name,
    item_description,
    item_adquisition_date,
    created_at
) VALUES (
    $1,         -- item_type_id
    $2,         -- item_code (generado en el backend)
    $3,         -- item_name
    $4,         -- item_description
    $5,         -- item_adquisition_date
    DEFAULT     -- created_at, usa la marca de tiempo actual
) RETURNING item_id
`

type CreateItemParams struct {
	ItemTypeID          int32       `json:"item_type_id"`
	ItemCode            string      `json:"item_code"`
	ItemName            string      `json:"item_name"`
	ItemDescription     pgtype.Text `json:"item_description"`
	ItemAdquisitionDate pgtype.Date `json:"item_adquisition_date"`
}

func (q *Queries) CreateItem(ctx context.Context, arg CreateItemParams) (int32, error) {
	row := q.db.QueryRow(ctx, createItem,
		arg.ItemTypeID,
		arg.ItemCode,
		arg.ItemName,
		arg.ItemDescription,
		arg.ItemAdquisitionDate,
	)
	var item_id int32
	err := row.Scan(&item_id)
	return item_id, err
}

const deleteAnalyzer = `-- name: DeleteAnalyzer :exec
DELETE FROM analyzers WHERE analyzer_id = $1
`

func (q *Queries) DeleteAnalyzer(ctx context.Context, analyzerID int32) error {
	_, err := q.db.Exec(ctx, deleteAnalyzer, analyzerID)
	return err
}

const getAnalyzer = `-- name: GetAnalyzer :one
SELECT analyzer_id, item_id, brand_id, model_id, analyzer_state_id, analyzer_serialnumber, analyzer_pollutant, analyzer_last_calibration, analyzer_last_maintenance FROM analyzers WHERE analyzer_id = $1
`

func (q *Queries) GetAnalyzer(ctx context.Context, analyzerID int32) (Analyzer, error) {
	row := q.db.QueryRow(ctx, getAnalyzer, analyzerID)
	var i Analyzer
	err := row.Scan(
		&i.AnalyzerID,
		&i.ItemID,
		&i.BrandID,
		&i.ModelID,
		&i.AnalyzerStateID,
		&i.AnalyzerSerialnumber,
		&i.AnalyzerPollutant,
		&i.AnalyzerLastCalibration,
		&i.AnalyzerLastMaintenance,
	)
	return i, err
}

const getAnalyzers = `-- name: GetAnalyzers :many
SELECT analyzer_id, item_id, brand_id, model_id, analyzer_state_id, analyzer_serialnumber, analyzer_pollutant, analyzer_last_calibration, analyzer_last_maintenance FROM analyzers
`

func (q *Queries) GetAnalyzers(ctx context.Context) ([]Analyzer, error) {
	rows, err := q.db.Query(ctx, getAnalyzers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Analyzer
	for rows.Next() {
		var i Analyzer
		if err := rows.Scan(
			&i.AnalyzerID,
			&i.ItemID,
			&i.BrandID,
			&i.ModelID,
			&i.AnalyzerStateID,
			&i.AnalyzerSerialnumber,
			&i.AnalyzerPollutant,
			&i.AnalyzerLastCalibration,
			&i.AnalyzerLastMaintenance,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBrandId = `-- name: GetBrandId :one
SELECT brand_id FROM brands WHERE brand_name = $1
`

func (q *Queries) GetBrandId(ctx context.Context, brandName string) (int32, error) {
	row := q.db.QueryRow(ctx, getBrandId, brandName)
	var brand_id int32
	err := row.Scan(&brand_id)
	return brand_id, err
}

const getItemTypeId = `-- name: GetItemTypeId :one
SELECT item_type_id FROM item_types WHERE type_name = $1
`

func (q *Queries) GetItemTypeId(ctx context.Context, typeName string) (int32, error) {
	row := q.db.QueryRow(ctx, getItemTypeId, typeName)
	var item_type_id int32
	err := row.Scan(&item_type_id)
	return item_type_id, err
}

const getItems = `-- name: GetItems :many
SELECT item_id, item_type_id, item_code, item_name, item_description, item_adquisition_date, created_at FROM items
`

func (q *Queries) GetItems(ctx context.Context) ([]Item, error) {
	rows, err := q.db.Query(ctx, getItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(
			&i.ItemID,
			&i.ItemTypeID,
			&i.ItemCode,
			&i.ItemName,
			&i.ItemDescription,
			&i.ItemAdquisitionDate,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getModelId = `-- name: GetModelId :one
SELECT model_id FROM models WHERE brand_id = $1 AND model_name = $2
`

type GetModelIdParams struct {
	BrandID   int32  `json:"brand_id"`
	ModelName string `json:"model_name"`
}

func (q *Queries) GetModelId(ctx context.Context, arg GetModelIdParams) (int32, error) {
	row := q.db.QueryRow(ctx, getModelId, arg.BrandID, arg.ModelName)
	var model_id int32
	err := row.Scan(&model_id)
	return model_id, err
}

const getStations = `-- name: GetStations :many
SELECT station_id, station_name, station_image_url, station_latitude, station_longitude, station_address, station_description, operational_since FROM stations
`

func (q *Queries) GetStations(ctx context.Context) ([]Station, error) {
	rows, err := q.db.Query(ctx, getStations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Station
	for rows.Next() {
		var i Station
		if err := rows.Scan(
			&i.StationID,
			&i.StationName,
			&i.StationImageUrl,
			&i.StationLatitude,
			&i.StationLongitude,
			&i.StationAddress,
			&i.StationDescription,
			&i.OperationalSince,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAnalyzer = `-- name: UpdateAnalyzer :exec
UPDATE analyzers SET
    item_id = $2,
    brand_id = $3,
    model_id = $4,
    analyzer_state_id = $5,
    analyzer_serialnumber = $6,
    analyzer_pollutant = $7,
    analyzer_last_calibration = $8,
    analyzer_last_maintenance = $9
WHERE analyzer_id = $1
`

type UpdateAnalyzerParams struct {
	AnalyzerID              int32       `json:"analyzer_id"`
	ItemID                  int32       `json:"item_id"`
	BrandID                 int32       `json:"brand_id"`
	ModelID                 int32       `json:"model_id"`
	AnalyzerStateID         int32       `json:"analyzer_state_id"`
	AnalyzerSerialnumber    string      `json:"analyzer_serialnumber"`
	AnalyzerPollutant       string      `json:"analyzer_pollutant"`
	AnalyzerLastCalibration pgtype.Date `json:"analyzer_last_calibration"`
	AnalyzerLastMaintenance pgtype.Date `json:"analyzer_last_maintenance"`
}

func (q *Queries) UpdateAnalyzer(ctx context.Context, arg UpdateAnalyzerParams) error {
	_, err := q.db.Exec(ctx, updateAnalyzer,
		arg.AnalyzerID,
		arg.ItemID,
		arg.BrandID,
		arg.ModelID,
		arg.AnalyzerStateID,
		arg.AnalyzerSerialnumber,
		arg.AnalyzerPollutant,
		arg.AnalyzerLastCalibration,
		arg.AnalyzerLastMaintenance,
	)
	return err
}

const updateInventaryCode = `-- name: UpdateInventaryCode :one
UPDATE items SET item_code = $1 WHERE item_id = $2 RETURNING item_code
`

type UpdateInventaryCodeParams struct {
	ItemCode string `json:"item_code"`
	ItemID   int32  `json:"item_id"`
}

func (q *Queries) UpdateInventaryCode(ctx context.Context, arg UpdateInventaryCodeParams) (string, error) {
	row := q.db.QueryRow(ctx, updateInventaryCode, arg.ItemCode, arg.ItemID)
	var item_code string
	err := row.Scan(&item_code)
	return item_code, err
}
