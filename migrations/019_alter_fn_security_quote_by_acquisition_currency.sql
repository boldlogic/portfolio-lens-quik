-- +goose Up 

-- +goose StatementBegin 

CREATE OR ALTER FUNCTION market.fnSecurityQuoteByAcquisitionCurrency (
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
                    SELECT CAST(accrued.interest / rate AS DECIMAL(19, 8))
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