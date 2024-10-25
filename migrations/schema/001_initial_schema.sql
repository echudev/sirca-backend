-- +goose Up

-- Crear tabla de referencia para item_state
CREATE TABLE IF NOT EXISTS item_states (
    item_state_id SERIAL PRIMARY KEY,
    state_name VARCHAR(20) NOT NULL UNIQUE
);


-- Crear tabla de referencia para pollutant_type
CREATE TABLE IF NOT EXISTS pollutants (
    pollutant_id SERIAL PRIMARY KEY,
    pollutant_name VARCHAR(40) NOT NULL UNIQUE
);


-- Crear tabla de referencia para part_state
CREATE TABLE IF NOT EXISTS part_states (
    part_state_id SERIAL PRIMARY KEY,
    part_state_name VARCHAR(20) NOT NULL UNIQUE
);


-- Crear tabla de referencia para gas_type
CREATE TABLE IF NOT EXISTS gas_types (
    gas_type_id SERIAL PRIMARY KEY,
    gas_type_name VARCHAR(40) NOT NULL UNIQUE
);


-- Crear tabla de referencia para cylinder_size
CREATE TABLE IF NOT EXISTS cylinder_sizes (
    cylinder_size_id SERIAL PRIMARY KEY,
    size_name VARCHAR(20) NOT NULL UNIQUE
);


-- Crear tabla de referencia para concentration_unit
CREATE TABLE IF NOT EXISTS concentration_units (
    concentration_unit_id SERIAL PRIMARY KEY,
    unit_name VARCHAR(10) NOT NULL UNIQUE
);


-- Crear brands table
CREATE TABLE IF NOT EXISTS brands (
    brand_id SERIAL PRIMARY KEY,
    brand_name VARCHAR(40) UNIQUE NOT NULL
);

-- Crear models table
CREATE TABLE IF NOT EXISTS models (
    model_id SERIAL PRIMARY KEY,
    brand_id INT NOT NULL,
    model_name VARCHAR(40) NOT NULL,
    UNIQUE (brand_id, model_name),
    FOREIGN KEY (brand_id) REFERENCES brands(brand_id) ON DELETE CASCADE
);

-- Crear items table
CREATE TABLE IF NOT EXISTS items (
    item_id SERIAL PRIMARY KEY,
    model_id INT NOT NULL,
    item_description TEXT NOT NULL,
    item_serial_number VARCHAR(100) NOT NULL UNIQUE, -- Restricción ajustada
    item_image_url TEXT,
    item_supplier TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (model_id) REFERENCES models(model_id) ON DELETE CASCADE
);

-- Crear analyzers table
CREATE TABLE IF NOT EXISTS analyzers (
    analyzer_id SERIAL PRIMARY KEY,
    item_id INT NOT NULL,
    analyzer_state_id INT NOT NULL, -- Llave foránea a item_states
    analyzer_pollutant_id INT NOT NULL, -- Llave foránea a pollutants
    analyzer_last_calibration DATE CHECK (analyzer_last_calibration <= CURRENT_DATE), -- Asegurar fechas pasadas
    analyzer_last_maintenance DATE CHECK (analyzer_last_maintenance <= CURRENT_DATE),
    FOREIGN KEY (item_id) REFERENCES items(item_id) ON DELETE CASCADE,
    FOREIGN KEY (analyzer_state_id) REFERENCES item_states(item_state_id),
    FOREIGN KEY (analyzer_pollutant_id) REFERENCES pollutants(pollutant_id)
);

-- Crear parts table
CREATE TABLE IF NOT EXISTS parts (
    part_id SERIAL PRIMARY KEY,
    item_id INT NOT NULL,
    part_state_id INT NOT NULL, -- Llave foránea a part_states
    FOREIGN KEY (item_id) REFERENCES items(item_id) ON DELETE CASCADE,
    FOREIGN KEY (part_state_id) REFERENCES part_states(part_state_id)
);

-- Crear items_parts table
CREATE TABLE IF NOT EXISTS items_parts (
    item_id INT NOT NULL,
    part_id INT NOT NULL,
    PRIMARY KEY (item_id, part_id),
    FOREIGN KEY (item_id) REFERENCES items(item_id) ON DELETE CASCADE,
    FOREIGN KEY (part_id) REFERENCES parts(part_id) ON DELETE CASCADE
);

-- Crear cylinders table
CREATE TABLE IF NOT EXISTS cylinders (
    cylinder_id SERIAL PRIMARY KEY,
    item_id INT NOT NULL,
    cylinder_gas_type_id INT NOT NULL, -- Llave foránea a gas_types
    cylinder_size_id INT NOT NULL, -- Llave foránea a cylinder_sizes
    cylinder_unit_id INT NOT NULL, -- Llave foránea a concentration_units
    cylinder_concentration DECIMAL(10, 2) CHECK (cylinder_concentration > 0), -- Asegurar concentraciones válidas
    cylinder_expiration_date DATE CHECK (cylinder_expiration_date >= CURRENT_DATE), -- Asegurar fechas futuras
    FOREIGN KEY (item_id) REFERENCES items(item_id) ON DELETE CASCADE,
    FOREIGN KEY (cylinder_gas_type_id) REFERENCES gas_types(gas_type_id),
    FOREIGN KEY (cylinder_size_id) REFERENCES cylinder_sizes(cylinder_size_id),
    FOREIGN KEY (cylinder_unit_id) REFERENCES concentration_units(concentration_unit_id)
);

-- Crear stations table
CREATE TABLE IF NOT EXISTS stations (
    station_id SERIAL PRIMARY KEY,
    station_name VARCHAR(100) NOT NULL, -- Tamaño restringido a lo necesario
    station_image_url TEXT,
    station_latitude DECIMAL(9,6),
    station_longitude DECIMAL(9,6),
    station_address TEXT,
    station_description TEXT,
    operational_since DATE CHECK (operational_since <= CURRENT_DATE) -- Asegurar que sea una fecha pasada
);

-- Crear inventory table
CREATE TABLE IF NOT EXISTS inventory (
    item_id INT NOT NULL,
    station_id INT NOT NULL,
    quantity INT NOT NULL CHECK (quantity >= 0), -- Asegurar que sea positivo
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

-- Crear índices adicionales para mejorar el rendimiento de consultas comunes
CREATE INDEX IF NOT EXISTS idx_models_brand_id ON models(brand_id);
CREATE INDEX IF NOT EXISTS idx_items_model_id ON items(model_id);
CREATE INDEX IF NOT EXISTS idx_analyzers_item_id ON analyzers(item_id);
CREATE INDEX IF NOT EXISTS idx_parts_item_id ON parts(item_id);
CREATE INDEX IF NOT EXISTS idx_inventory_item_id ON inventory(item_id);
CREATE INDEX IF NOT EXISTS idx_inventory_station_id ON inventory(station_id);
CREATE INDEX IF NOT EXISTS idx_traslados_item_station ON traslados(item_id, station_id_origen, station_id_destino);

-- +goose Down
-- Drop tables
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS cylinders;
DROP TABLE IF EXISTS items_parts;
DROP TABLE IF EXISTS parts;
DROP TABLE IF EXISTS analyzers;
DROP TABLE IF EXISTS item_states;
DROP TABLE IF EXISTS pollutants;
DROP TABLE IF EXISTS part_states;
DROP TABLE IF EXISTS gas_types;
DROP TABLE IF EXISTS cylinder_sizes;
DROP TABLE IF EXISTS concentration_units;
DROP TABLE IF EXISTS traslados;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS models;
DROP TABLE IF EXISTS stations;
DROP TABLE IF EXISTS brands;
