-- +goose Up

GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON SCHEMA::bcs TO quik_portfolio_writer;

GRANT SELECT ON quik.current_quotes TO quik_portfolio_writer;

GRANT SELECT ON ref.instrument_types TO quik_portfolio_writer;
GRANT SELECT ON ref.instrument_type_boards TO quik_portfolio_writer;
GRANT SELECT ON ref.trade_points TO quik_portfolio_writer;
GRANT SELECT ON ref.currencies TO quik_portfolio_writer;
GRANT SELECT, INSERT, REFERENCES ON ref.boards TO quik_portfolio_writer;
GRANT INSERT, REFERENCES ON ref.instruments TO quik_portfolio_writer;
GRANT SELECT, INSERT, REFERENCES ON ref.instrument_boards TO quik_portfolio_writer;
GRANT REFERENCES ON ref.trade_points TO quik_portfolio_writer;
GRANT REFERENCES ON ref.instrument_types TO quik_portfolio_writer;
GRANT REFERENCES ON ref.currencies TO quik_portfolio_writer;

-- +goose Down

REVOKE REFERENCES ON ref.currencies FROM quik_portfolio_writer;
REVOKE REFERENCES ON ref.instrument_types FROM quik_portfolio_writer;
REVOKE REFERENCES ON ref.trade_points FROM quik_portfolio_writer;
REVOKE SELECT, INSERT, REFERENCES ON ref.instrument_boards FROM quik_portfolio_writer;
REVOKE INSERT, REFERENCES ON ref.instruments FROM quik_portfolio_writer;
REVOKE SELECT, INSERT, REFERENCES ON ref.boards FROM quik_portfolio_writer;
REVOKE SELECT ON ref.currencies FROM quik_portfolio_writer;
REVOKE SELECT ON ref.trade_points FROM quik_portfolio_writer;
REVOKE SELECT ON ref.instrument_type_boards FROM quik_portfolio_writer;
REVOKE SELECT ON ref.instrument_types FROM quik_portfolio_writer;

REVOKE SELECT ON quik.current_quotes FROM quik_portfolio_writer;

REVOKE SELECT, INSERT, UPDATE, DELETE, REFERENCES ON SCHEMA::bcs FROM quik_portfolio_writer;
