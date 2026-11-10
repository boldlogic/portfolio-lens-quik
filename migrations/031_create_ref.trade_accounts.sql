-- +goose Up

CREATE TABLE ref.trade_accounts (
    trade_account_id INT         NOT NULL IDENTITY(1, 1),
    firm_id          TINYINT     NOT NULL,
    client_code      VARCHAR(12) NOT NULL,
    code             VARCHAR(12) NOT NULL,
    CONSTRAINT PK_ref_trade_accounts PRIMARY KEY CLUSTERED (trade_account_id),
    CONSTRAINT FK_ref_trade_accounts_firm FOREIGN KEY (firm_id)
        REFERENCES ref.firms (firm_id)
);

CREATE UNIQUE NONCLUSTERED INDEX UQ_ref_trade_accounts_firm_client_code
    ON ref.trade_accounts (firm_id, client_code, code);

EXEC sys.sp_addextendedproperty
    @name = N'MS_Description',
    @value = N'Торговые счета QUIK. Строка собирается из лимитов по бумагам: фирма, код клиента и торговый счет.',
    @level0type = N'SCHEMA', @level0name = N'ref',
    @level1type = N'TABLE', @level1name = N'trade_accounts';

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Внутренний идентификатор торгового счета.', @level0type = N'SCHEMA', @level0name = N'ref', @level1type = N'TABLE', @level1name = N'trade_accounts', @level2type = N'COLUMN', @level2name = N'trade_account_id';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Фирма. Ссылка на ref.firms.', @level0type = N'SCHEMA', @level0name = N'ref', @level1type = N'TABLE', @level1name = N'trade_accounts', @level2type = N'COLUMN', @level2name = N'firm_id';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код клиента (QUIK: ClientCode).', @level0type = N'SCHEMA', @level0name = N'ref', @level1type = N'TABLE', @level1name = N'trade_accounts', @level2type = N'COLUMN', @level2name = N'client_code';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Торговый счет (QUIK: TrdAcc, в лимитах trade_account).', @level0type = N'SCHEMA', @level0name = N'ref', @level1type = N'TABLE', @level1name = N'trade_accounts', @level2type = N'COLUMN', @level2name = N'code';

GRANT SELECT, INSERT, UPDATE ON ref.trade_accounts TO quik_portfolio_writer;
