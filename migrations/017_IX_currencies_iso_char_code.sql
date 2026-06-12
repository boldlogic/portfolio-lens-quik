-- +goose Up
CREATE UNIQUE NONCLUSTERED INDEX IX_currencies_iso_char_code ON ref.currencies (iso_char_code) INCLUDE (
    iso_code
)
;
 
