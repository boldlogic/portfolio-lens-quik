-- +goose Up

GRANT SELECT ON quik.money_limits TO quik_portfolio_reader;
GRANT SELECT ON quik.security_limits TO quik_portfolio_reader;
GRANT SELECT ON quik.security_limits_otc TO quik_portfolio_reader;
GRANT SELECT ON quik.current_quotes TO quik_portfolio_reader;

GRANT SELECT ON SCHEMA::ref TO quik_portfolio_reader;

GRANT SELECT ON market.fx_cbr_rates TO quik_portfolio_reader;
GRANT SELECT ON market.fnFxRateCross TO quik_portfolio_reader;
GRANT SELECT ON market.fnSecurityQuoteByAcquisitionCurrency TO quik_portfolio_reader;

GRANT REFERENCES ON TYPE::app.client_code_list TO quik_portfolio_reader;

GRANT INSERT ON quik.money_limits TO quik_portfolio_writer;
GRANT INSERT ON quik.security_limits TO quik_portfolio_writer;
GRANT INSERT ON quik.security_limits_otc TO quik_portfolio_writer;
GRANT SELECT ON quik.money_limits TO quik_portfolio_writer;
GRANT SELECT ON quik.security_limits TO quik_portfolio_writer;
GRANT SELECT ON ref.firms TO quik_portfolio_writer;
GRANT INSERT, UPDATE ON ref.firms TO quik_portfolio_writer;

GRANT SELECT, INSERT, DELETE ON quik.money_limits TO quik_portfolio_worker;
GRANT SELECT, INSERT, DELETE ON quik.security_limits TO quik_portfolio_worker;
GRANT SELECT, INSERT, DELETE ON quik.security_limits_otc TO quik_portfolio_worker;

GRANT SELECT ON quik.current_quotes TO quik_currency_worker;
GRANT SELECT, INSERT, UPDATE ON ref.currencies TO quik_currency_worker;
GRANT SELECT, INSERT, UPDATE ON ref.external_codes TO quik_currency_worker;
GRANT SELECT, INSERT, UPDATE ON ref.external_systems TO quik_currency_worker;
GRANT SELECT, INSERT, UPDATE ON market.fx_cbr_rates TO quik_currency_worker;
