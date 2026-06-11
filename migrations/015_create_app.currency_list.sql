-- +goose Up

CREATE TYPE app.currency_list 
AS TABLE (
    iso_code SMALLINT NOT NULL,
    iso_char_code CHAR(3) NOT NULL,
    currency_name NVARCHAR(100) NULL,
    lat_name NVARCHAR(100) NULL,
    minor_units INT NULL,
    PRIMARY KEY CLUSTERED (iso_code)
);

GRANT EXECUTE ON TYPE::app.currency_list TO quik_currency_worker;