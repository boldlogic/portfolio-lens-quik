-- +goose Up

-- +goose StatementBegin
IF EXISTS (
    SELECT 1
    FROM ref.external_codes
    WHERE ext_code_type_id IN (4, 5, 6)
)
BEGIN
    THROW 51000, N'Типы внешних кодов 4, 5 или 6 уже заняты; проверьте справочник перед миграцией T-Invest', 1;
END;
-- +goose StatementEnd

CREATE TABLE ref.external_code_types (
    ext_code_type_id TINYINT NOT NULL,
    code             VARCHAR(32) NOT NULL,
    name             NVARCHAR(100) NOT NULL,
    CONSTRAINT PK_ref_external_code_types PRIMARY KEY CLUSTERED (ext_code_type_id),
    CONSTRAINT UQ_ref_external_code_types_code UNIQUE (code)
);

INSERT INTO ref.external_code_types (ext_code_type_id, code, name)
VALUES
    (1, 'currency', N'Код валюты источника'),
    (3, 'instrument_code', N'Код инструмента источника'),
    (4, 'figi', N'FIGI'),
    (5, 'isin', N'ISIN'),
    (6, 'instrument_uid', N'UID инструмента T-Invest');

INSERT INTO ref.external_systems (ext_system)
VALUES ('TINVEST');

CREATE TABLE portfolio.source_accounts (
    source_account_id   BIGINT         NOT NULL IDENTITY(1, 1),
    ext_system_id       TINYINT        NOT NULL,
    external_account_id NVARCHAR(128)  NOT NULL,
    name                NVARCHAR(200)  NOT NULL DEFAULT N'',
    account_type        NVARCHAR(64)   NOT NULL DEFAULT N'',
    CONSTRAINT PK_portfolio_source_accounts PRIMARY KEY CLUSTERED (source_account_id),
    CONSTRAINT FK_portfolio_source_accounts_system FOREIGN KEY (ext_system_id)
        REFERENCES ref.external_systems (ext_system_id),
    CONSTRAINT UQ_portfolio_source_accounts_external UNIQUE (ext_system_id, external_account_id)
);

ALTER TABLE portfolio.positions
    ADD source_account_id BIGINT NULL;

ALTER TABLE portfolio.positions
    ADD CONSTRAINT FK_portfolio_positions_source_account FOREIGN KEY (source_account_id)
        REFERENCES portfolio.source_accounts (source_account_id);

ALTER TABLE portfolio.positions
    DROP CONSTRAINT CK_portfolio_positions_holder;

ALTER TABLE portfolio.positions
    ALTER COLUMN balance DECIMAL(28, 9) NULL;

ALTER TABLE portfolio.positions
    ADD CONSTRAINT CK_portfolio_positions_holder CHECK (
        (client_id IS NOT NULL AND trade_account_id IS NULL AND source_account_id IS NULL)
        OR (client_id IS NULL AND trade_account_id IS NOT NULL AND source_account_id IS NULL)
        OR (client_id IS NULL AND trade_account_id IS NULL AND source_account_id IS NOT NULL)
    );

CREATE UNIQUE NONCLUSTERED INDEX UQ_portfolio_positions_source_account
    ON portfolio.positions (load_date, source_account_id, instrument_id)
    WHERE source_account_id IS NOT NULL;

GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON portfolio.source_accounts TO quik_portfolio_writer;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON portfolio.positions TO quik_portfolio_writer;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON ref.external_code_types TO quik_portfolio_writer;
GRANT INSERT ON ref.instruments TO quik_portfolio_writer;
GRANT UPDATE ON ref.external_codes TO quik_portfolio_writer;
GRANT SELECT ON ref.boards TO quik_portfolio_writer;
GRANT SELECT ON ref.instrument_type_boards TO quik_portfolio_writer;

-- +goose Down

-- +goose StatementBegin
IF EXISTS (
    SELECT 1
    FROM portfolio.positions
    WHERE source_account_id IS NOT NULL
)
    THROW 51001, N'Сначала удалите или перенесите позиции T-Invest', 1;
-- +goose StatementEnd

-- +goose StatementBegin
IF EXISTS (SELECT 1 FROM portfolio.source_accounts)
    THROW 51002, N'Сначала удалите или перенесите счета T-Invest', 1;
-- +goose StatementEnd

-- +goose StatementBegin
IF EXISTS (
    SELECT 1
    FROM ref.external_codes ec
    JOIN ref.external_systems es ON es.ext_system_id = ec.ext_system_id
    WHERE es.ext_system = 'TINVEST'
)
    THROW 51003, N'Сначала удалите внешние коды T-Invest', 1;
-- +goose StatementEnd

REVOKE SELECT, INSERT, UPDATE, DELETE, REFERENCES ON ref.external_code_types FROM quik_portfolio_writer;
REVOKE SELECT, INSERT, UPDATE, DELETE, REFERENCES ON portfolio.positions FROM quik_portfolio_writer;
REVOKE SELECT, INSERT, UPDATE, DELETE, REFERENCES ON portfolio.source_accounts FROM quik_portfolio_writer;
REVOKE INSERT ON ref.instruments FROM quik_portfolio_writer;
REVOKE UPDATE ON ref.external_codes FROM quik_portfolio_writer;
REVOKE SELECT ON ref.instrument_type_boards FROM quik_portfolio_writer;
REVOKE SELECT ON ref.boards FROM quik_portfolio_writer;

DROP INDEX UQ_portfolio_positions_source_account ON portfolio.positions;
ALTER TABLE portfolio.positions DROP CONSTRAINT CK_portfolio_positions_holder;
ALTER TABLE portfolio.positions DROP CONSTRAINT FK_portfolio_positions_source_account;
ALTER TABLE portfolio.positions DROP COLUMN source_account_id;
ALTER TABLE portfolio.positions ALTER COLUMN balance DECIMAL(19, 4) NULL;

ALTER TABLE portfolio.positions
    ADD CONSTRAINT CK_portfolio_positions_holder CHECK (
        (client_id IS NOT NULL AND trade_account_id IS NULL)
        OR (client_id IS NULL AND trade_account_id IS NOT NULL)
    );

DROP TABLE portfolio.source_accounts;
DROP TABLE ref.external_code_types;

DELETE FROM ref.external_systems WHERE ext_system = 'TINVEST';
