-- +goose Up

INSERT INTO bcs.agreement_accounts (
    agreement_code, broker_account, forts_account, depo_account, client_id,
    is_iia, is_closed, login_id
)
SELECT N'2713365/23-' + NCHAR(1080) + NCHAR(1080) + NCHAR(1089),
       N'2064274', N'82P4k1e', N'CB2628056', NULL, 1, 0, NULL
WHERE NOT EXISTS (
    SELECT 1 FROM bcs.agreement_accounts
    WHERE agreement_code = N'2713365/23-' + NCHAR(1080) + NCHAR(1080) + NCHAR(1089)
);

INSERT INTO bcs.agreement_accounts (
    agreement_code, broker_account, forts_account, depo_account, client_id,
    is_iia, is_closed, login_id
)
SELECT N'2713363/23', N'2064275', N'82P4k1f', N'CB2628057', NULL,
       0, 0, NULL
WHERE NOT EXISTS (
    SELECT 1 FROM bcs.agreement_accounts
    WHERE agreement_code = N'2713363/23'
);

-- +goose Down

IF EXISTS (
    SELECT 1 FROM bcs.report_imports
    WHERE agreement_code IN (N'2713365/23-' + NCHAR(1080) + NCHAR(1080) + NCHAR(1089), N'2713363/23')
) OR EXISTS (
    SELECT 1 FROM bcs.own_trades
    WHERE agreement_code IN (N'2713365/23-' + NCHAR(1080) + NCHAR(1080) + NCHAR(1089), N'2713363/23')
) OR EXISTS (
    SELECT 1 FROM bcs.cash_operations
    WHERE agreement_code IN (N'2713365/23-' + NCHAR(1080) + NCHAR(1080) + NCHAR(1089), N'2713363/23')
)
    THROW 51000, 'Cannot remove BCS agreement mappings while report history refers to them.', 1;

DELETE FROM bcs.agreement_accounts
WHERE agreement_code IN (N'2713365/23-' + NCHAR(1080) + NCHAR(1080) + NCHAR(1089), N'2713363/23');
