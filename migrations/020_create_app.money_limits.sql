-- +goose Up

CREATE TYPE app.money_limits
AS TABLE (
    client_code   varchar(12)    NOT NULL,
    currency_code varchar(4)     NOT NULL,
    position_code varchar(4)     NOT NULL,
    settle_code   varchar(5)     NOT NULL DEFAULT 'Tx',
    firm_code     varchar(12)    NOT NULL,
    balance       decimal(19, 4) NOT NULL default 0,
    PRIMARY KEY CLUSTERED (client_code,
        currency_code,
        position_code,
        settle_code,
        firm_code)
);

CREATE TYPE app.security_limits
AS TABLE (
    client_code   varchar(12)    NOT NULL,
    sec_code                  varchar(12)    NOT NULL,
    trade_account             varchar(12)    NOT NULL DEFAULT 'OTC',
    settle_code               varchar(5)     NOT NULL DEFAULT 'Tx',
    firm_code                 varchar(12)    NOT NULL,
    balance                   decimal(19, 4) NOT NULL default 0,
    acquisition_currency_code varchar(4)     NULL, 
    isin                      varchar(12)    NULL,
    PRIMARY KEY CLUSTERED ( client_code,
        sec_code,
        trade_account,
        settle_code,
        firm_code)
);

GRANT EXECUTE ON TYPE::app.money_limits TO quik_portfolio_writer;
GRANT EXECUTE ON TYPE::app.security_limits TO quik_portfolio_writer;