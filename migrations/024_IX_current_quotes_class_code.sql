-- +goose Up

DROP INDEX IX_current_quotes_class_code
    ON quik.current_quotes;

CREATE NONCLUSTERED INDEX IX_current_quotes_class_instrument ON quik.current_quotes (class_code, instrument_id, rw) INCLUDE (
    instrument_class,
    sec_code,
    currency,
    isin,
    registration_number,
    full_name,
    short_name,
    base_currency,
	quote_currency,
    counter_currency,
    face_value,
    maturity_date,
    coupon_duration
)
;
 
