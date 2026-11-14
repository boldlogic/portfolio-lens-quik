-- +goose Up

ALTER TABLE bcs.report_imports
    ADD report_period_end DATE NULL;

ALTER TABLE bcs.own_trades
    ADD report_leg_key NVARCHAR(64) NULL,
        repo_part TINYINT NULL,
        report_period_end DATE NULL,
        report_import_id BIGINT NULL;

ALTER TABLE bcs.own_trades
    ADD CONSTRAINT CK_bcs_own_trades_repo_part CHECK (repo_part IS NULL OR repo_part IN (1, 2));

ALTER TABLE bcs.own_trades
    ADD CONSTRAINT FK_bcs_own_trades_report_import
        FOREIGN KEY (report_import_id) REFERENCES bcs.report_imports (import_id);

ALTER TABLE bcs.cash_operations
    ADD report_period_end DATE NULL,
        report_import_id BIGINT NULL;

ALTER TABLE bcs.cash_operations
    ADD CONSTRAINT FK_bcs_cash_operations_report_import
        FOREIGN KEY (report_import_id) REFERENCES bcs.report_imports (import_id);

DROP INDEX UQ_bcs_own_trades_agreement_report_number ON bcs.own_trades;
DROP INDEX UQ_bcs_own_trades_client_report_number ON bcs.own_trades;

CREATE UNIQUE NONCLUSTERED INDEX UQ_bcs_own_trades_agreement_report_leg
    ON bcs.own_trades (agreement_code, report_trade_number, report_leg_key)
    WHERE agreement_code IS NOT NULL AND report_trade_number IS NOT NULL AND report_leg_key IS NOT NULL;

CREATE UNIQUE NONCLUSTERED INDEX UQ_bcs_own_trades_client_report_leg
    ON bcs.own_trades (client_code, report_trade_number, report_leg_key)
    WHERE agreement_code IS NULL AND client_code IS NOT NULL AND report_trade_number IS NOT NULL AND report_leg_key IS NOT NULL;

-- +goose Down

IF EXISTS (
    SELECT agreement_code, report_trade_number
    FROM bcs.own_trades
    WHERE report_trade_number IS NOT NULL
    GROUP BY agreement_code, report_trade_number
    HAVING COUNT(*) > 1
)
    THROW 51000, 'Cannot revert report trade legs while report numbers repeat within an agreement.', 1;

DROP INDEX UQ_bcs_own_trades_client_report_leg ON bcs.own_trades;
DROP INDEX UQ_bcs_own_trades_agreement_report_leg ON bcs.own_trades;

CREATE UNIQUE NONCLUSTERED INDEX UQ_bcs_own_trades_agreement_report_number
    ON bcs.own_trades (agreement_code, report_trade_number)
    WHERE agreement_code IS NOT NULL AND report_trade_number IS NOT NULL;

CREATE UNIQUE NONCLUSTERED INDEX UQ_bcs_own_trades_client_report_number
    ON bcs.own_trades (client_code, report_trade_number)
    WHERE agreement_code IS NULL AND client_code IS NOT NULL AND report_trade_number IS NOT NULL;

ALTER TABLE bcs.cash_operations DROP CONSTRAINT FK_bcs_cash_operations_report_import;
ALTER TABLE bcs.cash_operations DROP COLUMN report_period_end, report_import_id;

ALTER TABLE bcs.own_trades DROP CONSTRAINT FK_bcs_own_trades_report_import;
ALTER TABLE bcs.own_trades DROP CONSTRAINT CK_bcs_own_trades_repo_part;
ALTER TABLE bcs.own_trades DROP COLUMN report_leg_key, repo_part, report_period_end, report_import_id;

ALTER TABLE bcs.report_imports DROP COLUMN report_period_end;
