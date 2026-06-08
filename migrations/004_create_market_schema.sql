-- +goose Up

CREATE SCHEMA market;

CREATE TABLE market.fx_cbr_rates (
    date               DATE              NOT NULL,
    quote_iso_code     SMALLINT          NOT NULL,
    base_iso_code      SMALLINT          NOT NULL,
    rate_quote_per_base DECIMAL(18, 8)   NULL,
    rate_base_per_quote DECIMAL(18, 8)   NULL,
    created_at         DATETIMEOFFSET(7) NOT NULL DEFAULT SYSDATETIMEOFFSET(),
    updated_at         DATETIMEOFFSET(7) NOT NULL DEFAULT SYSDATETIMEOFFSET(),
    ext_system_id      TINYINT           NULL,
    CONSTRAINT PK_market_fx_cbr_rates PRIMARY KEY CLUSTERED (date, quote_iso_code, base_iso_code),
    CONSTRAINT FK_market_fx_cbr_rates_quote_iso FOREIGN KEY (quote_iso_code)
        REFERENCES ref.currencies (iso_code),
    CONSTRAINT FK_market_fx_cbr_rates_base_iso FOREIGN KEY (base_iso_code)
        REFERENCES ref.currencies (iso_code),
    CONSTRAINT FK_market_fx_cbr_rates_ext_system FOREIGN KEY (ext_system_id)
        REFERENCES ref.external_systems (ext_system_id)
);

-- +goose StatementBegin
CREATE FUNCTION market.fnFxRateCross (
    @fromIsoCode VARCHAR(4),
    @toIsoCode   VARCHAR(4),
    @date        DATE
)
RETURNS TABLE
AS
RETURN (
    WITH norm AS (
        SELECT
            from_iso = CASE WHEN UPPER(TRIM(ISNULL(@fromIsoCode, ''))) IN ('SUR', 'RUR') THEN 'RUB' ELSE UPPER(TRIM(ISNULL(@fromIsoCode, ''))) END,
            to_iso   = CASE WHEN UPPER(TRIM(ISNULL(@toIsoCode, ''))) IN ('SUR', 'RUR') THEN 'RUB' ELSE UPPER(RTRIM(ISNULL(@toIsoCode, ''))) END
    )
    SELECT
        rate = CAST(
            ISNULL(src.rate_quote_per_base, 1.0) / ISNULL(tgt.rate_quote_per_base, 1.0) AS DECIMAL(18, 8)
        ),
        rate_date = COALESCE(tgt.rate_date, src.rate_date)
    FROM norm n
    OUTER APPLY (
        SELECT TOP 1
            fx.rate_quote_per_base,
            fx.[date] AS rate_date
        FROM market.fx_cbr_rates fx
        JOIN ref.currencies c ON c.iso_code = fx.base_iso_code
        WHERE n.from_iso <> 'RUB'
            AND c.iso_char_code = n.from_iso
            AND fx.quote_iso_code = 643
            AND fx.[date] <= @date
        ORDER BY fx.[date] DESC
    ) q_from
    CROSS APPLY (
        SELECT
            rate_quote_per_base = CASE
                WHEN n.from_iso = 'RUB' THEN CAST(1.0 AS DECIMAL(18, 8))
                ELSE ISNULL(q_from.rate_quote_per_base, 1.0)
            END,
            rate_date = CASE
                WHEN n.from_iso = 'RUB' THEN @date
                ELSE q_from.rate_date
            END
    ) src
    OUTER APPLY (
        SELECT TOP 1
            fx.rate_quote_per_base,
            fx.[date] AS rate_date
        FROM market.fx_cbr_rates fx
        JOIN ref.currencies c ON c.iso_code = fx.base_iso_code
        WHERE n.to_iso <> 'RUB'
            AND c.iso_char_code = n.to_iso
            AND fx.quote_iso_code = 643
            AND fx.[date] <= @date
        ORDER BY fx.[date] DESC
    ) q_to
    CROSS APPLY (
        SELECT
            rate_quote_per_base = CASE
                WHEN n.to_iso = 'RUB' THEN CAST(1.0 AS DECIMAL(18, 8))
                ELSE ISNULL(q_to.rate_quote_per_base, 1.0)
            END,
            rate_date = CASE
                WHEN n.to_iso = 'RUB' THEN @date
                ELSE q_to.rate_date
            END
    ) tgt
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION market.fnSecurityQuoteByAcquisitionCurrency (
    @sec_code                  VARCHAR(12),
    @acquisition_currency_code VARCHAR(4)
)
RETURNS TABLE
AS
RETURN (
    SELECT TOP 1
        market_value = CAST(ROUND(instr.price + ai_instr.ai, 4) AS DECIMAL(19, 4)),
        price        = CAST(ROUND(instr.price, 4) AS DECIMAL(19, 4)),
        ai           = CAST(ROUND(ai_instr.ai, 4) AS DECIMAL(19, 4)),
        short_name   = TRIM(q.short_name),
        q.quote_date,
        instr.currency
    FROM quik.current_quotes q
    OUTER APPLY (
        SELECT
            interest = ISNULL(q.accrued_int, CAST(0 AS DECIMAL(19, 8))),
            currency = CASE WHEN q.instrument_type = N'Облигации' THEN ISNULL(q.counter_currency, q.currency) ELSE NULL END
    ) accrued
    OUTER APPLY (
        SELECT
            price = CASE WHEN q.instrument_type = N'Облигации' THEN
                (ISNULL(q.face_value, CAST(0 AS DECIMAL(19, 8))) / CAST(100 AS DECIMAL(19, 8)))
                * CASE
                    WHEN ISNULL(q.last_price, CAST(0 AS DECIMAL(19, 8))) <> 0 THEN q.last_price
                    ELSE ISNULL(q.close_price, CAST(0 AS DECIMAL(19, 8)))
                END
            ELSE
                CASE
                    WHEN ISNULL(q.last_price, CAST(0 AS DECIMAL(19, 8))) <> 0 THEN q.last_price
                    ELSE q.close_price
                END
            END,
            currency = CASE WHEN q.instrument_type = N'Облигации' THEN COALESCE(
                NULLIF(TRIM(q.base_currency), N''),
                NULLIF(TRIM(q.quote_currency), N''),
                NULLIF(TRIM(q.currency), N'')
            ) ELSE ISNULL(q.counter_currency, q.base_currency) END
    ) instr
    OUTER APPLY (
        SELECT
            ai = CASE WHEN q.instrument_type = N'Облигации'
                AND accrued.currency <> instr.currency THEN (
                    SELECT CAST(accrued.interest * rate AS DECIMAL(19, 8))
                    FROM market.fnFxRateCross(accrued.currency, instr.currency, q.quote_date)
                ) ELSE accrued.interest END
    ) ai_instr
    WHERE q.sec_code = @sec_code
    ORDER BY
        CASE
            WHEN @acquisition_currency_code = q.base_currency
                AND @acquisition_currency_code = q.counter_currency THEN 0
            WHEN @acquisition_currency_code = q.base_currency THEN 1
            WHEN @acquisition_currency_code = q.counter_currency THEN 2
            ELSE 3
        END
);
-- +goose StatementEnd
