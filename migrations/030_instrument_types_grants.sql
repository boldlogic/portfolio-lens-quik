-- +goose Up

GRANT SELECT ON ref.instrument_types TO quik_currency_worker;
GRANT SELECT ON ref.instrument_type_boards TO quik_currency_worker;
GRANT SELECT ON ref.boards TO quik_currency_worker;