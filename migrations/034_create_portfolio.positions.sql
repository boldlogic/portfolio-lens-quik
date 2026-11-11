-- +goose Up

CREATE SCHEMA portfolio;

CREATE TABLE portfolio.positions (
    position_id      BIGINT         NOT NULL IDENTITY(1, 1),
    load_date        DATE           NOT NULL,
    client_id        INT            NULL,
    trade_account_id INT            NULL,
    instrument_id    BIGINT         NOT NULL,
    balance          DECIMAL(19, 4) NULL,
    CONSTRAINT PK_portfolio_positions PRIMARY KEY CLUSTERED (position_id),
    CONSTRAINT FK_portfolio_positions_client FOREIGN KEY (client_id)
        REFERENCES ref.clients (client_id),
    CONSTRAINT FK_portfolio_positions_trade_account FOREIGN KEY (trade_account_id)
        REFERENCES ref.trade_accounts (trade_account_id),
    CONSTRAINT FK_portfolio_positions_instrument FOREIGN KEY (instrument_id)
        REFERENCES ref.instruments (instrument_id),
    CONSTRAINT CK_portfolio_positions_holder CHECK (
        (client_id IS NOT NULL AND trade_account_id IS NULL)
        OR (client_id IS NULL AND trade_account_id IS NOT NULL)
    )
);

CREATE UNIQUE NONCLUSTERED INDEX UQ_portfolio_positions_money
    ON portfolio.positions (load_date, client_id, instrument_id)
    WHERE client_id IS NOT NULL;

CREATE UNIQUE NONCLUSTERED INDEX UQ_portfolio_positions_security
    ON portfolio.positions (load_date, trade_account_id, instrument_id)
    WHERE trade_account_id IS NOT NULL;

EXEC sys.sp_addextendedproperty
    @name = N'MS_Description',
    @value = N'Снимок остатка. Деньги ссылаются на клиента, бумаги на торговый счет.',
    @level0type = N'SCHEMA', @level0name = N'portfolio',
    @level1type = N'TABLE', @level1name = N'positions';

GRANT SELECT, INSERT, UPDATE, DELETE ON portfolio.positions TO quik_portfolio_worker;

GRANT SELECT ON ref.firms TO quik_portfolio_worker;
GRANT SELECT ON ref.instruments TO quik_portfolio_worker;
GRANT SELECT ON ref.instrument_types TO quik_portfolio_worker;
GRANT SELECT ON ref.currencies TO quik_portfolio_worker;
GRANT SELECT ON ref.external_codes TO quik_portfolio_worker;
GRANT SELECT ON ref.external_systems TO quik_portfolio_worker;

GRANT SELECT ON ref.instruments TO quik_portfolio_writer;
GRANT SELECT ON ref.external_systems TO quik_portfolio_writer;
GRANT SELECT, INSERT ON ref.external_codes TO quik_portfolio_writer;

-- +goose Down

REVOKE SELECT, INSERT ON ref.external_codes FROM quik_portfolio_writer;
REVOKE SELECT ON ref.external_systems FROM quik_portfolio_writer;
REVOKE SELECT ON ref.instruments FROM quik_portfolio_writer;

REVOKE SELECT ON ref.external_systems FROM quik_portfolio_worker;
REVOKE SELECT ON ref.external_codes FROM quik_portfolio_worker;
REVOKE SELECT ON ref.currencies FROM quik_portfolio_worker;
REVOKE SELECT ON ref.instrument_types FROM quik_portfolio_worker;
REVOKE SELECT ON ref.instruments FROM quik_portfolio_worker;
REVOKE SELECT ON ref.firms FROM quik_portfolio_worker;

DROP TABLE portfolio.positions;
DROP SCHEMA portfolio;
