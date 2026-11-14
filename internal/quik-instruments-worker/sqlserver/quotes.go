package sqlserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/model"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	"github.com/shopspring/decimal"
)

const selectQuotesSQL = `
WITH type_boards AS (
	SELECT
		b.board_id,
		b.code,
		b.trade_point_id
	FROM ref.boards b
	JOIN ref.instrument_type_boards tb
		ON tb.board_id = b.board_id
	JOIN ref.instrument_types t
		ON t.type_id = tb.type_id
	WHERE t.title = @p1
		AND b.trade_point_id IS NOT NULL
),
src_quotes AS (
	SELECT TOP (@p2)
		q.rw,
		q.sec_code,
		b.trade_point_id,
		b.board_id,
		q.registration_number,
		q.full_name,
		q.short_name,
		q.isin,
		q.face_value,
		q.maturity_date,
		q.coupon_duration,
		currency = norm_currency.currency,
		base_currency = norm_base_currency.base_currency,
		quote_currency = norm_quote_currency.quote_currency,
		counter_currency = norm_counter_currency.counter_currency,
		q.instrument_id
	FROM quik.current_quotes q
	JOIN type_boards b
		ON b.code = q.class_code
	CROSS APPLY (
		SELECT currency = CASE
			WHEN q.currency IN ('SUR', 'RUR') THEN 'RUB'
			WHEN q.currency = 'USDX' THEN 'USD'
			ELSE q.currency
		END
	) norm_currency
	CROSS APPLY (
		SELECT base_currency = CASE
			WHEN q.base_currency IN ('SUR', 'RUR') THEN 'RUB'
			WHEN q.base_currency = 'USDX' THEN 'USD'
			ELSE q.base_currency
		END
	) norm_base_currency
	CROSS APPLY (
		SELECT quote_currency = CASE
			WHEN q.quote_currency IN ('SUR', 'RUR') THEN 'RUB'
			WHEN q.quote_currency = 'USDX' THEN 'USD'
			ELSE q.quote_currency
		END
	) norm_quote_currency
	CROSS APPLY (
		SELECT counter_currency = CASE
			WHEN q.counter_currency IN ('SUR', 'RUR') THEN 'RUB'
			WHEN q.counter_currency = 'USDX' THEN 'USD'
			ELSE q.counter_currency
		END
	) norm_counter_currency
	WHERE CONVERT(binary(8), q.rw) > @p3
		AND q.sec_code IS NOT NULL
		AND LTRIM(RTRIM(q.sec_code)) <> N''
	ORDER BY q.rw
),
currency_keys AS (
	SELECT DISTINCT currency FROM src_quotes WHERE currency IS NOT NULL
	UNION
	SELECT DISTINCT base_currency FROM src_quotes WHERE base_currency IS NOT NULL
	UNION
	SELECT DISTINCT quote_currency FROM src_quotes WHERE quote_currency IS NOT NULL
	UNION
	SELECT DISTINCT counter_currency FROM src_quotes WHERE counter_currency IS NOT NULL
),
currencies AS (
	SELECT
		k.currency,
		iso_code = COALESCE(c_by_ext.iso_code, c_by_iso.iso_code)
	FROM currency_keys k
	LEFT JOIN ref.external_systems es
		ON es.ext_system = 'QUIK'
	LEFT JOIN ref.external_codes ec
		ON ec.ext_system_id = es.ext_system_id
		AND ec.ext_code_type_id = 1
		AND ec.ext_code = k.currency
	LEFT JOIN ref.currencies c_by_ext
		ON c_by_ext.iso_code = ec.internal_id
	LEFT JOIN ref.currencies c_by_iso
		ON c_by_iso.iso_char_code = k.currency
)
SELECT
	s.rw,
	s.sec_code,
	s.trade_point_id,
	s.board_id,
	s.registration_number,
	s.full_name,
	s.short_name,
	s.isin,
	s.face_value,
	s.maturity_date,
	s.coupon_duration,
	currency_id = c1.iso_code,
	base_currency_id = c2.iso_code,
	quote_currency_id = c3.iso_code,
	counter_currency_id = c4.iso_code,
	s.instrument_id
FROM src_quotes s
LEFT JOIN currencies c1 ON c1.currency = s.currency
LEFT JOIN currencies c2 ON c2.currency = s.base_currency
LEFT JOIN currencies c3 ON c3.currency = s.quote_currency
LEFT JOIN currencies c4 ON c4.currency = s.counter_currency
ORDER BY s.rw;`

