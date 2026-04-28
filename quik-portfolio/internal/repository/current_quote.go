package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/boldlogic/packages/shutdown"
	qmodels "github.com/boldlogic/portfolio-lens-quik/quik-portfolio/internal/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type CurrentQuote struct {
	Ticker          string
	ISIN            *string
	LastPrice       *decimal.Decimal
	ClosePrice      *decimal.Decimal
	AccruedInt      *decimal.Decimal
	FaceValue       *decimal.Decimal
	InstrumentType  string
	PriceCurrency   string
	AccruedCurrency string
	QuoteDate       time.Time
}

const selectCurrentQuotes = `
SELECT
    q.ticker,
    q.isin,
    q.last_price,
    q.close_price,
    q.accrued_int,
    q.face_value,
    q.instrument_type,
    CASE WHEN q.instrument_type = N'Облигации'
         THEN COALESCE(
                NULLIF(LTRIM(RTRIM(q.base_currency)), N''),
                NULLIF(LTRIM(RTRIM(q.quote_currency)), N''),
                NULLIF(LTRIM(RTRIM(q.currency)), N'')
              )
         ELSE ISNULL(q.counter_currency, q.base_currency)
    END as price_currency,
    CASE WHEN q.instrument_type = N'Облигации'
         THEN ISNULL(q.counter_currency, q.currency)
         ELSE ISNULL(q.base_currency, q.currency)
    END as accrued_currency,
    q.quote_date
FROM quik.current_quotes q
WHERE q.quote_date = CAST(GETDATE() AS DATE)
`

func (r *Repository) SelectCurrentQuotes(ctx context.Context) ([]qmodels.CurrentQuote, error) {
	rows, err := r.Db.QueryContext(ctx, selectCurrentQuotes)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("failed to query current quotes", zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	defer rows.Close()

	var result []qmodels.CurrentQuote
	for rows.Next() {
		var q qmodels.CurrentQuote
		var ticker, isin sql.NullString
		var lastPrice, closePrice, accruedInt, faceValue sql.NullFloat64
		var instrumentType, priceCurrency, accruedCurrency sql.NullString
		var quoteDate sql.NullTime

		err := rows.Scan(&ticker, &isin, &lastPrice, &closePrice, &accruedInt, &faceValue, &instrumentType, &priceCurrency, &accruedCurrency, &quoteDate)
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}
			r.Logger.Error("failed to scan current quote", zap.Error(err))
			return nil, models.ErrRetrievingData
		}

		if ticker.Valid {
			q.Ticker = ticker.String
		}
		if isin.Valid && isin.String != "" {
			q.ISIN = &isin.String
		}
		if lastPrice.Valid {
			d := decimal.NewFromFloat(lastPrice.Float64)
			q.LastPrice = &d
		}
		if closePrice.Valid {
			d := decimal.NewFromFloat(closePrice.Float64)
			q.ClosePrice = &d
		}
		if accruedInt.Valid {
			d := decimal.NewFromFloat(accruedInt.Float64)
			q.AccruedInt = &d
		}
		if faceValue.Valid {
			d := decimal.NewFromFloat(faceValue.Float64)
			q.FaceValue = &d
		}
		if instrumentType.Valid {
			q.InstrumentType = instrumentType.String
		}
		if priceCurrency.Valid {
			q.PriceCurrency = priceCurrency.String
		}
		if accruedCurrency.Valid {
			q.AccruedCurrency = accruedCurrency.String
		}
		if quoteDate.Valid {
			q.QuoteDate = quoteDate.Time
		} else {
			q.QuoteDate = time.Now()
		}

		result = append(result, q)
	}

	if err := rows.Err(); err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("error iterating current quotes", zap.Error(err))
		return nil, models.ErrRetrievingData
	}

	r.Logger.Debug("current quotes retrieved", zap.Int("count", len(result)))
	return result, nil
}

const quoteKeysBatchSize = 400

// SelectCurrentQuotesForKeys выбирает котировки за сегодня только по списку тикеров/ISIN (верхний регистр).
func (r *Repository) SelectCurrentQuotesForKeys(ctx context.Context, keys []string) ([]qmodels.CurrentQuote, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	seen := make(map[string]struct{}, len(keys))
	norm := make([]string, 0, len(keys))
	for _, k := range keys {
		k = strings.ToUpper(strings.TrimSpace(k))
		if k == "" {
			continue
		}
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		norm = append(norm, k)
	}
	if len(norm) == 0 {
		return nil, nil
	}

	var all []qmodels.CurrentQuote
	for i := 0; i < len(norm); i += quoteKeysBatchSize {
		end := i + quoteKeysBatchSize
		if end > len(norm) {
			end = len(norm)
		}
		batch := norm[i:end]
		part, err := r.selectCurrentQuotesForKeysBatch(ctx, batch)
		if err != nil {
			return nil, err
		}
		all = append(all, part...)
	}
	r.Logger.Debug("current quotes for keys retrieved", zap.Int("keys", len(norm)), zap.Int("rows", len(all)))
	return all, nil
}

