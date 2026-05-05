CREATE OR ALTER FUNCTION dbo.fnFxRateToRub (
    @isoCode NVARCHAR(3),
    @date    DATE
)
RETURNS TABLE
AS
RETURN (
    SELECT
        rate      = ISNULL(r.rate_quote_per_base, 1),
        rate_date = r.rate_date
    FROM (SELECT 1 AS d) t
    OUTER APPLY (
        SELECT TOP 1
            fx.rate_quote_per_base,
            fx.[date] AS rate_date
        FROM dbo.fx_cbr_rates fx
        JOIN dbo.currencies c ON c.iso_code = fx.base_iso_code
        WHERE c.iso_char_code   = @isoCode
          AND fx.quote_iso_code = 643
          AND fx.[date]        <= @date
        ORDER BY fx.[date] DESC
    ) r
);
