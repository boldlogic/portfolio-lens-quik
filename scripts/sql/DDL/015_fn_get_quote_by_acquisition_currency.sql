CREATE OR ALTER FUNCTION dbo.fnGetQuoteByAcquisitionCurrency 
(
@ticker CHAR(12), 
@currency CHAR(4)
) 
RETURNS 
TABLE 
AS RETURN (
    SELECT
        TOP 1 market_value = ROUND(instr.price + ai_instr.ai, 4),
        price = ROUND(instr.price, 4),
        ai = ROUND(ai_instr.ai, 4),
        short_name=trim(q.short_name),
        q.quote_date,
        instr.currency
    FROM
        quik.current_quotes q
        outer apply (
            select
                interest = isnull(q.accrued_int, 0),
                currency = case when q.instrument_type = 'Облигации' then isnull(q.counter_currency, q.currency) else null end
        ) accrued
        outer apply (
            select
                price = case when q.instrument_type = 'Облигации' then (isnull(q.face_value, 0) / 100.0) * (
                    case when isnull(q.last_price, 0) <> 0 then q.last_price else isnull(q.close_price, 0) end
                ) else (
                    case when isnull(q.last_price, 0) <> 0 then q.last_price else q.close_price end
                ) end,
                currency = case when q.instrument_type = 'Облигации' then coalesce(
                    nullif(trim(q.base_currency), ''),
                    nullif(trim(q.quote_currency), ''),
                    nullif(trim(q.currency), '')
                ) else isnull(q.counter_currency, q.base_currency) end
        ) instr
        outer apply (
            select
                ai = case when q.instrument_type = 'Облигации'
                and accrued.currency <> instr.currency then (
                    select
                        accrued.interest * rate
                    from
                        dbo.fnFxRateCross (accrued.currency, instr.currency, q.quote_date)
                ) else accrued.interest end
        ) ai_instr
    WHERE
        q.ticker = @ticker
    ORDER BY
        case when @currency = q.base_currency
        and @currency = q.counter_currency then 0 when @currency = q.base_currency then 1 when @currency = q.counter_currency then 2 else 3 end
);