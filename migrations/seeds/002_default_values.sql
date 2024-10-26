-- +goose Up

-- Insertar valores iniciales en item_states
INSERT INTO item_states (state_name)
VALUES
    ('active'),
    ('inactive'),
    ('maintenance')
ON CONFLICT DO NOTHING;

-- Insertar valores iniciales en pollutants
INSERT INTO pollutants (pollutant_name)
VALUES
    ('particulate'),
    ('ozone'),
    ('nitrogen oxides'),
    ('carbon monoxide'),
    ('sulfur dioxide'),
    ('hydrogen sulfide')
ON CONFLICT DO NOTHING;

-- Insertar valores iniciales en part_states
INSERT INTO part_states (part_state_name)
VALUES
    ('new'),
    ('used'),
    ('broken'),
    ('obsolete')
ON CONFLICT DO NOTHING;

-- Insertar valores iniciales en gas_types
INSERT INTO cylinder_gas (gas_type_name)
VALUES
    ('CO'),
    ('NO'),
    ('SO2'),
    ('H2S'),
    ('NO2'),
ON CONFLICT DO NOTHING;

-- Insertar valores iniciales en cylinder_sizes
INSERT INTO cylinder_volumes (cylinder_volume)
VALUES
    ('83,4 CF'),
    ('140 CF'),
ON CONFLICT DO NOTHING;

-- Insertar valores iniciales en cylinder_connectors
INSERT INTO cylinder_connectors (cylinder_connector)
VALUES
    ('CGA-330'),
    ('CGA-350'),
    ('CGA-660')
ON CONFLICT DO NOTHING;


-- Insertar valores iniciales en concentration_units
INSERT INTO concentration_units (unit_name)
VALUES
    ('ppm'),
    ('ppb'),
    ('ppt'),
    ('mg/m3'),
    ('g/m3')
ON CONFLICT DO NOTHING;

-- +goose Down

-- Eliminar datos de referencia (opcional)
DELETE FROM item_states WHERE state_name IN ('active', 'inactive', 'maintenance');
DELETE FROM pollutants WHERE pollutant_name IN ('particulate', 'ozone', 'nitrogen oxides', 'carbon monoxide', 'sulfur dioxide', 'hydrogen sulfide');
DELETE FROM part_states WHERE part_state_name IN ('new', 'used', 'broken', 'obsolete');
DELETE FROM gas_types WHERE gas_type_name IN ('nitrogen', 'oxygen', 'argon', 'carbon dioxide', 'hydrogen', 'helium', 'methane', 'water');
DELETE FROM cylinder_sizes WHERE size_name IN ('small', 'medium', 'large');
DELETE FROM concentration_units WHERE unit_name IN ('ppm', 'ppb', 'ppt', 'mg/m3', 'g/m3');
