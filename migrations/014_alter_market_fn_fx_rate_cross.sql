-- +goose Up 

-- +goose StatementBegin 

CREATE OR ALTER FUNCTION market.fnFxRateCross (
    @fromIsoCode VARCHAR(4),
    @toIsoCode   VARCHAR(4),
    @date        DATE
) 
RETURNS TABLE 
AS 
RETURN (
    WITH norm AS (
        SELECT
            from_iso = CASE 
                            WHEN UPPER(ISNULL(@fromIsoCode, '')) IN ('SUR', 'RUR') THEN 'RUB' 
                            WHEN UPPER(ISNULL(@fromIsoCode, '')) = 'USDX' THEN 'USD' 
                            ELSE UPPER(ISNULL(@fromIsoCode, '')) 
                        END,
            to_iso = CASE 
                            WHEN UPPER(ISNULL(@toIsoCode, '')) IN ('SUR', 'RUR') THEN 'RUB' 
                            WHEN UPPER(ISNULL(@toIsoCode, '')) = 'USDX' THEN 'USD' 
                            ELSE UPPER(ISNULL(@toIsoCode, '')) 
                        END
    )
    SELECT
        from_iso_code = f.iso_code,
        to_iso_code = t.iso_code,
        rate = cast(
            isnull(qr.rate_quote_per_base, 1) AS DECIMAL(18, 8)
        ),
        rate_date = qr.rate_date
    FROM
        norm n
        JOIN ref.currencies f ON f.iso_char_code = n.from_iso
        JOIN ref.currencies t ON t.iso_char_code = n.to_iso
        outer apply (
            select
                top 1 
                fx.rate_quote_per_base,
                rate_date = fx.[date]
            FROM
                market.fx_cbr_rates fx
            where
                fx.base_iso_code = t.iso_code
                and fx.quote_iso_code = f.iso_code
                and fx.[date] <= @date
            order by
               fx.[date] desc
        ) qr
);

-- +goose StatementEnd

CREATE NONCLUSTERED INDEX IX_fx_cbr_rates_pair_date
ON market.fx_cbr_rates (quote_iso_code, base_iso_code, [date] DESC)
INCLUDE (rate_quote_per_base);