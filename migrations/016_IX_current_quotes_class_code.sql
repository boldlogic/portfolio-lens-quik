-- +goose Up
CREATE NONCLUSTERED INDEX IX_current_quotes_class_code ON quik.current_quotes (class_code) INCLUDE (
    sec_code,
    full_name,
    short_name,
    currency
)
;
 
