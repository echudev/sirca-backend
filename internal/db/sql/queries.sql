-- name: GetItems :many
SELECT item_id, model_id, item_description, item_serial_number, created_at
FROM items;
