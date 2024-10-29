-- name: CreateItem :one
INSERT INTO items (
    item_type_id,
    item_code,
    item_name,
    item_description,
    created_at
) VALUES (
    $1,         -- item_type_id
    $2,         -- item_code (generado en el backend)
    $3,         -- item_name
    $4,         -- item_description
    DEFAULT     -- created_at, usa la marca de tiempo actual
) RETURNING item_id;



-- name: GetItems :many
SELECT
    i.item_id,
    i.item_name,
    i.item_description,
    i.created_at,
    it.type_name AS item_type
FROM
    items i
JOIN
    item_types it ON i.item_type_id = it.item_type_id;


-- name: GetStations :many
SELECT station_id, station_name, station_image_url, operational_since, station_latitude, station_longitude, station_address
FROM stations;


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
