CREATE DATABASE sircadb;

-- Tabla base: items (superclase)
CREATE TABLE items (
    item_id SERIAL PRIMARY KEY,
    item_brand VARCHAR(40) NOT NULL CHECK (char_length(brand) > 0),
    item_model VARCHAR(40) NOT NULL CHECK (char_length(model) > 0),
    item_description TEXT NOT NULL,
    item_serial_number VARCHAR(255) NOT NULL CHECK (char_length(serial_number) > 0),
    item_image_url TEXT,
    item_supplier TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- Registro de creación
    UNIQUE (serial_number)
);

-- Tabla: analyzers (hereda de items)
CREATE TABLE analyzers (
    analyzer_id SERIAL PRIMARY KEY,
    item_id INT REFERENCES items(item_id) ON DELETE CASCADE,
    analyzer_state ENUM('active', 'inactive', 'maintenance') NOT NULL,
    analyzer_pollutant ENUM('particulate', 'ozone', 'nitrogen', 'carbon monoxide', 'sulfur dioxide', 'hydrogen sulfide') NOT NULL,
    analyzer_last_calibration DATE,
    analuzer_last_maintenance DATE
);

-- Tabla: parts (hereda de items)
CREATE TABLE parts (
    part_id SERIAL PRIMARY KEY,
    item_id INT REFERENCES items(item_id) ON DELETE CASCADE,
    part_state ENUM('new', 'used', 'broken', 'obsolete') NOT NULL
);

-- Tabla intermedia de muchos a muchos entre items y parts (items_parts)
CREATE TABLE items_parts (
    item_id INT REFERENCES items(item_id) ON DELETE CASCADE,  -- El ítem al que pertenece la parte
    part_id INT REFERENCES parts(part_id) ON DELETE CASCADE,  -- La parte que pertenece al ítem
    PRIMARY KEY (item_id, part_id)
);

-- Table: cylinders
CREATE TABLE cylinders (
    cylinder_id SERIAL PRIMARY KEY,
    item_id INT REFERENCES items(item_id) ON DELETE CASCADE,
    cylinder_gas ENUM('nitrogen', 'oxygen', 'argon', 'carbon dioxide', 'hydrogen', 'methane', 'water', 'other') NOT NULL,
    cylinder_size ENUM('small', 'medium', 'large') NOT NULL,
    cyliinder_unit ENUM('ppm','ppb','ppt','mg/m3','g/m3') NOT NULL,
    cylinder_concentration DECIMAL(10, 2),
    cylinder_expiration_date DATE
);

-- Table: stations
CREATE TABLE stations (
    station_id SERIAL PRIMARY KEY,
    station_name VARCHAR(255) NOT NULL,
    station_image_url TEXT,
    station_description TEXT,
    operational_since DATE
);

-- Table: inventory
CREATE TABLE inventory (
    item_id INT REFERENCES items(item_id) ON DELETE CASCADE,
    station_id INT REFERENCES stations(station_id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity >= 0),
    PRIMARY KEY (item_id, station_id)
);