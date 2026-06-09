-- +goose Up

GRANT SELECT ON quik.security_limits_otc TO quik_portfolio_writer;

-- +goose Down

REVOKE SELECT ON quik.security_limits_otc FROM quik_portfolio_writer;
