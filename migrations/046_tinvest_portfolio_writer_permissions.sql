-- +goose Up
-- +goose StatementBegin
IF DATABASE_PRINCIPAL_ID(N't_invest_portfolio_writer') IS NULL
    THROW 51000, N'Сначала создайте пользователя и роль T-Invest bootstrap-скриптом', 1;

GRANT SELECT, INSERT, UPDATE, DELETE ON OBJECT::tinvest.current_state TO t_invest_portfolio_writer;
GRANT SELECT, INSERT, UPDATE ON OBJECT::tinvest.history_record TO t_invest_portfolio_writer;
GRANT SELECT, INSERT, UPDATE ON OBJECT::tinvest.instrument TO t_invest_portfolio_writer;

GRANT SELECT, INSERT, UPDATE, DELETE ON OBJECT::portfolio.source_accounts TO t_invest_portfolio_writer;
GRANT SELECT, INSERT, UPDATE, DELETE ON OBJECT::portfolio.positions TO t_invest_portfolio_writer;

GRANT SELECT ON OBJECT::ref.external_systems TO t_invest_portfolio_writer;
GRANT SELECT ON OBJECT::ref.boards TO t_invest_portfolio_writer;
GRANT SELECT ON OBJECT::ref.instrument_type_boards TO t_invest_portfolio_writer;
GRANT SELECT, INSERT, UPDATE ON OBJECT::ref.instruments TO t_invest_portfolio_writer;
GRANT SELECT, INSERT, UPDATE ON OBJECT::ref.external_codes TO t_invest_portfolio_writer;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
REVOKE SELECT, INSERT, UPDATE, DELETE ON OBJECT::tinvest.current_state FROM t_invest_portfolio_writer;
REVOKE SELECT, INSERT, UPDATE ON OBJECT::tinvest.history_record FROM t_invest_portfolio_writer;
REVOKE SELECT, INSERT, UPDATE ON OBJECT::tinvest.instrument FROM t_invest_portfolio_writer;

REVOKE SELECT, INSERT, UPDATE, DELETE ON OBJECT::portfolio.source_accounts FROM t_invest_portfolio_writer;
REVOKE SELECT, INSERT, UPDATE, DELETE ON OBJECT::portfolio.positions FROM t_invest_portfolio_writer;

REVOKE SELECT ON OBJECT::ref.external_systems FROM t_invest_portfolio_writer;
REVOKE SELECT ON OBJECT::ref.boards FROM t_invest_portfolio_writer;
REVOKE SELECT ON OBJECT::ref.instrument_type_boards FROM t_invest_portfolio_writer;
REVOKE SELECT, INSERT, UPDATE ON OBJECT::ref.instruments FROM t_invest_portfolio_writer;
REVOKE SELECT, INSERT, UPDATE ON OBJECT::ref.external_codes FROM t_invest_portfolio_writer;
-- +goose StatementEnd
