CREATE DATABASE sircadb;

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    brand TEXT NOT NULL CHECK (char_length(brand) > 0),
    model TEXT NOT NULL CHECK (char_length(model) > 0),
    description TEXT NOT NULL,
    serial_number TEXT NOT NULL CHECK (char_length(serial_number) > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (serial_number)
);
