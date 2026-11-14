-- +goose Up

ALTER TABLE bcs.agreement_accounts ADD trade_account_id INT NULL;

ALTER TABLE bcs.agreement_accounts
    ADD CONSTRAINT FK_bcs_agreement_accounts_trade FOREIGN KEY (trade_account_id)
        REFERENCES ref.trade_accounts (trade_account_id);

-- +goose StatementBegin
UPDATE a
SET trade_account_id = x.trade_account_id
FROM bcs.agreement_accounts AS a
JOIN (
    SELECT client_id, MIN(trade_account_id) AS trade_account_id
    FROM ref.trade_accounts
    WHERE venue = 'exchange'
    GROUP BY client_id
    HAVING COUNT(*) = 1
) AS x ON x.client_id = a.client_id
WHERE a.client_id IS NOT NULL;
-- +goose StatementEnd

ALTER TABLE portfolio.positions ADD source_updated_at DATETIMEOFFSET(7) NULL;

ALTER TABLE portfolio.positions ADD ext_system_id TINYINT NULL;

ALTER TABLE portfolio.positions
    ADD CONSTRAINT FK_portfolio_positions_ext_system FOREIGN KEY (ext_system_id)
        REFERENCES ref.external_systems (ext_system_id);

-- +goose Down

ALTER TABLE portfolio.positions DROP CONSTRAINT FK_portfolio_positions_ext_system;
ALTER TABLE portfolio.positions DROP COLUMN ext_system_id;
ALTER TABLE portfolio.positions DROP COLUMN source_updated_at;

ALTER TABLE bcs.agreement_accounts DROP CONSTRAINT FK_bcs_agreement_accounts_trade;
ALTER TABLE bcs.agreement_accounts DROP COLUMN trade_account_id;
