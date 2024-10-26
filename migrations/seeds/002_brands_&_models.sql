-- +goose Up

-- Insertar marcas en la tabla brands
INSERT INTO brands (brand_name) VALUES
('Thermo'),
('Ecotech'),
('Acoem'),
('Teledyne'),
('Horiba'),
('Siemens'),
('TSI'),
('MonitorEurope'),
('MonitorLabs'),
('Environnement'),
('ENVEA'),
('MetOne'),
('DAVIS')
ON CONFLICT DO NOTHING;


-- Insertar modelos en la tabla models

-- Thermo
INSERT INTO models (brand_id, model_name) VALUES
((SELECT brand_id FROM brands WHERE brand_name = 'Thermo'), 'F64C14'),
((SELECT brand_id FROM brands WHERE brand_name = 'Thermo'), '48C'),
((SELECT brand_id FROM brands WHERE brand_name = 'Thermo'), '42C'),
((SELECT brand_id FROM brands WHERE brand_name = 'Thermo'), '48i'),
((SELECT brand_id FROM brands WHERE brand_name = 'Thermo'), '42i'),
((SELECT brand_id FROM brands WHERE brand_name = 'Thermo'), '49i'),
((SELECT brand_id FROM brands WHERE brand_name = 'Thermo'), '5030i')
ON CONFLICT DO NOTHING;


-- Ecotech
INSERT INTO models (brand_id, model_name) VALUES
((SELECT brand_id FROM brands WHERE brand_name = 'Ecotech'), 'EC9830'),
((SELECT brand_id FROM brands WHERE brand_name = 'Ecotech'), 'EC9841'),
((SELECT brand_id FROM brands WHERE brand_name = 'Ecotech'), 'Serinus 10'),
((SELECT brand_id FROM brands WHERE brand_name = 'Ecotech'), 'Serinus 30'),
((SELECT brand_id FROM brands WHERE brand_name = 'Ecotech'), 'Serinus 40'),
((SELECT brand_id FROM brands WHERE brand_name = 'Ecotech'), 'Serinus 50'),
((SELECT brand_id FROM brands WHERE brand_name = 'Ecotech'), 'Serinus CAL3000'),
((SELECT brand_id FROM brands WHERE brand_name = 'Ecotech'), 'SpirantBam 1020')
ON CONFLICT DO NOTHING;

-- MonitorLabs
INSERT INTO models (brand_id, model_name) VALUES
((SELECT brand_id FROM brands WHERE brand_name = 'MonitorLabs'), 'ML9810B'),
((SELECT brand_id FROM brands WHERE brand_name = 'MonitorLabs'), 'ML9830B'),
((SELECT brand_id FROM brands WHERE brand_name = 'MonitorLabs'), 'ML9841B')
ON CONFLICT DO NOTHING;

-- MonitorEurope
INSERT INTO models (brand_id, model_name) VALUES
((SELECT brand_id FROM brands WHERE brand_name = 'MonitorEurope'), 'ME9810'),
((SELECT brand_id FROM brands WHERE brand_name = 'MonitorEurope'), 'ME9830'),
((SELECT brand_id FROM brands WHERE brand_name = 'MonitorEurope'), 'ME9841')
ON CONFLICT DO NOTHING;

-- Acoem
INSERT INTO models (brand_id, model_name) VALUES
((SELECT brand_id FROM brands WHERE brand_name = 'Acoem'), 'Serinus 10'),
((SELECT brand_id FROM brands WHERE brand_name = 'Acoem'), 'Serinus 30'),
((SELECT brand_id FROM brands WHERE brand_name = 'Acoem'), 'Serinus 40'),
((SELECT brand_id FROM brands WHERE brand_name = 'Acoem'), 'Serinus 50'),
((SELECT brand_id FROM brands WHERE brand_name = 'Acoem'), 'SpirantBam 1020')
ON CONFLICT DO NOTHING;

-- Teledyne
INSERT INTO models (brand_id, model_name) VALUES
((SELECT brand_id FROM brands WHERE brand_name = 'Teledyne'), 'T100'),
((SELECT brand_id FROM brands WHERE brand_name = 'Teledyne'), 'T100u'),
((SELECT brand_id FROM brands WHERE brand_name = 'Teledyne'), 'T200'),
((SELECT brand_id FROM brands WHERE brand_name = 'Teledyne'), 'T300'),
((SELECT brand_id FROM brands WHERE brand_name = 'Teledyne'), 'T400'),
((SELECT brand_id FROM brands WHERE brand_name = 'Teledyne'), 'T700')
ON CONFLICT DO NOTHING;

-- MetOne
INSERT INTO models (brand_id, model_name) VALUES
((SELECT brand_id FROM brands WHERE brand_name = 'MetOne'), 'BAM-1020'),
((SELECT brand_id FROM brands WHERE brand_name = 'MetOne'), 'ES-405')
ON CONFLICT DO NOTHING;

-- ENVEA
INSERT INTO models (brand_id, model_name) VALUES
((SELECT brand_id FROM brands WHERE brand_name = 'ENVEA'), 'CO12e'),
((SELECT brand_id FROM brands WHERE brand_name = 'ENVEA'), 'AF22e'),
((SELECT brand_id FROM brands WHERE brand_name = 'ENVEA'), 'AC32e'),
((SELECT brand_id FROM brands WHERE brand_name = 'ENVEA'), 'MP101M')
ON CONFLICT DO NOTHING;

-- Environnement
INSERT INTO models (brand_id, model_name) VALUES
((SELECT brand_id FROM brands WHERE brand_name = 'Environnement'), 'CO12m'),
((SELECT brand_id FROM brands WHERE brand_name = 'Environnement'), 'AF22m'),
((SELECT brand_id FROM brands WHERE brand_name = 'Environnement'), 'AC32m'),
((SELECT brand_id FROM brands WHERE brand_name = 'Environnement'), 'ZAG7001'),
((SELECT brand_id FROM brands WHERE brand_name = 'Environnement'), 'MGC101p')
ON CONFLICT DO NOTHING;

-- DAVIS
INSERT INTO models (brand_id, model_name) VALUES
((SELECT brand_id FROM brands WHERE brand_name = 'DAVIS'), 'VantagePro2')
ON CONFLICT DO NOTHING;


-- +goose Down

-- Eliminar los modelos de todas las marcas en caso de rollback
DELETE FROM models WHERE model_name IN (
    'F64C14', '48C', '42C', '48i', '42i', '49i', '5030i',
    'EC9830', 'EC9841', 'Serinus 10', 'Serinus 30', 'Serinus 40', 'Serinus 50', 'Serinus CAL3000', 'SpirantBam 1020',
    'ML9810B', 'ML9830B', 'ML9841B',
    'ME9810', 'ME9830', 'ME9841',
    'T100', 'T100u', 'T200', 'T300', 'T400', 'T700',
    'BAM-1020', 'ES-405',
    'CO12e', 'AF22e', 'AC32e', 'MP101M',
    'CO12m', 'AF22m', 'AC32m', 'ZAG7001', 'MGC101p',
    'VantagePro2'
);

-- Eliminar las marcas en caso de rollback
DELETE FROM brands WHERE brand_name IN (
    'Thermo', 'Ecotech', 'Acoem', 'Teledyne', 'Horiba', 'Siemens', 'TSI', 'MonitorEurope', 'MonitorLabs',
    'Environnement', 'ENVEA', 'MetOne', 'DAVIS'
);
