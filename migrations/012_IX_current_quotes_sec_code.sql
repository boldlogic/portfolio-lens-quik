-- +goose Up
CREATE NONCLUSTERED INDEX IX_current_quotes_sec_code ON quik.current_quotes (sec_code) INCLUDE (
    quote_date,
    short_name,
    instrument_type,
    face_value,
    currency,
    base_currency,
    quote_currency,
    counter_currency,
    last_price,
    close_price,
    waprice,
    accrued_int
)
;
 
