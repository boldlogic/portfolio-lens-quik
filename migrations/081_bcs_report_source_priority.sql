-- +goose Up

ALTER TABLE bcs.own_trades
    ADD report_authoritative BIT NOT NULL CONSTRAINT DF_bcs_own_trades_report_authoritative DEFAULT (0),
        report_time_known    BIT NOT NULL CONSTRAINT DF_bcs_own_trades_report_time_known DEFAULT (1),
        report_trade_number  NVARCHAR(64) NULL;

CREATE UNIQUE NONCLUSTERED INDEX UQ_bcs_own_trades_report_trade_number
    ON bcs.own_trades (client_code, report_trade_number)
    WHERE report_trade_number IS NOT NULL;

ALTER TABLE bcs.cash_operations
    ADD report_authoritative BIT NOT NULL CONSTRAINT DF_bcs_cash_operations_report_authoritative DEFAULT (0);

ALTER TABLE bcs.cash_operations
    ALTER COLUMN status NVARCHAR(32) NULL;

-- +goose Down

UPDATE bcs.cash_operations SET status = N'Unknown' WHERE status IS NULL;

ALTER TABLE bcs.cash_operations
    ALTER COLUMN status NVARCHAR(32) NOT NULL;

ALTER TABLE bcs.cash_operations DROP CONSTRAINT DF_bcs_cash_operations_report_authoritative;
ALTER TABLE bcs.cash_operations DROP COLUMN report_authoritative;

ALTER TABLE bcs.own_trades DROP CONSTRAINT DF_bcs_own_trades_report_time_known;
ALTER TABLE bcs.own_trades DROP CONSTRAINT DF_bcs_own_trades_report_authoritative;
DROP INDEX UQ_bcs_own_trades_report_trade_number ON bcs.own_trades;
ALTER TABLE bcs.own_trades DROP COLUMN report_trade_number;
ALTER TABLE bcs.own_trades DROP COLUMN report_time_known;
ALTER TABLE bcs.own_trades DROP COLUMN report_authoritative;
