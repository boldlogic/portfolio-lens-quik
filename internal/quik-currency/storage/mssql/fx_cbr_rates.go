package mssql

import "context"

const (
	mergeFxCBRRates = `
		WITH src_quotes AS (
			SELECT
				q.quote_date,
				norm.iso_char_code,
				rate_quote_per_base = CAST(COALESCE(q.close_price, q.last_price, q.crossrate) AS DECIMAL(18, 8)),
				rn = ROW_NUMBER() OVER (
					PARTITION BY q.quote_date, norm.iso_char_code
					ORDER BY q.rw DESC
				)
			FROM quik.current_quotes q
			CROSS APPLY (
				SELECT iso_char_code = UPPER(COALESCE(NULLIF(q.currency, ''), NULLIF(q.sec_code, '')))
			) raw
			CROSS APPLY (
				SELECT iso_char_code = CASE
					WHEN raw.iso_char_code IN ('SUR', 'RUR') THEN 'RUB'
					WHEN raw.iso_char_code = 'USDX' THEN 'USD'
					ELSE raw.iso_char_code
				END
			) norm
			WHERE q.class_code = 'CROSSRATE'
				AND COALESCE(q.close_price, q.last_price, q.crossrate) > 0
		), rub_rates AS (
			SELECT
				date = sq.quote_date,
				quote_iso_code = CAST(643 AS SMALLINT),
				base_iso_code = COALESCE(c_by_ext.iso_code, c_by_iso.iso_code),
				sq.rate_quote_per_base,
				rate_base_per_quote = CAST(1 / sq.rate_quote_per_base AS DECIMAL(18, 8)),
				es.ext_system_id
			FROM src_quotes sq
			JOIN ref.external_systems es
				ON es.ext_system = 'QUIK'
			LEFT JOIN ref.external_codes ec
				ON ec.ext_code_type_id = 1
				AND ec.ext_code = sq.iso_char_code
				AND ec.ext_system_id = es.ext_system_id
			LEFT JOIN ref.currencies c_by_ext
				ON c_by_ext.iso_code = ec.internal_id
			LEFT JOIN ref.currencies c_by_iso
				ON c_by_iso.iso_char_code = sq.iso_char_code
			WHERE sq.rn = 1
				AND COALESCE(c_by_ext.iso_code, c_by_iso.iso_code) <> 643
		), src AS (
			SELECT
				date,
				quote_iso_code,
				base_iso_code,
				rate_quote_per_base,
				rate_base_per_quote,
				ext_system_id
			FROM rub_rates
		
			UNION ALL
		
			SELECT
				date,
				quote_iso_code = base_iso_code,
				base_iso_code = quote_iso_code,
				rate_quote_per_base = rate_base_per_quote,
				rate_base_per_quote = rate_quote_per_base,
				ext_system_id
			FROM rub_rates
		
			UNION ALL
		
			SELECT
				date = base.date,
				quote_iso_code = quote.base_iso_code,
				base_iso_code = base.base_iso_code,
				rate_quote_per_base = CAST(base.rate_quote_per_base / NULLIF(quote.rate_quote_per_base, 0) AS DECIMAL(18, 8)),
				rate_base_per_quote = CAST(quote.rate_quote_per_base / NULLIF(base.rate_quote_per_base, 0) AS DECIMAL(18, 8)),
				base.ext_system_id
			FROM rub_rates base
			JOIN rub_rates quote
				ON quote.date = base.date
			WHERE base.base_iso_code <> quote.base_iso_code
		)
		
			MERGE INTO market.fx_cbr_rates AS tgt
				USING src ON tgt.date = src.date
					AND tgt.quote_iso_code = src.quote_iso_code
					AND tgt.base_iso_code  = src.base_iso_code
				WHEN MATCHED 
				AND (tgt.rate_quote_per_base <> src.rate_quote_per_base
					OR tgt.rate_base_per_quote <> src.rate_base_per_quote)
				THEN UPDATE SET
					tgt.rate_quote_per_base      = src.rate_quote_per_base,
					tgt.rate_base_per_quote = src.rate_base_per_quote,
					tgt.updated_at          = SYSDATETIMEOFFSET(),
					tgt.ext_system_id=src.ext_system_id
				WHEN NOT MATCHED BY TARGET THEN INSERT (
					date, quote_iso_code, base_iso_code,
					rate_quote_per_base, rate_base_per_quote,
					created_at, updated_at,ext_system_id
				) VALUES (
					src.date, src.quote_iso_code, src.base_iso_code,
					src.rate_quote_per_base, src.rate_base_per_quote,
					SYSDATETIMEOFFSET(), SYSDATETIMEOFFSET(), src.ext_system_id
				);`
)

func (repo CurrencyRepo) MergeFxCBRRates(ctx context.Context) error {
	_, err := repo.runner.ExecContext(ctx, mergeFxCBRRates)
	if err != nil {
		return err
	}
	return nil
}

// func (repo CurrencyRepo) SelectFXCBRRates(ctx context.Context)
