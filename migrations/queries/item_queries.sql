-- name: SaveItem :one
WITH type_lookup AS (
    SELECT item_type_id FROM item_types WHERE type_name = $$type_name$$
)
INSERT INTO items (
    item_type_id,
    item_name,
    item_description
)
VALUES (
    (SELECT item_type_id FROM type_lookup),
    $$item_name$$,
    $$item_description$$
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
