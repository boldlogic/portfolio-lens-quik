-- +goose Up

ALTER TABLE bcs.agreement_accounts
    ALTER COLUMN broker_account NVARCHAR(32) NULL;

ALTER TABLE bcs.agreement_accounts
    ADD is_iia BIT NOT NULL CONSTRAINT DF_bcs_agreement_accounts_is_iia DEFAULT (0),
        is_closed BIT NOT NULL CONSTRAINT DF_bcs_agreement_accounts_is_closed DEFAULT (0);

ALTER TABLE bcs.own_trades
    ADD agreement_code NVARCHAR(32) NULL,
        is_iia BIT NOT NULL CONSTRAINT DF_bcs_own_trades_is_iia DEFAULT (0);

ALTER TABLE bcs.own_trades
    ALTER COLUMN client_code NVARCHAR(32) NULL;

ALTER TABLE bcs.cash_operations
    ADD agreement_code NVARCHAR(32) NULL,
        client_code NVARCHAR(32) NULL;

ALTER TABLE bcs.report_imports
    ADD agreement_code NVARCHAR(32) NULL,
        broker_account NVARCHAR(32) NULL,
        merge_version INT NOT NULL CONSTRAINT DF_bcs_report_imports_merge_version DEFAULT (1);

ALTER TABLE bcs.own_trades DROP CONSTRAINT UQ_bcs_own_trades_num;
DROP INDEX UQ_bcs_own_trades_report_trade_number ON bcs.own_trades;
ALTER TABLE bcs.cash_operations DROP CONSTRAINT UQ_bcs_cash_operations_broker;

ALTER TABLE bcs.own_trades
    ADD CONSTRAINT CK_bcs_own_trades_account_scope
        CHECK (agreement_code IS NOT NULL OR client_code IS NOT NULL);

CREATE UNIQUE NONCLUSTERED INDEX UQ_bcs_own_trades_agreement_trade
    ON bcs.own_trades (agreement_code, trade_num)
    WHERE agreement_code IS NOT NULL;

CREATE UNIQUE NONCLUSTERED INDEX UQ_bcs_own_trades_client_trade
    ON bcs.own_trades (client_code, trade_num)
    WHERE agreement_code IS NULL;

CREATE UNIQUE NONCLUSTERED INDEX UQ_bcs_own_trades_agreement_report_number
    ON bcs.own_trades (agreement_code, report_trade_number)
    WHERE agreement_code IS NOT NULL AND report_trade_number IS NOT NULL;

CREATE UNIQUE NONCLUSTERED INDEX UQ_bcs_own_trades_client_report_number
    ON bcs.own_trades (client_code, report_trade_number)
    WHERE agreement_code IS NULL AND client_code IS NOT NULL AND report_trade_number IS NOT NULL;

CREATE UNIQUE NONCLUSTERED INDEX UQ_bcs_cash_operations_agreement_broker
    ON bcs.cash_operations (agreement_code, broker_id)
    WHERE agreement_code IS NOT NULL;

CREATE UNIQUE NONCLUSTERED INDEX UQ_bcs_cash_operations_broker_unscoped
    ON bcs.cash_operations (broker_id)
    WHERE agreement_code IS NULL;

DECLARE @client_id INT;
DECLARE @client_count INT;
SELECT @client_id = MIN(client_id), @client_count = COUNT(*)
FROM ref.clients
WHERE code = N'1475023';
IF @client_count <> 1 SET @client_id = NULL;

IF NOT EXISTS (SELECT 1 FROM bcs.agreement_accounts WHERE agreement_code = N'1857324/22-аб')
    INSERT INTO bcs.agreement_accounts (agreement_code, broker_account, client_id, is_iia)
    VALUES (N'1857324/22-аб', N'1475023', @client_id, 0);
