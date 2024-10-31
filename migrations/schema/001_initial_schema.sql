-- +goose Up

-- Crear la tabla de item_types para definir los tipos de items
CREATE TABLE IF NOT EXISTS item_types (
    item_type_id SERIAL PRIMARY KEY,
    type_name VARCHAR(20) NOT NULL UNIQUE
);

-- Crear tabla de referencia para analyzer_states
CREATE TABLE IF NOT EXISTS analyzer_states (
    analyzer_state_id SERIAL PRIMARY KEY,
    state_name VARCHAR(20) NOT NULL UNIQUE
);

-- Crear tabla de referencia para part_states
CREATE TABLE IF NOT EXISTS part_states (
    part_state_id SERIAL PRIMARY KEY,
    part_state_name VARCHAR(20) NOT NULL UNIQUE
);

-- Crear tabla de referencia para brands
CREATE TABLE IF NOT EXISTS brands (
    brand_id SERIAL PRIMARY KEY,
    brand_name VARCHAR(40) UNIQUE NOT NULL
);

-- Crear tabla de referencia para models
CREATE TABLE IF NOT EXISTS models (
    model_id SERIAL PRIMARY KEY,
    brand_id INT NOT NULL,
    model_name VARCHAR(40) NOT NULL,
    UNIQUE (brand_id, model_name),
    FOREIGN KEY (brand_id) REFERENCES brands(brand_id) ON DELETE CASCADE
);

-- Crear tabla items (tabla central)
CREATE TABLE IF NOT EXISTS items (
    item_id SERIAL PRIMARY KEY,
    item_type_id INT NOT NULL,
    item_code VARCHAR(40) NOT NULL UNIQUE,
    item_name VARCHAR(100) NOT NULL,
    item_description TEXT,
    item_adquisition_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (item_type_id) REFERENCES item_types(item_type_id) ON DELETE RESTRICT
);

-- Create the analyzers table
CREATE TABLE IF NOT EXISTS analyzers (
    analyzer_id SERIAL PRIMARY KEY,
    item_id INT NOT NULL UNIQUE,
    brand_id INT NOT NULL,
    model_id INT NOT NULL,
    analyzer_state_id INT NOT NULL,
    analyzer_serialnumber VARCHAR(40) NOT NULL UNIQUE,
    analyzer_pollutant VARCHAR(40) NOT NULL,
    analyzer_last_calibration DATE,
    analyzer_last_maintenance DATE,
    FOREIGN KEY (item_id) REFERENCES items(item_id) ON DELETE CASCADE,
    FOREIGN KEY (brand_id) REFERENCES brands(brand_id),
    FOREIGN KEY (model_id) REFERENCES models(model_id) ON DELETE CASCADE,
    FOREIGN KEY (analyzer_state_id) REFERENCES analyzer_states(analyzer_state_id) ON DELETE CASCADE
);

-- Crear tabla parts
CREATE TABLE IF NOT EXISTS parts (
    part_id SERIAL PRIMARY KEY,
    item_id INT NOT NULL,
    part_number VARCHAR(30) NOT NULL,
    part_serialnumber VARCHAR(40) NOT NULL,
    part_state_id INT NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items(item_id) ON DELETE CASCADE,
    FOREIGN KEY (part_state_id) REFERENCES part_states(part_state_id) ON DELETE CASCADE
);

-- Crear tabla items_parts para asociar items y parts
CREATE TABLE IF NOT EXISTS items_parts (
    item_id INT NOT NULL,
    part_id INT NOT NULL,
    PRIMARY KEY (item_id, part_id),
    FOREIGN KEY (item_id) REFERENCES items(item_id) ON DELETE CASCADE,
    FOREIGN KEY (part_id) REFERENCES parts(part_id) ON DELETE CASCADE
);

-- Crear tabla cylinders
CREATE TABLE IF NOT EXISTS cylinders (
    cylinder_id SERIAL PRIMARY KEY,
    item_id INT NOT NULL,
    cylinder_number VARCHAR(30) NOT NULL,
    cylinder_concentration DECIMAL(10, 2),
    cylinder_connector VARCHAR(20) NOT NULL,
    cylinder_expiration_date DATE NOT NULL,
    cylinder_certificate TEXT,
    FOREIGN KEY (item_id) REFERENCES items(item_id) ON DELETE CASCADE
);

