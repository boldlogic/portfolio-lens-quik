-- +goose Up
GRANT UPDATE ON quik.money_limits TO quik_portfolio_writer;
GRANT UPDATE ON quik.security_limits TO quik_portfolio_writer;
GRANT UPDATE ON quik.security_limits_otc TO quik_portfolio_writer;