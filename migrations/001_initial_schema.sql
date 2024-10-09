-- +goose Up
-- Create the database (Note: Goose typically doesn't handle database creation)
-- You might need to create the database manually or use a separate script

-- Create ENUM types
CREATE TYPE item_state AS ENUM('active', 'inactive', 'maintenance');
CREATE TYPE pollutant_type AS ENUM('particulate', 'ozone', 'nitrogen oxides', 'carbon monoxide', 'sulfur dioxide', 'hydrogen sulfide');
CREATE TYPE part_state AS ENUM('new', 'used', 'broken', 'obsolete');
CREATE TYPE gas_type AS ENUM('nitrogen', 'oxygen', 'argon', 'carbon dioxide', 'hydrogen', 'methane', 'water', 'other');
CREATE TYPE cylinder_size AS ENUM('small', 'medium', 'large');
CREATE TYPE concentration_unit AS ENUM('ppm','ppb','ppt','mg/m3','g/m3');

-- Create brands table
CREATE TABLE IF NOT EXISTS brands (
    brand_id SERIAL PRIMARY KEY,
    brand_name VARCHAR(40) UNIQUE NOT NULL
);

-- Create models table
CREATE TABLE IF NOT EXISTS models (
    model_id SERIAL PRIMARY KEY,
    brand_id INT REFERENCES brands(brand_id),
    model_name VARCHAR(40) NOT NULL,
    UNIQUE (brand_id, model_name)
);

-- Create items table
CREATE TABLE IF NOT EXISTS items (
    item_id SERIAL PRIMARY KEY,
    model_id INT REFERENCES models(model_id),
    item_description TEXT NOT NULL,
    item_serial_number VARCHAR(255) NOT NULL UNIQUE,
    item_image_url TEXT,
    item_supplier TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create analyzers table
CREATE TABLE IF NOT EXISTS analyzers (
    analyzer_id SERIAL PRIMARY KEY,
    item_id INT REFERENCES items(item_id) ON DELETE CASCADE,
    analyzer_state item_state NOT NULL,
    analyzer_pollutant pollutant_type NOT NULL,
    analyzer_last_calibration DATE,
    analyzer_last_maintenance DATE
);

-- Create parts table
CREATE TABLE IF NOT EXISTS parts (
    part_id SERIAL PRIMARY KEY,
    item_id INT REFERENCES items(item_id) ON DELETE CASCADE,
    part_state part_state NOT NULL
);

-- Create items_parts table
CREATE TABLE IF NOT EXISTS items_parts (
    item_id INT REFERENCES items(item_id) ON DELETE CASCADE,
    part_id INT REFERENCES parts(part_id) ON DELETE CASCADE,
    PRIMARY KEY (item_id, part_id)
);

-- Create cylinders table
CREATE TABLE IF NOT EXISTS cylinders (
    cylinder_id SERIAL PRIMARY KEY,
    item_id INT REFERENCES items(item_id) ON DELETE CASCADE,
    cylinder_gas gas_type NOT NULL,
    cylinder_size cylinder_size NOT NULL,
    cylinder_unit concentration_unit NOT NULL,
    cylinder_concentration DECIMAL(10, 2),
    cylinder_expiration_date DATE
);

-- Create stations table
CREATE TABLE IF NOT EXISTS stations (
    station_id SERIAL PRIMARY KEY,
    station_name VARCHAR(255) NOT NULL,
    station_image_url TEXT,
    station_description TEXT,
    operational_since DATE
);

-- Create inventory table
CREATE TABLE IF NOT EXISTS inventory (
    item_id INT REFERENCES items(item_id) ON DELETE CASCADE,
    station_id INT REFERENCES stations(station_id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity >= 0),
    PRIMARY KEY (item_id, station_id)
);

-- +goose Down
-- Drop tables
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS stations;
DROP TABLE IF EXISTS cylinders;
DROP TABLE IF EXISTS items_parts;
DROP TABLE IF EXISTS parts;
DROP TABLE IF EXISTS analyzers;
DROP TABLE IF EXISTS items;

-- Drop ENUM types
DROP TYPE IF EXISTS concentration_unit;
DROP TYPE IF EXISTS cylinder_size;
DROP TYPE IF EXISTS gas_type;
DROP TYPE IF EXISTS part_state;
DROP TYPE IF EXISTS pollutant_type;
DROP TYPE IF EXISTS item_state;

-- Note: Dropping the database should be done manually if needed
