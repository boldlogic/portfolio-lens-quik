package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"go.uber.org/zap"
)

const (
	selectMoneyPositions = `
		WITH cte AS (
			SELECT
				li.load_date,
				li.source_date,
				li.client_code,
				currency_code = case when UPPER(TRIM(li.ccy)) in ('SUR', 'RUR') THEN 'RUB' ELSE UPPER(TRIM(li.ccy)) END,
				li.settle_code,
				li.firm_code,
				li.firm_name,
				li.balance,
				settle_max = MAX(li.settle_code) OVER (
					PARTITION BY li.load_date,
					li.client_code,
					li.ccy,
					li.position_code,
					li.firm_code
				)
			FROM
				quik.money_limits li
			WHERE
				li.load_date = cast(@p1 as date)
		)
		SELECT
			c.load_date,
			c.source_date,
			c.client_code,
			c.firm_code,
			c.firm_name,
			c.currency_code,
			currency_name=COALESCE(cv_ec.currency_name,   cv_iso.currency_name),
			c.balance,
			market_value=round(c.balance  *f.rate/ ft.rate,isnull(ct.minor_units, 2))
			FROM cte c
		LEFT JOIN dbo.external_codes ec_ccy 
			ON ec_ccy.ext_system_id = (select
				ext_system_id
			from
				dbo.external_systems
			where
				ext_system = 'QUIK')
			AND ec_ccy.ext_code_type_id = 1
			AND ec_ccy.ext_code = c.currency_code
		LEFT JOIN currencies cv_ec  ON cv_ec.iso_code  = ec_ccy.internal_id
		LEFT JOIN currencies cv_iso ON cv_iso.iso_char_code = c.currency_code
		JOIN currencies ct ON ct.iso_char_code = @p2
		CROSS APPLY dbo.fnFxRateToRub(ISNULL(COALESCE(cv_ec.iso_char_code,   cv_iso.iso_char_code), ''), c.load_date) f
		CROSS APPLY dbo.fnFxRateToRub(@p2,            c.load_date) ft
		WHERE c.settle_code = c.settle_max	
	`
)

func scanMoneyToPosition(row *sql.Rows) (moneyPosition, error) {
	var out moneyPosition
	err := row.Scan(
		&out.LoadDate,
		&out.SourceDate,
		&out.ClientCode,
		&out.FirmCode,
		&out.FirmName,
		&out.CurrencyCode,
		&out.CurrencyName,
		&out.Balance,
		&out.MV,
	)
	if err != nil {
		return moneyPosition{}, err
	}

	return out, nil
}

func (r *Repository) SelectMoneyPortfolio(ctx context.Context, date time.Time, targetCcy string) (result []quik.Position, err error) {
	start := time.Now()
	defer func() { r.metrics.ObserveRepository("SelectMoneyPortfolio", time.Since(start), err) }()

	pos, err := selectRows(
		ctx,
		r.Db,
		selectMoneyPositions,
		scanMoneyToPosition,
		date,
		targetCcy)
	if err != nil {
		r.Logger.Error("", zap.Error(err))

		return nil, err
	}
	result = toPosition(pos)

	return result, nil
}
