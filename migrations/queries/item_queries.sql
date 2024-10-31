-- name: CreateItem :one
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
) RETURNING item_id;



-- name: GetItems :many
SELECT * FROM items;

-- name: UpdateInventaryCode :one
UPDATE items SET item_code = $1 WHERE item_id = $2 RETURNING item_code;

-- name: GetStations :many
SELECT *FROM stations;

-- name: GetAnalyzers :many
SELECT * FROM analyzers;

-- name: GetItemTypeId :one
SELECT item_type_id FROM item_types WHERE type_name = $1;

-- name: GetBrandId :one
SELECT brand_id FROM brands WHERE brand_name = $1;

-- name: GetModelId :one
SELECT model_id FROM models WHERE brand_id = $1 AND model_name = $2;

-- name: CreateAnalyzer :one
INSERT INTO analyzers (
    item_id,
    brand_id,
    model_id,
    analyzer_state_id,
    analyzer_serialnumber,
    analyzer_pollutant,
    analyzer_last_calibration,
    analyzer_last_maintenance
) VALUES (
    $1,           -- item_id obtenido del primer INSERT en items
    $2,           -- brand_id
    $3,           -- model_id
    $4,           -- analyzer_state_id
    $5,           -- analyzer_serialnumber
    $6,           -- analyzer_pollutant
    $7,           -- analyzer_last_calibration
    $8            -- analyzer_last_maintenance
) RETURNING analyzer_id ;
