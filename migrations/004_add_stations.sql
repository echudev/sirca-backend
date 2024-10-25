-- +goose Up

-- Inserta datos de estaciones en la tabla stations
INSERT INTO stations (station_name, station_image_url, station_description, operational_since)
VALUES
('Centenario',
'https://buenosaires.gob.ar/sites/default/files/styles/full_width/public/media/image/2015/12/18/507a431a07df170da69ff82797733fa33ba62cc7.jpg',
'',
'2005-06-01'),
('La Boca',
'https://buenosaires.gob.ar/sites/default/files/styles/full_width/public/media/image/2015/12/18/6a16045fe7f267b358b87d9505f0787c9b2ad8aa.jpg',
'',
'2009-01-01'),
('Córdoba',
'https://buenosaires.gob.ar/sites/default/files/styles/full_width/public/media/image/2015/12/18/ad83eca8d878eec0ba6b1af4c2ff33c900235c36.jpg',
'',
'2009-05-01'),
('CIFA',
'',
'',
'2016-01-01');

-- +goose Down

-- Elimina datos de estaciones en la tabla stations
DELETE FROM stations WHERE station_name IN ('Centenario', 'La Boca', 'Córdoba', 'CIFA');