ELSE
    UPDATE bcs.agreement_accounts
    SET broker_account = N'1475023', client_id = COALESCE(client_id, @client_id), is_iia = 0, is_closed = 0
    WHERE agreement_code = N'1857324/22-аб';

IF NOT EXISTS (SELECT 1 FROM bcs.agreement_accounts WHERE agreement_code = N'1951062/22-аб-иис')
    INSERT INTO bcs.agreement_accounts (agreement_code, broker_account, client_id, is_iia, is_closed)
    VALUES (N'1951062/22-аб-иис', NULL, NULL, 1, 1);
ELSE
    UPDATE bcs.agreement_accounts
    SET broker_account = NULL, client_id = NULL, is_iia = 1, is_closed = 1
    WHERE agreement_code = N'1951062/22-аб-иис';

-- +goose Down

IF EXISTS (SELECT trade_num FROM bcs.own_trades GROUP BY trade_num HAVING COUNT(*) > 1)
    THROW 51000, 'Cannot revert agreement-scoped trades while trade numbers repeat across agreements.', 1;

IF EXISTS (SELECT broker_id FROM bcs.cash_operations GROUP BY broker_id HAVING COUNT(*) > 1)
    THROW 51000, 'Cannot revert agreement-scoped operations while broker IDs repeat across agreements.', 1;

IF EXISTS (SELECT 1 FROM bcs.own_trades WHERE client_code IS NULL OR is_iia = 1)
    THROW 51000, 'Cannot revert while trades rely on agreement-only or IIS account scope.', 1;

DELETE FROM bcs.agreement_accounts
WHERE agreement_code IN (N'1857324/22-аб', N'1951062/22-аб-иис');

DROP INDEX UQ_bcs_cash_operations_broker_unscoped ON bcs.cash_operations;
DROP INDEX UQ_bcs_cash_operations_agreement_broker ON bcs.cash_operations;
DROP INDEX UQ_bcs_own_trades_client_report_number ON bcs.own_trades;
DROP INDEX UQ_bcs_own_trades_agreement_report_number ON bcs.own_trades;
DROP INDEX UQ_bcs_own_trades_client_trade ON bcs.own_trades;
DROP INDEX UQ_bcs_own_trades_agreement_trade ON bcs.own_trades;

ALTER TABLE bcs.cash_operations
    ADD CONSTRAINT UQ_bcs_cash_operations_broker UNIQUE (broker_id);

ALTER TABLE bcs.own_trades
    ADD CONSTRAINT UQ_bcs_own_trades_num UNIQUE (trade_num);

CREATE UNIQUE NONCLUSTERED INDEX UQ_bcs_own_trades_report_trade_number
    ON bcs.own_trades (client_code, report_trade_number)
    WHERE report_trade_number IS NOT NULL;

ALTER TABLE bcs.own_trades DROP CONSTRAINT CK_bcs_own_trades_account_scope;
ALTER TABLE bcs.own_trades DROP CONSTRAINT DF_bcs_own_trades_is_iia;
ALTER TABLE bcs.own_trades DROP COLUMN agreement_code, is_iia;
ALTER TABLE bcs.own_trades ALTER COLUMN client_code NVARCHAR(32) NOT NULL;

ALTER TABLE bcs.cash_operations DROP COLUMN agreement_code, client_code;
ALTER TABLE bcs.agreement_accounts DROP CONSTRAINT DF_bcs_agreement_accounts_is_iia;
ALTER TABLE bcs.agreement_accounts DROP CONSTRAINT DF_bcs_agreement_accounts_is_closed;
ALTER TABLE bcs.agreement_accounts DROP COLUMN is_iia, is_closed;
ALTER TABLE bcs.agreement_accounts ALTER COLUMN broker_account NVARCHAR(32) NOT NULL;

ALTER TABLE bcs.report_imports DROP CONSTRAINT DF_bcs_report_imports_merge_version;
ALTER TABLE bcs.report_imports DROP COLUMN agreement_code, broker_account, merge_version;
