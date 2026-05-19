package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"go.uber.org/zap"
)

const (
	selectSecurityLimitsByDate = `
			WITH cte AS (
			SELECT
				li.load_date,
				li.source_date,
				li.client_code,
				li.ticker,
				li.trade_account,
				li.firm_code,
				li.settle_code,
				li.firm_name,
				li.balance,
				li.acquisition_ccy,
				li.isin,
				settle_max = MAX(li.settle_code) OVER (
					PARTITION BY li.load_date, li.client_code, li.ticker, li.trade_account, li.firm_code
				)
			FROM quik.security_limits li
			WHERE li.load_date = cast(@p1 as date)
		)
		SELECT
			cte.load_date,
			cte.source_date,
			cte.client_code,
			cte.ticker,
			cte.settle_code,
			cte.trade_account,
			cte.firm_code,
			cte.firm_name,
			cte.balance,
			cte.acquisition_ccy,
			cte.isin,
			short_disp.short_name
		FROM cte
		OUTER APPLY (
			SELECT TOP 1 ltrim(rtrim(q.short_name)) AS short_name
			FROM quik.current_quotes q
			WHERE q.ticker = cte.ticker
			ORDER BY CASE
				WHEN cte.acquisition_ccy = q.base_currency AND cte.acquisition_ccy = q.counter_currency THEN 0
				WHEN cte.acquisition_ccy = q.base_currency THEN 1
				WHEN cte.acquisition_ccy = q.counter_currency THEN 2
				ELSE 3 END
		) short_disp
		WHERE 1=1
		ORDER BY cte.load_date, cte.client_code, cte.ticker, cte.trade_account, cte.firm_code
`
)

func (r *Repository) SelectSecurityLimits(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error) {
	var result []quik.SecurityLimit

	rows, err := r.Db.QueryContext(ctx, selectSecurityLimitsByDate, date)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("ошибка выполнения запроса на получение позиций по бумагам", zap.Time("load_date", date), zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	defer rows.Close()

	for rows.Next() {
		row := quik.SecurityLimit{}
		var shortName sql.NullString
		err = rows.Scan(
			&row.LoadDate,
			&row.SourceDate,
			&row.ClientCode,
			&row.Ticker,
			&row.SettleCode,
			&row.TradeAccount,
			&row.FirmCode,
			&row.FirmName,
			&row.Balance,
			&row.AcquisitionCcy,
			&row.ISIN,
			&shortName,
		)
		if shortName.Valid {
			row.ShortName = shortName.String
		}
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}
			r.Logger.Error("ошибка чтения позиций по бумагам", zap.Time("load_date", date), zap.Error(err))
			return nil, models.ErrRetrievingData
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		r.Logger.Error("ошибка чтения позиций по бумагам", zap.Time("load_date", date), zap.Error(rows.Err()))
		return nil, models.ErrRetrievingData
	}
	r.Logger.Debug("результаты получения позиций по бумагам", zap.Time("load_date", date), zap.Int("count", len(result)))
	if len(result) == 0 {
		r.Logger.Warn("позиции по бумагам не найдены", zap.Time("load_date", date))
	}
	return result, nil
}
