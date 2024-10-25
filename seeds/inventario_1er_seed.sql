-- Insertar marcas
INSERT INTO brands (brand_name) VALUES
('Ecotech'), ('Thermo'), ('Monitor Europe'), ('Teledyne'), ('Envea'), ('Environnement'), ('Otros');





-- Hacé de nuevo los seeds, sacaste los enums de la base de datos
-- primero tirate un goose down
-- despues goose up
-- preguntale a gpt como es el algoritmo para modificar base de datos paso a paso






-- Insertar modelos
INSERT INTO models (brand_id, model_name) VALUES
((SELECT brand_id FROM brands WHERE brand_name = 'Ecotech'), 'Serinus 30'),
((SELECT brand_id FROM brands WHERE brand_name = 'Ecotech'), 'Serinus 10'),
((SELECT brand_id FROM brands WHERE brand_name = 'Ecotech'), 'Model Z3'),
((SELECT brand_id FROM brands WHERE brand_name = 'Thermo'), 'Model X2'),
((SELECT brand_id FROM brands WHERE brand_name = 'Thermo'), 'Model W4');

-- Insertar items
INSERT INTO items (model_id, item_description, item_serial_number, item_image_url, item_supplier) VALUES
((SELECT model_id FROM models WHERE model_name = 'Serinus 30'), 'Analizador de monóxido de carbono', '22-0710', NULL, NULL),
((SELECT model_id FROM models WHERE model_name = 'Serinus 10'), 'Analizador de Ozono', '16-1099', NULL, NULL),
((SELECT model_id FROM models WHERE model_name = 'Model Z3'), 'Advanced air quality monitor', 'SN003', 'http://example.com/imageZ3.jpg', 'Supplier 3'),
((SELECT model_id FROM models WHERE model_name = 'Model X2'), 'Calibration gas cylinder', 'SN004', 'http://example.com/imageX2.jpg', 'Supplier 1'),
((SELECT model_id FROM models WHERE model_name = 'Model W4'), 'Replacement filter', 'SN005', 'http://example.com/imageW4.jpg', 'Supplier 4');

-- Seed data for parts table
INSERT INTO parts (item_id, part_state) VALUES
(4, 'new'),
(5, 'used');

-- Seed data for items_parts table
INSERT INTO items_parts (item_id, part_id) VALUES
(1, 5),
(2, 5),
(3, 4);

-- Seed data for cylinders table
INSERT INTO cylinders (item_id, cylinder_gas, cylinder_size, cylinder_unit, cylinder_concentration, cylinder_expiration_date) VALUES
(4, 'nitrogen', 'medium', 'ppm', 100.00, '2024-12-31');

-- Seed data for stations table
INSERT INTO stations (station_name, station_image_url, station_description, operational_since) VALUES
('Downtown Station', 'http://example.com/downtown.jpg', 'Main city center monitoring station', '2020-01-01'),
('Suburban Station', 'http://example.com/suburban.jpg', 'Residential area monitoring station', '2021-03-15'),
('Industrial Zone Station', 'http://example.com/industrial.jpg', 'Heavy industry area monitoring station', '2019-07-01');

-- Seed data for inventory table
INSERT INTO inventory (item_id, station_id, quantity) VALUES
(1, 1, 2),
(2, 1, 1),
(3, 2, 1),
(4, 3, 5),
(5, 2, 10);