-- Crear tabla stations
CREATE TABLE IF NOT EXISTS stations (
    station_id SERIAL PRIMARY KEY,
    station_name VARCHAR(100) NOT NULL,
    station_image_url TEXT,
    station_latitude DECIMAL(9,6),
    station_longitude DECIMAL(9,6),
    station_address TEXT,
    station_description TEXT,
    operational_since DATE
);

-- Crear tabla inventory
CREATE TABLE IF NOT EXISTS inventory (
    item_id INT NOT NULL,
    station_id INT NOT NULL,
    quantity INT NOT NULL CHECK (quantity >= 0),
    last_update TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(40) NOT NULL,
    PRIMARY KEY (item_id, station_id),
    FOREIGN KEY (item_id) REFERENCES items(item_id) ON DELETE CASCADE,
    FOREIGN KEY (station_id) REFERENCES stations(station_id) ON DELETE CASCADE
);

-- Crear tabla traslados para registrar el movimiento de items entre estaciones
CREATE TABLE IF NOT EXISTS traslados (
    traslado_id SERIAL PRIMARY KEY,
    item_id INT NOT NULL,
    station_id_origen INT NOT NULL,
    station_id_destino INT NOT NULL,
    cantidad INT NOT NULL CHECK (cantidad > 0),
    fecha_traslado TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (item_id) REFERENCES items(item_id) ON DELETE CASCADE,
    FOREIGN KEY (station_id_origen) REFERENCES stations(station_id) ON DELETE CASCADE,
    FOREIGN KEY (station_id_destino) REFERENCES stations(station_id) ON DELETE CASCADE,
    CHECK (station_id_origen <> station_id_destino)
);

-- Insertar valores iniciales en item_types
INSERT INTO item_types (type_name)
VALUES
    ('analyzer'),
    ('cylinder'),
    ('part')
ON CONFLICT DO NOTHING;

-- Insertar valores iniciales en analyzer_states
INSERT INTO analyzer_states (state_name)
VALUES
    ('active'),
    ('inactive'),
    ('maintenance')
ON CONFLICT DO NOTHING;

-- Insertar valores iniciales en part_states
INSERT INTO part_states (part_state_name)
VALUES
    ('new'),
    ('used'),
    ('broken'),
    ('obsolete')
ON CONFLICT DO NOTHING;

-- Crear índices adicionales para mejorar el rendimiento de consultas comunes
CREATE INDEX IF NOT EXISTS idx_models_brand_id ON models(brand_id);
CREATE INDEX IF NOT EXISTS idx_analyzers_item_id ON analyzers(item_id);
CREATE INDEX IF NOT EXISTS idx_parts_item_id ON parts(item_id);
CREATE INDEX IF NOT EXISTS idx_inventory_item_id ON inventory(item_id);
CREATE INDEX IF NOT EXISTS idx_inventory_station_id ON inventory(station_id);
CREATE INDEX IF NOT EXISTS idx_traslados_item_station ON traslados(item_id, station_id_origen, station_id_destino);

-- +goose Down

DELETE FROM inventory;
DELETE FROM traslados;
DELETE FROM cylinders;
DELETE FROM items_parts;
DELETE FROM parts;
DELETE FROM analyzers;  -- Eliminar antes de borrar models
DELETE FROM models;
DELETE FROM brands;
DELETE FROM items;
DELETE FROM analyzer_states;
DELETE FROM part_states;
DELETE FROM stations;
DELETE FROM item_types;

-- Eliminar las tablas en el orden inverso para evitar errores de dependencia
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS traslados;
DROP TABLE IF EXISTS cylinders;
DROP TABLE IF EXISTS items_parts;
DROP TABLE IF EXISTS parts;
DROP TABLE IF EXISTS analyzers;
DROP TABLE IF EXISTS analyzer_states;
DROP TABLE IF EXISTS part_states;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS models;
DROP TABLE IF EXISTS brands;
DROP TABLE IF EXISTS stations;
DROP TABLE IF EXISTS item_types;
