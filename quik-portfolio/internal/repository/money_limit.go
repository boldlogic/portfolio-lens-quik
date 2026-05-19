package repository

import (
	"context"
	"time"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"go.uber.org/zap"
)

const (
	selectMoneyLimitsByDate = `
		WITH cte AS (
			SELECT
				li.load_date,
				li.source_date,
				li.client_code,
				li.ccy,
				li.position_code,
				li.firm_code,
				li.settle_code,
				li.firm_name,
				li.balance,
				settle_max = MAX(li.settle_code) OVER (
					PARTITION BY li.load_date, li.client_code, li.ccy, li.position_code, li.firm_code
				)
			FROM quik.money_limits li
			WHERE li.load_date = cast(@p1 as date)
		)
		SELECT
			c.load_date,
			c.source_date,
			c.client_code,
			c.ccy,
			c.position_code,
			c.settle_code,
			c.firm_code,
			c.firm_name,
			c.balance
		FROM cte c
		WHERE 1=1
		ORDER BY c.load_date, c.client_code, c.ccy, c.position_code, c.firm_code;
`
)

func (r *Repository) SelectMoneyLimits(ctx context.Context, date time.Time) ([]quik.MoneyLimit, error) {
	var result []quik.MoneyLimit

	rows, err := r.Db.QueryContext(ctx, selectMoneyLimitsByDate, date)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("ошибка выполнения запроса на получение позиций по деньгам", zap.Time("load_date", date), zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	defer rows.Close()

	for rows.Next() {
		row := quik.MoneyLimit{}
		err = rows.Scan(
			&row.LoadDate,
			&row.SourceDate,
			&row.ClientCode,
			&row.Currency,
			&row.PositionCode,
			&row.SettleCode,
			&row.FirmCode,
			&row.FirmName,
			&row.Balance,
		)
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}
			r.Logger.Error("ошибка чтения позиций по деньгам", zap.Time("load_date", date), zap.Error(err))
			return nil, models.ErrRetrievingData
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		r.Logger.Error("ошибка чтения позиций по деньгам", zap.Time("load_date", date), zap.Error(rows.Err()))
		return nil, models.ErrRetrievingData
	}

	r.Logger.Debug("результаты получения позиций по деньгам", zap.Time("load_date", date), zap.Int("count", len(result)))

	if len(result) == 0 {
		r.Logger.Warn("позиции по деньгам не найдены", zap.Time("load_date", date))
	}
	return result, nil
}
