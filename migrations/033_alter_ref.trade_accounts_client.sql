-- +goose Up

INSERT INTO ref.clients (firm_id, code)
SELECT DISTINCT t.firm_id, t.client_code
FROM ref.trade_accounts t
WHERE NOT EXISTS (
    SELECT 1
    FROM ref.clients c
    WHERE c.firm_id = t.firm_id
        AND c.code = t.client_code
);

ALTER TABLE ref.trade_accounts ADD client_id INT NULL;
ALTER TABLE ref.trade_accounts ADD venue VARCHAR(8) NULL;

UPDATE t
SET client_id = c.client_id
FROM ref.trade_accounts t
JOIN ref.clients c
    ON c.firm_id = t.firm_id
    AND c.code = t.client_code;

UPDATE t
SET venue = CASE
    WHEN EXISTS (
        SELECT 1
        FROM quik.security_limits s
        JOIN ref.firms f ON f.firm_id = t.firm_id
        WHERE s.firm_code = f.code
            AND s.client_code = t.client_code
            AND s.trade_account = t.code
    ) THEN 'exchange'
    ELSE 'otc'
END
FROM ref.trade_accounts t;

-- +goose StatementBegin
IF EXISTS (
    SELECT 1
    FROM ref.trade_accounts
    WHERE client_id IS NULL OR venue IS NULL
)
BEGIN
    THROW 51000, N'Не удалось привязать торговый счет к клиенту', 1;
END
-- +goose StatementEnd

ALTER TABLE ref.trade_accounts ALTER COLUMN client_id INT NOT NULL;
ALTER TABLE ref.trade_accounts ALTER COLUMN venue VARCHAR(8) NOT NULL;

ALTER TABLE ref.trade_accounts DROP CONSTRAINT FK_ref_trade_accounts_firm;

DROP INDEX UQ_ref_trade_accounts_firm_client_code ON ref.trade_accounts;

ALTER TABLE ref.trade_accounts DROP COLUMN firm_id;
ALTER TABLE ref.trade_accounts DROP COLUMN client_code;

ALTER TABLE ref.trade_accounts
    ADD CONSTRAINT FK_ref_trade_accounts_client FOREIGN KEY (client_id)
        REFERENCES ref.clients (client_id);

ALTER TABLE ref.trade_accounts
    ADD CONSTRAINT CK_ref_trade_accounts_venue CHECK (venue IN ('exchange', 'otc'));

CREATE UNIQUE NONCLUSTERED INDEX UQ_ref_trade_accounts_client_venue_code
    ON ref.trade_accounts (client_id, venue, code);

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Клиент. Ссылка на ref.clients.', @level0type = N'SCHEMA', @level0name = N'ref', @level1type = N'TABLE', @level1name = N'trade_accounts', @level2type = N'COLUMN', @level2name = N'client_id';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Вид счета: exchange или otc.', @level0type = N'SCHEMA', @level0name = N'ref', @level1type = N'TABLE', @level1name = N'trade_accounts', @level2type = N'COLUMN', @level2name = N'venue';

GRANT SELECT ON ref.trade_accounts TO quik_portfolio_worker;

-- +goose Down

ALTER TABLE ref.trade_accounts ADD firm_id TINYINT NULL;
ALTER TABLE ref.trade_accounts ADD client_code VARCHAR(12) NULL;

UPDATE t
SET firm_id = c.firm_id,
    client_code = c.code
FROM ref.trade_accounts t
JOIN ref.clients c ON c.client_id = t.client_id;

-- +goose StatementBegin
IF EXISTS (
    SELECT 1
    FROM ref.trade_accounts
    WHERE firm_id IS NULL OR client_code IS NULL
)
BEGIN
    THROW 51000, N'Не удалось вернуть фирму и код клиента на торговый счет', 1;
END
-- +goose StatementEnd

ALTER TABLE ref.trade_accounts ALTER COLUMN firm_id TINYINT NOT NULL;
ALTER TABLE ref.trade_accounts ALTER COLUMN client_code VARCHAR(12) NOT NULL;

ALTER TABLE ref.trade_accounts DROP CONSTRAINT FK_ref_trade_accounts_client;
ALTER TABLE ref.trade_accounts DROP CONSTRAINT CK_ref_trade_accounts_venue;

DROP INDEX UQ_ref_trade_accounts_client_venue_code ON ref.trade_accounts;

ALTER TABLE ref.trade_accounts DROP COLUMN client_id;
ALTER TABLE ref.trade_accounts DROP COLUMN venue;

ALTER TABLE ref.trade_accounts
    ADD CONSTRAINT FK_ref_trade_accounts_firm FOREIGN KEY (firm_id)
        REFERENCES ref.firms (firm_id);

CREATE UNIQUE NONCLUSTERED INDEX UQ_ref_trade_accounts_firm_client_code
    ON ref.trade_accounts (firm_id, client_code, code);

REVOKE SELECT ON ref.trade_accounts FROM quik_portfolio_worker;
