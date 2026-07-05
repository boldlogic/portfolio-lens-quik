-- +goose Up

GRANT SELECT, UPDATE ON quik.current_quotes TO quik_instruments_worker;
GRANT SELECT ON ref.currencies TO quik_instruments_worker;
GRANT SELECT ON ref.external_codes TO quik_instruments_worker;
GRANT SELECT ON ref.external_systems TO quik_instruments_worker;
GRANT SELECT ON ref.instrument_types TO quik_instruments_worker;
GRANT SELECT ON ref.instrument_type_boards TO quik_instruments_worker;
GRANT SELECT ON ref.boards TO quik_instruments_worker;
GRANT SELECT ON ref.trade_points TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.instruments TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.instrument_boards TO quik_instruments_worker;
