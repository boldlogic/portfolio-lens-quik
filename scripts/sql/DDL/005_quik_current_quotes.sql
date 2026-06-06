IF NOT EXISTS (SELECT 1 FROM sys.schemas WHERE name = N'quik')
    EXEC (N'CREATE SCHEMA quik');
GO

IF OBJECT_ID(N'quik.current_quotes', N'U') IS NULL
BEGIN

CREATE TABLE quik.current_quotes (
        quote_date date NOT NULL default (getdate ()),
        instrument_class     varchar(80) NOT NULL,
        sec_code               varchar(12)   NULL,
        registration_number  varchar(30)   NULL,
        full_name            varchar(100)  NULL,
        short_name           varchar(50)   NULL,
        class_code            varchar(12)   NULL,
        class_name            varchar(60)   NULL,
        instrument_type       varchar(15)   NULL,
        instrument_subtype    varchar(60)   NULL,
        isin                  varchar(12)   NULL,
        face_value            DECIMAL(19,8)      NULL,
        currency         varchar(4)    NULL,
        base_currency         varchar(4)    NULL,
        quote_currency        varchar(4)    NULL,
        counter_currency      varchar(4)    NULL,
        maturity_date         DATE       NULL,
        coupon_duration       INT        NULL,
        ---
        last_price DECIMAL(19,8)      NULL, --Цена последней сделки 
        close_price DECIMAL(19,8)      NULL,
        waprice DECIMAL(19,8) null,--Средневзвешенная цена
        accrued_int DECIMAL(19,8) null,--Накопленный купонный доход
        trading_status  varchar(15)   NULL,
        ---
        rw                    ROWVERSION NOT NULL,
        instrument_id         BIGINT     NULL,
    );
END
GO

IF OBJECT_ID(N'quik.current_quotes', N'U') IS NOT NULL
BEGIN
    DROP INDEX IF EXISTS IX_current_quotes_instrument_class
    ON quik.current_quotes;

    CREATE NONCLUSTERED INDEX IX_current_quotes_instrument_class
    ON quik.current_quotes (instrument_class);
END
GO

IF OBJECT_ID(N'quik.current_quotes', N'U') IS NOT NULL
BEGIN
    DROP INDEX IF EXISTS IX_current_quotes_instrument_id_null
    ON quik.current_quotes;

    CREATE NONCLUSTERED INDEX IX_current_quotes_instrument_id_null
    ON quik.current_quotes (instrument_id)
    WHERE instrument_id IS NULL;
END
GO