const countBlankSecCodesSQL = `
SELECT COUNT(*)
FROM quik.current_quotes q
JOIN ref.boards b
	ON b.code = q.class_code
JOIN ref.instrument_type_boards tb
	ON tb.board_id = b.board_id
JOIN ref.instrument_types t
	ON t.type_id = tb.type_id
WHERE t.title = @p1
	AND b.trade_point_id IS NOT NULL
	AND (
		q.sec_code IS NULL
		OR LTRIM(RTRIM(q.sec_code)) = N''
	);`

func (r *Repository) SelectQuotes(ctx context.Context, title quik.InstrumentTypeTitle, limit int, afterRowVersion []byte) ([]model.SourceQuote, error) {
	rows, err := r.Db.QueryContext(ctx, selectQuotesSQL, string(title), limit, rowVersionParam(afterRowVersion))
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения котировок: %w", err)
	}
	defer rows.Close()

	var out []model.SourceQuote
	for rows.Next() {
		item, err := scanSourceQuote(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка обхода котировок: %w", err)
	}
	return out, nil
}

func scanSourceQuote(row *sql.Rows) (model.SourceQuote, error) {
	var out model.SourceQuote
	var registrationNumber sql.NullString
	var fullName sql.NullString
	var shortName sql.NullString
	var isin sql.NullString
	var faceValue *decimal.Decimal
	var maturityDate sql.NullTime
	var couponDuration sql.NullInt32
	var currencyID sql.NullInt16
	var baseCurrencyID sql.NullInt16
	var quoteCurrencyID sql.NullInt16
	var counterCurrencyID sql.NullInt16
	var linkedInstrumentID sql.NullInt64

	err := row.Scan(
		&out.SourceRowVersion,
		&out.SecCode,
		&out.TradePointID,
		&out.BoardID,
		&registrationNumber,
		&fullName,
		&shortName,
		&isin,
		&faceValue,
		&maturityDate,
		&couponDuration,
		&currencyID,
		&baseCurrencyID,
		&quoteCurrencyID,
		&counterCurrencyID,
		&linkedInstrumentID,
	)
	if err != nil {
		return model.SourceQuote{}, fmt.Errorf("ошибка чтения котировки: %w", err)
	}

	out.RegistrationNumber = stringPtr(registrationNumber)
	out.FullName = stringPtr(fullName)
	out.ShortName = stringPtr(shortName)
	out.ISIN = stringPtr(isin)
	out.FaceValue = decimalPtr(faceValue)
	out.MaturityDate = timePtr(maturityDate)
	out.CouponDuration = int32Ptr(couponDuration)
	out.CurrencyID = int16Ptr(currencyID)
	out.BaseCurrencyID = int16Ptr(baseCurrencyID)
	out.QuoteCurrencyID = int16Ptr(quoteCurrencyID)
	out.CounterCurrencyID = int16Ptr(counterCurrencyID)
	out.LinkedInstrumentID = int64Ptr(linkedInstrumentID)
	return out, nil
}

func (r *Repository) CountBlankSecCodes(ctx context.Context, title quik.InstrumentTypeTitle) (int, error) {
	var count int
	err := r.Db.QueryRowContext(ctx, countBlankSecCodesSQL, string(title)).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка подсчёта котировок без кода: %w", err)
	}
	return count, nil
}

func rowVersionParam(rowVersion []byte) []byte {
	if len(rowVersion) == 0 {
		return []byte{0, 0, 0, 0, 0, 0, 0, 0}
	}
	return rowVersion
}

func int64Ptr(v sql.NullInt64) *int64 {
	if !v.Valid {
		return nil
	}
	return &v.Int64
}
