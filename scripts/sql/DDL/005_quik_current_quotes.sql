IF NOT EXISTS (SELECT 1 FROM sys.schemas WHERE name = N'quik')
    EXEC (N'CREATE SCHEMA quik');
GO

IF OBJECT_ID(N'quik.current_quotes', N'U') IS NULL
BEGIN
    CREATE TABLE quik.current_quotes (
        quote_date date NOT NULL default (getdate ()),
        instrument_class     CHAR(80) NOT NULL,
        ticker               CHAR(12)   NULL,
        registration_number  CHAR(30)   NULL,
        full_name            CHAR(100)  NULL,
        short_name           CHAR(50)   NULL,
        class_code            CHAR(15)   NULL,
        class_name            CHAR(60)   NULL,
        instrument_type       CHAR(15)   NULL,
        instrument_subtype    CHAR(60)   NULL,
        isin                  CHAR(12)   NULL,
        face_value            FLOAT      NULL,
        currency         CHAR(4)    NULL,
        base_currency         CHAR(4)    NULL,
        quote_currency        CHAR(4)    NULL,
        counter_currency      CHAR(4)    NULL,
        maturity_date         DATE       NULL,
        coupon_duration       INT        NULL,
        rw                    ROWVERSION NOT NULL,
        instrument_id         BIGINT     NULL,
        ---
        last_price FLOAT      NULL,
        close_price FLOAT      NULL,
        accrued_int float null,
        trading_status  CHAR(15)   NULL,
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
