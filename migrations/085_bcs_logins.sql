-- +goose Up

CREATE TABLE bcs.logins (
    login_id   BIGINT       NOT NULL IDENTITY(1, 1),
    label      NVARCHAR(64) NOT NULL,
    token_kind NVARCHAR(8)  NOT NULL,
    CONSTRAINT PK_bcs_logins PRIMARY KEY CLUSTERED (login_id),
    CONSTRAINT UQ_bcs_logins_label UNIQUE (label),
    CONSTRAINT CK_bcs_logins_kind CHECK (token_kind IN (N'read', N'trade'))
);

ALTER TABLE bcs.agreement_accounts
    ADD login_id BIGINT NULL,
        CONSTRAINT FK_bcs_agreement_accounts_login FOREIGN KEY (login_id)
            REFERENCES bcs.logins (login_id);

ALTER TABLE bcs.depo_limits ADD login_id BIGINT NULL,
    CONSTRAINT FK_bcs_depo_limits_login FOREIGN KEY (login_id) REFERENCES bcs.logins (login_id);
ALTER TABLE bcs.money_limits ADD login_id BIGINT NULL,
    CONSTRAINT FK_bcs_money_limits_login FOREIGN KEY (login_id) REFERENCES bcs.logins (login_id);
ALTER TABLE bcs.future_holdings ADD login_id BIGINT NULL,
    CONSTRAINT FK_bcs_future_holdings_login FOREIGN KEY (login_id) REFERENCES bcs.logins (login_id);
ALTER TABLE bcs.futures_limits ADD login_id BIGINT NULL,
    CONSTRAINT FK_bcs_futures_limits_login FOREIGN KEY (login_id) REFERENCES bcs.logins (login_id);

ALTER TABLE bcs.cash_operations ADD login_id BIGINT NULL,
    CONSTRAINT FK_bcs_cash_operations_login FOREIGN KEY (login_id) REFERENCES bcs.logins (login_id);

DROP INDEX UQ_bcs_cash_operations_broker_unscoped ON bcs.cash_operations;

CREATE UNIQUE NONCLUSTERED INDEX UQ_bcs_cash_operations_broker_unscoped
    ON bcs.cash_operations (broker_id)
    WHERE agreement_code IS NULL AND login_id IS NULL;

CREATE UNIQUE NONCLUSTERED INDEX UQ_bcs_cash_operations_login_broker
    ON bcs.cash_operations (login_id, broker_id)
    WHERE agreement_code IS NULL AND login_id IS NOT NULL;

-- +goose Down

DROP INDEX UQ_bcs_cash_operations_login_broker ON bcs.cash_operations;
DROP INDEX UQ_bcs_cash_operations_broker_unscoped ON bcs.cash_operations;

CREATE UNIQUE NONCLUSTERED INDEX UQ_bcs_cash_operations_broker_unscoped
    ON bcs.cash_operations (broker_id)
    WHERE agreement_code IS NULL;

ALTER TABLE bcs.cash_operations DROP CONSTRAINT FK_bcs_cash_operations_login;
ALTER TABLE bcs.cash_operations DROP COLUMN login_id;

ALTER TABLE bcs.futures_limits DROP CONSTRAINT FK_bcs_futures_limits_login;
ALTER TABLE bcs.futures_limits DROP COLUMN login_id;
ALTER TABLE bcs.future_holdings DROP CONSTRAINT FK_bcs_future_holdings_login;
ALTER TABLE bcs.future_holdings DROP COLUMN login_id;
ALTER TABLE bcs.money_limits DROP CONSTRAINT FK_bcs_money_limits_login;
ALTER TABLE bcs.money_limits DROP COLUMN login_id;
ALTER TABLE bcs.depo_limits DROP CONSTRAINT FK_bcs_depo_limits_login;
ALTER TABLE bcs.depo_limits DROP COLUMN login_id;

ALTER TABLE bcs.agreement_accounts DROP CONSTRAINT FK_bcs_agreement_accounts_login;
ALTER TABLE bcs.agreement_accounts DROP COLUMN login_id;

DROP TABLE bcs.logins;
