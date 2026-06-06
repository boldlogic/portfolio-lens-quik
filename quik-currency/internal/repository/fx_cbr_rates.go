package repository

import (
	"context"
	"time"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"go.uber.org/zap"
)

const (
	selectUnknownFxCBRRatesQuikCurrencies = `
WITH qc AS (
	SELECT
		q.quote_date,
		iso_char_code = UPPER(COALESCE(
			NULLIF(LTRIM(RTRIM(q.currency)), ''),
			NULLIF(LTRIM(RTRIM(q.ticker)), '')
		))
	FROM quik.current_quotes q
	WHERE q.class_code = 'CROSSRATE'
		AND COALESCE(q.close_price, q.last_price) IS NOT NULL
)
SELECT
	qc.quote_date,
	qc.iso_char_code
FROM qc
LEFT JOIN dbo.currencies c ON RTRIM(c.iso_char_code) = qc.iso_char_code
WHERE qc.iso_char_code IS NOT NULL
	AND qc.iso_char_code NOT IN ('RUR', 'SUR', 'RUB')
	AND c.iso_code IS NULL
GROUP BY qc.quote_date, qc.iso_char_code;`
	
mergeFxCBRRatesQuik = `
	WITH qc AS (
		SELECT
			q.quote_date,
			iso_char_code = RTRIM(COALESCE(q.currency, q.ticker)),
			rate_quote_per_base = COALESCE(q.close_price, q.last_price),
			rate_base_per_quote = 1 / NULLIF(COALESCE(q.close_price, q.last_price), 0)

		FROM quik.current_quotes q
		WHERE q.class_code = 'CROSSRATE'
		AND COALESCE(q.close_price, q.last_price) IS NOT NULL
		AND RTRIM(COALESCE(q.currency, q.ticker)) NOT IN ('RUR','SUR','RUB')
	), src as (

	SELECT
		date=qc.quote_date,
		base_iso_code=c.iso_code,
			rate_quote_per_base=CAST(qc.rate_quote_per_base AS DECIMAL(18,8)),
			rate_base_per_quote=CAST(qc.rate_base_per_quote AS DECIMAL(18,8)),
			es.ext_system_id
	FROM qc

	LEFT JOIN dbo.external_systems es 
		ON es.ext_system='QUIK'

	LEFT JOIN dbo.external_codes ec
		ON ec.ext_code_type_id = 1
		AND ec.ext_code = qc.iso_char_code
		AND es.ext_system_id=ec.ext_system_id

	LEFT JOIN dbo.currencies c1
		ON c1.iso_code = ec.internal_id

	LEFT JOIN dbo.currencies c2
		ON c2.iso_char_code = qc.iso_char_code

	CROSS APPLY (
		SELECT iso_code = COALESCE(c1.iso_code, c2.iso_code)
	) x

	JOIN dbo.currencies c
		ON c.iso_code = x.iso_code)

	MERGE INTO dbo.fx_cbr_rates AS tgt
		USING src ON tgt.date = src.date
			AND tgt.quote_iso_code = 643
			AND tgt.base_iso_code  = src.base_iso_code
		WHEN MATCHED 
		AND tgt.rate_quote_per_base  <> src.rate_quote_per_base
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
			src.date, 643, src.base_iso_code,
			src.rate_quote_per_base, src.rate_base_per_quote,
			SYSDATETIMEOFFSET(), SYSDATETIMEOFFSET(), src.ext_system_id
		);`
)

type unknownFxCBRRateQuikCurrency struct {
	quoteDate   time.Time
	isoCharCode string
}

func (r *Repository) MergeFxCBRRatesQuik(ctx context.Context) error {
	if err := r.logUnknownFxCBRRateQuikCurrencies(ctx); err != nil {
		return err
	}

	_, err := r.Db.ExecContext(ctx, mergeFxCBRRatesQuik)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.Logger.Error("ошибка сохранения кросс-курсов валют из QUIK", zap.Error(err))
		return models.ErrSavingData
	}

	return nil
}

func (r *Repository) logUnknownFxCBRRateQuikCurrencies(ctx context.Context) error {
	rows, err := r.Db.QueryContext(ctx, selectUnknownFxCBRRatesQuikCurrencies)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.Logger.Error("ошибка получения несопоставленных валют кросс-курсов QUIK", zap.Error(err))
		return models.ErrRetrievingData
	}
	defer rows.Close()

	for rows.Next() {
		var currency unknownFxCBRRateQuikCurrency
		if err = rows.Scan(&currency.quoteDate, &currency.isoCharCode); err != nil {
			if shutdown.IsExceeded(err) {
				return err
			}
			r.Logger.Error("ошибка чтения несопоставленной валюты кросс-курса QUIK", zap.Error(err))
			return models.ErrRetrievingData
		}
		r.Logger.Warn("кросс-курс QUIK пропущен: валюта отсутствует в справочнике currencies",
			zap.Time("quote_date", currency.quoteDate),
			zap.String("iso_char_code", currency.isoCharCode))
	}

	if err = rows.Err(); err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.Logger.Error("ошибка обхода несопоставленных валют кросс-курсов QUIK", zap.Error(err))
		return models.ErrRetrievingData
	}

	return nil
}
