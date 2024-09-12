-- name: CreateItem :one
INSERT INTO items (brand, model, description, serial_number) 
VALUES ($1, $2, $3, $4)
RETURNING id, brand, model, description, created_at;

-- name: GetItem :one
SELECT id, brand, model, created_at 
FROM items 
WHERE id = $1;

-- name: ListItems :many
SELECT id, brand, model, description, serial_number, created_at
FROM items;

-- name: UpdateItem :one
UPDATE items 
SET brand = $2, model = $3 
WHERE id = $1
RETURNING id, brand, model, description, created_at;

-- name: DeleteItem :exec
DELETE FROM items 
WHERE id = $1;
