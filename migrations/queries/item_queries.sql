-- name: GetItems :many
SELECT item_id, model_id, item_description, item_serial_number, created_at
FROM items;


-- name: GetStations :many
SELECT station_id, station_name, station_image_url, operational_since, station_latitude, station_longitude, station_address
FROM stations;
