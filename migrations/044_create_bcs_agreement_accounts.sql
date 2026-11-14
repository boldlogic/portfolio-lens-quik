-- +goose Up

CREATE TABLE bcs.agreement_accounts (
    agreement_account_id BIGINT        NOT NULL IDENTITY(1, 1),
    agreement_code       NVARCHAR(32)  NOT NULL,
    broker_account       NVARCHAR(32)  NOT NULL,
    forts_account        NVARCHAR(32)  NULL,
    depo_account         NVARCHAR(32)  NULL,
    client_id            INT           NULL,
    CONSTRAINT PK_bcs_agreement_accounts PRIMARY KEY CLUSTERED (agreement_account_id),
    CONSTRAINT UQ_bcs_agreement_accounts_code UNIQUE (agreement_code),
    CONSTRAINT FK_bcs_agreement_accounts_client FOREIGN KEY (client_id)
        REFERENCES ref.clients (client_id)
);

-- +goose StatementBegin
DECLARE @n INT;
DECLARE @client_id INT;

SELECT @n = COUNT(*)
FROM ref.clients
WHERE code = '1475023';

IF @n = 1
    SELECT @client_id = client_id
    FROM ref.clients
    WHERE code = '1475023';

INSERT INTO bcs.agreement_accounts (
    agreement_code, broker_account, forts_account, depo_account, client_id
)
VALUES (
    N'1688487/22', N'1475023', N'82LL5sd', N'СВ1675519', @client_id
);
-- +goose StatementEnd

-- +goose Down

DROP TABLE IF EXISTS bcs.agreement_accounts;