func (r *Repository) selectCurrentQuotesForKeysBatch(ctx context.Context, keys []string) ([]qmodels.CurrentQuote, error) {
	n := len(keys)
	tPH := make([]string, n)
	iPH := make([]string, n)
	for j := 0; j < n; j++ {
		tPH[j] = fmt.Sprintf("@p%d", j+1)
		iPH[j] = fmt.Sprintf("@p%d", n+j+1)
	}
	q := fmt.Sprintf(`
SELECT
    q.ticker,
    q.isin,
    q.last_price,
    q.close_price,
    q.accrued_int,
    q.face_value,
    q.instrument_type,
    CASE WHEN q.instrument_type = N'Облигации'
         THEN COALESCE(
                NULLIF(LTRIM(RTRIM(q.base_currency)), N''),
                NULLIF(LTRIM(RTRIM(q.quote_currency)), N''),
                NULLIF(LTRIM(RTRIM(q.currency)), N'')
              )
         ELSE ISNULL(q.counter_currency, q.base_currency)
    END as price_currency,
    CASE WHEN q.instrument_type = N'Облигации'
         THEN ISNULL(q.counter_currency, q.currency)
         ELSE ISNULL(q.base_currency, q.currency)
    END as accrued_currency,
    q.quote_date
FROM quik.current_quotes q
WHERE q.quote_date = CAST(GETDATE() AS DATE)
AND (
  UPPER(LTRIM(RTRIM(ISNULL(q.ticker, N'')))) IN (%s)
  OR UPPER(LTRIM(RTRIM(ISNULL(q.isin, N'')))) IN (%s)
)
`, strings.Join(tPH, ", "), strings.Join(iPH, ", "))

	args := make([]any, 0, 2*n)
	for _, k := range keys {
		args = append(args, k)
	}
	for _, k := range keys {
		args = append(args, k)
	}

	rows, err := r.Db.QueryContext(ctx, q, args...)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("failed to query current quotes for keys", zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	defer rows.Close()

	var result []qmodels.CurrentQuote
	for rows.Next() {
		var q qmodels.CurrentQuote
		var ticker, isin sql.NullString
		var lastPrice, closePrice, accruedInt, faceValue sql.NullFloat64
		var instrumentType, priceCurrency, accruedCurrency sql.NullString
		var quoteDate sql.NullTime

		err := rows.Scan(&ticker, &isin, &lastPrice, &closePrice, &accruedInt, &faceValue, &instrumentType, &priceCurrency, &accruedCurrency, &quoteDate)
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}
			r.Logger.Error("failed to scan current quote", zap.Error(err))
			return nil, models.ErrRetrievingData
		}

		if ticker.Valid {
			q.Ticker = ticker.String
		}
		if isin.Valid && isin.String != "" {
			q.ISIN = &isin.String
		}
		if lastPrice.Valid {
			d := decimal.NewFromFloat(lastPrice.Float64)
			q.LastPrice = &d
		}
		if closePrice.Valid {
			d := decimal.NewFromFloat(closePrice.Float64)
			q.ClosePrice = &d
		}
		if accruedInt.Valid {
			d := decimal.NewFromFloat(accruedInt.Float64)
			q.AccruedInt = &d
		}
		if faceValue.Valid {
			d := decimal.NewFromFloat(faceValue.Float64)
			q.FaceValue = &d
		}
		if instrumentType.Valid {
			q.InstrumentType = instrumentType.String
		}
		if priceCurrency.Valid {
			q.PriceCurrency = priceCurrency.String
		}
		if accruedCurrency.Valid {
			q.AccruedCurrency = accruedCurrency.String
		}
		if quoteDate.Valid {
			q.QuoteDate = quoteDate.Time
		} else {
			q.QuoteDate = time.Now()
		}

		result = append(result, q)
	}

	if err := rows.Err(); err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("error iterating current quotes for keys", zap.Error(err))
		return nil, models.ErrRetrievingData
	}

	return result, nil
}
