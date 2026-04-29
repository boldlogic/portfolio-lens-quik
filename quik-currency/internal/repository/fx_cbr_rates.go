package repository

import (
	"context"
	"time"

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
LEFT JOIN quik.currencies c ON RTRIM(c.iso_char_code) = qc.iso_char_code
WHERE qc.iso_char_code IS NOT NULL
	AND qc.iso_char_code NOT IN ('RUR', 'SUR', 'RUB')
	AND c.iso_code IS NULL
GROUP BY qc.quote_date, qc.iso_char_code;`
	mergeFxCBRRatesQuik = `
WITH qc_raw AS (
	SELECT
		q.quote_date,
		iso_char_code = UPPER(COALESCE(
			NULLIF(LTRIM(RTRIM(q.currency)), ''),
			NULLIF(LTRIM(RTRIM(q.ticker)), '')
		)),
		rate = COALESCE(q.close_price, q.last_price)
	FROM quik.current_quotes q
	WHERE q.class_code = 'CROSSRATE'
		AND COALESCE(q.close_price, q.last_price) IS NOT NULL
),
qc AS (
	SELECT
		qc_raw.quote_date,
		qc_raw.iso_char_code,
		rate_quote_per_base = CAST(MAX(qc_raw.rate) AS DECIMAL(18,8)),
		rate_base_per_quote = CAST(1 / NULLIF(MAX(qc_raw.rate), 0) AS DECIMAL(18,8))
	FROM qc_raw
	WHERE qc_raw.iso_char_code IS NOT NULL
		AND qc_raw.iso_char_code NOT IN ('RUR', 'SUR', 'RUB')
	GROUP BY qc_raw.quote_date, qc_raw.iso_char_code
),
src AS (
	SELECT
		date = qc.quote_date,
		base_iso_code = c.iso_code,
		qc.rate_quote_per_base,
		qc.rate_base_per_quote
	FROM qc
	INNER JOIN quik.currencies c ON RTRIM(c.iso_char_code) = qc.iso_char_code
)
MERGE INTO quik.fx_cbr_rates AS tgt
USING src ON tgt.date = src.date
	AND tgt.quote_iso_code = 643
	AND tgt.base_iso_code = src.base_iso_code
WHEN MATCHED
	AND (
		tgt.rate_quote_per_base <> src.rate_quote_per_base
		OR tgt.rate_base_per_quote <> src.rate_base_per_quote
	)
THEN UPDATE SET
	tgt.rate_quote_per_base = src.rate_quote_per_base,
	tgt.rate_base_per_quote = src.rate_base_per_quote,
	tgt.updated_at = SYSDATETIMEOFFSET()
WHEN NOT MATCHED BY TARGET THEN INSERT (
	date,
	quote_iso_code,
	base_iso_code,
	rate_quote_per_base,
	rate_base_per_quote,
	created_at,
	updated_at
) VALUES (
	src.date,
	643,
	src.base_iso_code,
	src.rate_quote_per_base,
	src.rate_base_per_quote,
	SYSDATETIMEOFFSET(),
	SYSDATETIMEOFFSET()
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
		if r.isShutdown(err) {
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
		if r.isShutdown(err) {
			return err
		}
		r.Logger.Error("ошибка получения несопоставленных валют кросс-курсов QUIK", zap.Error(err))
		return models.ErrRetrievingData
	}
	defer rows.Close()

	for rows.Next() {
		var currency unknownFxCBRRateQuikCurrency
		if err = rows.Scan(&currency.quoteDate, &currency.isoCharCode); err != nil {
			if r.isShutdown(err) {
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
		if r.isShutdown(err) {
			return err
		}
		r.Logger.Error("ошибка обхода несопоставленных валют кросс-курсов QUIK", zap.Error(err))
		return models.ErrRetrievingData
	}

	return nil
}
