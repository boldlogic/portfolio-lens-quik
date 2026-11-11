-- +goose Up

CREATE TABLE ref.clients (
    client_id INT         NOT NULL IDENTITY(1, 1),
    firm_id   TINYINT     NOT NULL,
    code      VARCHAR(12) NOT NULL,
    CONSTRAINT PK_ref_clients PRIMARY KEY CLUSTERED (client_id),
    CONSTRAINT FK_ref_clients_firm FOREIGN KEY (firm_id)
        REFERENCES ref.firms (firm_id)
);

CREATE UNIQUE NONCLUSTERED INDEX UQ_ref_clients_firm_code
    ON ref.clients (firm_id, code);

EXEC sys.sp_addextendedproperty
    @name = N'MS_Description',
    @value = N'Клиент QUIK: фирма и код клиента. Торговый счет здесь не хранится.',
    @level0type = N'SCHEMA', @level0name = N'ref',
    @level1type = N'TABLE', @level1name = N'clients';

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Внутренний идентификатор клиента.', @level0type = N'SCHEMA', @level0name = N'ref', @level1type = N'TABLE', @level1name = N'clients', @level2type = N'COLUMN', @level2name = N'client_id';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Фирма. Ссылка на ref.firms.', @level0type = N'SCHEMA', @level0name = N'ref', @level1type = N'TABLE', @level1name = N'clients', @level2type = N'COLUMN', @level2name = N'firm_id';
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'Код клиента (QUIK: ClientCode).', @level0type = N'SCHEMA', @level0name = N'ref', @level1type = N'TABLE', @level1name = N'clients', @level2type = N'COLUMN', @level2name = N'code';

GRANT SELECT, INSERT ON ref.clients TO quik_portfolio_writer;
GRANT SELECT ON ref.clients TO quik_portfolio_worker;

-- +goose Down

DROP TABLE ref.clients;
