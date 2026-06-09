-- +goose Up

GRANT EXECUTE ON TYPE::app.client_code_list TO quik_portfolio_reader;

-- +goose Down

REVOKE EXECUTE ON TYPE::app.client_code_list FROM quik_portfolio_reader;
