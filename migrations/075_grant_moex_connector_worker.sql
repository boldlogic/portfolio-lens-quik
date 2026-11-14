-- +goose Up

GRANT SELECT, INSERT, UPDATE, DELETE ON SCHEMA::moex TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON market.moex_quotes TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON market.moex_historical_rates TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON market.moex_indices TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON market.moex_index_components TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON moex.common_board_map TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.external_codes TO quik_instruments_worker;
GRANT SELECT ON ref.external_systems TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.instrument_boards TO quik_instruments_worker;
GRANT SELECT, INSERT, UPDATE ON ref.instruments TO quik_instruments_worker;
GRANT SELECT ON ref.boards TO quik_instruments_worker;
GRANT SELECT ON ref.currencies TO quik_instruments_worker;
GRANT SELECT ON ref.instrument_types TO quik_instruments_worker;

-- +goose Down

REVOKE SELECT, INSERT, UPDATE, DELETE ON SCHEMA::moex FROM quik_instruments_worker;
REVOKE SELECT, INSERT, UPDATE ON market.moex_quotes FROM quik_instruments_worker;
REVOKE SELECT, INSERT, UPDATE ON market.moex_historical_rates FROM quik_instruments_worker;
REVOKE SELECT, INSERT, UPDATE ON market.moex_indices FROM quik_instruments_worker;
REVOKE SELECT, INSERT, UPDATE ON market.moex_index_components FROM quik_instruments_worker;
REVOKE SELECT, INSERT, UPDATE ON moex.common_board_map FROM quik_instruments_worker;
REVOKE SELECT, INSERT, UPDATE ON ref.external_codes FROM quik_instruments_worker;
REVOKE SELECT ON ref.external_systems FROM quik_instruments_worker;
REVOKE SELECT, INSERT, UPDATE ON ref.instrument_boards FROM quik_instruments_worker;
REVOKE SELECT, INSERT, UPDATE ON ref.instruments FROM quik_instruments_worker;
REVOKE SELECT ON ref.boards FROM quik_instruments_worker;
REVOKE SELECT ON ref.currencies FROM quik_instruments_worker;
REVOKE SELECT ON ref.instrument_types FROM quik_instruments_worker;
