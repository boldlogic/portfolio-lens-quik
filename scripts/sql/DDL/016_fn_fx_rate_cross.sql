CREATE OR ALTER FUNCTION dbo.fnFxRateCross 
(
    @fromIsoCode varchar(4),
    @toIsoCode varchar(4),
    @date DATE
) 
RETURNS 
TABLE 
AS RETURN 
(
    WITH norm AS (
        SELECT
            from_iso = CASE WHEN UPPER(TRIM(ISNULL(@fromIsoCode, ''))) IN ('SUR', 'RUR') THEN 'RUB' ELSE UPPER(TRIM(ISNULL(@fromIsoCode, ''))) END,
            to_iso = CASE WHEN UPPER(TRIM(ISNULL(@toIsoCode, ''))) IN ('SUR', 'RUR') THEN 'RUB' ELSE UPPER(RTRIM(ISNULL(@toIsoCode, ''))) END
    )
    SELECT
        rate = cast (
            ISNULL(src.rate_quote_per_base, 1.0) / ISNULL(tgt.rate_quote_per_base, 1.0)  AS DECIMAL(18, 8)),
        rate_date = COALESCE(tgt.rate_date, src.rate_date)
		FROM
        norm n
        OUTER APPLY (
            SELECT
                TOP 1 fx.rate_quote_per_base,
                fx.[date] AS rate_date
            FROM
                dbo.fx_cbr_rates fx
                JOIN dbo.currencies c ON c.iso_code = fx.base_iso_code
            WHERE
                n.from_iso <> 'RUB'
                AND c.iso_char_code = n.from_iso
                AND fx.quote_iso_code = 643
                AND fx.[date] <= @date
            ORDER BY
                fx.[date] DESC
        ) q_from
        CROSS APPLY (
            SELECT
                rate_quote_per_base = CASE 
                    WHEN n.from_iso = 'RUB' THEN CAST(1.0 AS DECIMAL(18, 8)) 
                    ELSE ISNULL(q_from.rate_quote_per_base, 1.0) 
                    END,
                rate_date = CASE 
                    WHEN n.from_iso = 'RUB' 
                    THEN @date 
                    ELSE q_from.rate_date 
                    END
        ) src
        OUTER APPLY (
            SELECT
                TOP 1 fx.rate_quote_per_base,
                fx.[date] AS rate_date
            FROM
                dbo.fx_cbr_rates fx
                JOIN dbo.currencies c ON c.iso_code = fx.base_iso_code
            WHERE
                n.to_iso <> 'RUB'
                AND c.iso_char_code = n.to_iso
                AND fx.quote_iso_code = 643
                AND fx.[date] <= @date
            ORDER BY
                fx.[date] DESC
        ) q_to
        CROSS APPLY (
            SELECT
                rate_quote_per_base = CASE WHEN n.to_iso = 'RUB' THEN CAST(1.0 AS DECIMAL(18, 8)) ELSE ISNULL(q_to.rate_quote_per_base, 1.0) END,
                rate_date = CASE WHEN n.to_iso = 'RUB' THEN @date ELSE q_to.rate_date END
        ) tgt
);