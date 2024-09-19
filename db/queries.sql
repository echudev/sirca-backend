-- name: CreateItem :one
INSERT INTO items (item_brand, item_model, item_description, item_serial_number)
VALUES ($1, $2, $3, $4)
RETURNING item_id, item_brand, item_model, item_description, created_at;

-- name: GetItem :one
SELECT item_id, item_brand, item_model, created_at 
FROM items 
WHERE item_id = $1;

-- name: ListItems :many
SELECT item_brand, item_model, item_description, item_serial_number, created_at
FROM items;

-- name: UpdateItem :one
UPDATE items 
SET item_brand = $2, item_model = $3 
WHERE item_id = $1
RETURNING item_id, item_brand, item_model, item_description, created_at;

-- name: DeleteItem :exec
DELETE FROM items 
WHERE item_id = $1;

-- name: ListAnalyzers :many
SELECT i.item_brand, i.item_model, i.item_serial_number, a.analyzer_state 
FROM analyzers a JOIN items i ON a.item_id = i.item_id;