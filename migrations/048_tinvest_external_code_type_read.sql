-- +goose Up
GRANT SELECT ON OBJECT::ref.external_code_types TO t_invest_portfolio_writer;

-- +goose Down
REVOKE SELECT ON OBJECT::ref.external_code_types FROM t_invest_portfolio_writer;
