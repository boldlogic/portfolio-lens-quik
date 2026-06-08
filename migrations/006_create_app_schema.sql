-- +goose Up

CREATE SCHEMA app;

CREATE TYPE app.client_code_list AS TABLE (
    client_code VARCHAR(12) NOT NULL,
    PRIMARY KEY (client_code)
);
