-- +goose Up

CREATE SCHEMA quik;

CREATE TABLE quik.money_limits (
    load_date     date           NOT NULL DEFAULT (GETDATE()),
    client_code   varchar(12)    NOT NULL,
    currency_code varchar(4)     NOT NULL,
    position_code varchar(4)     NOT NULL,
    settle_code   varchar(5)     NOT NULL,
    firm_code     varchar(12)    NOT NULL,
    firm_name     varchar(128)   NULL,
    balance       decimal(19, 4) NULL,
    source_date   date           NOT NULL DEFAULT (GETDATE()),
    ts            timestamp      NOT NULL,
    CONSTRAINT PK_quik_money_limits PRIMARY KEY CLUSTERED (
        load_date,
        client_code,
        currency_code,
        position_code,
        settle_code,
        firm_code
    )
);

CREATE INDEX idx_quik_money_limits_load_date
    ON quik.money_limits (load_date);

CREATE TABLE quik.security_limits (
    load_date                 date           NOT NULL DEFAULT (GETDATE()),
    client_code               varchar(12)    NOT NULL,
    sec_code                  varchar(12)    NOT NULL, -- изменилось
    trade_account             varchar(12)    NOT NULL,
    settle_code               varchar(5)     NOT NULL DEFAULT 'Tx',
    firm_code                 varchar(12)    NOT NULL,
    firm_name                 varchar(128)   NULL,
    balance                   decimal(19, 4) NULL,
    acquisition_currency_code varchar(4)     NULL, -- изменилось
    isin                      varchar(12)    NULL,
    sec_name                  varchar(128)   NULL,
    source_date               date           NOT NULL DEFAULT (GETDATE()),
    ts                        timestamp      NOT NULL,
    CONSTRAINT PK_quik_security_limits PRIMARY KEY CLUSTERED (
        load_date,
        client_code,
        sec_code,
        trade_account,
        settle_code,
        firm_code
    )
);

CREATE TABLE quik.security_limits_otc (
    load_date                 date           NOT NULL DEFAULT (GETDATE()),
    client_code               varchar(12)    NOT NULL,
    sec_code                  varchar(12)    NOT NULL, -- изменилось
    trade_account             varchar(12)    NOT NULL DEFAULT 'OTC',
    settle_code               varchar(5)     NOT NULL DEFAULT 'Tx',
    firm_code                 varchar(12)    NOT NULL,
    firm_name                 varchar(128)   NULL,
    balance                   decimal(19, 4) NULL,
    acquisition_currency_code varchar(4)     NULL, -- изменилось
    isin                      varchar(12)    NULL,
    sec_name                  varchar(128)   NULL,
    source_date               date           NOT NULL DEFAULT (GETDATE()),
    ts                        timestamp      NOT NULL,
    CONSTRAINT PK_quik_security_limits_otc PRIMARY KEY CLUSTERED (
        load_date,
        client_code,
        sec_code,
        trade_account,
        settle_code,
        firm_code
    )
);

CREATE TABLE quik.current_quotes (
    quote_date           date           NOT NULL DEFAULT (GETDATE()),
    instrument_class     varchar(80)    NOT NULL,
    sec_code             varchar(12)    NULL,
    registration_number  varchar(30)    NULL,
    full_name            varchar(100)   NULL,
    short_name           varchar(50)    NULL,
    class_code           varchar(12)    NULL,
    class_name           varchar(60)    NULL,
    instrument_type      varchar(15)    NULL,
    instrument_subtype   varchar(60)    NULL,
    isin                 varchar(12)    NULL,
    face_value           decimal(19, 8) NULL,
    currency             varchar(4)     NULL,
    base_currency        varchar(4)     NULL,
    quote_currency       varchar(4)     NULL,
    counter_currency     varchar(4)     NULL,
    maturity_date        date           NULL,
    coupon_duration      int            NULL,
    last_price           decimal(19, 8) NULL, -- Цена последней сделки
    close_price          decimal(19, 8) NULL,
    waprice              decimal(19, 8) NULL, -- Средневзвешенная цена
    accrued_int          decimal(19, 8) NULL, -- Накопленный купонный доход
    trading_status       varchar(15)    NULL,
    rw                   rowversion     NOT NULL,
    instrument_id        bigint         NULL
);

CREATE NONCLUSTERED INDEX IX_current_quotes_instrument_class
    ON quik.current_quotes (instrument_class);

CREATE NONCLUSTERED INDEX IX_current_quotes_instrument_id_null
    ON quik.current_quotes (instrument_id)
    WHERE instrument_id IS NULL;
