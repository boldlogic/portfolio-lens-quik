package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	"github.com/shopspring/decimal"
)

type moneyPortfolioRow struct {
	LoadDate                    time.Time
	SourceDate                  time.Time
	FxRateDate                  sql.NullTime
	ClientCode                  string
	FirmCode                    string
	FirmName                    string
	CurrencyCode                string
	CurrencyName                sql.NullString
	Balance                     decimal.Decimal
	MarketValueInTargetCurrency decimal.Decimal
}

func scanMoneyPortfolioRow(row *sql.Rows) (moneyPortfolioRow, error) {
	var out moneyPortfolioRow
	err := row.Scan(
		&out.LoadDate,
		&out.SourceDate,
		&out.FxRateDate,
		&out.ClientCode,
		&out.FirmCode,
		&out.FirmName,
		&out.CurrencyCode,
		&out.CurrencyName,
		&out.Balance,
		&out.MarketValueInTargetCurrency,
	)
	if err != nil {
		return moneyPortfolioRow{}, fmt.Errorf("%w: %w", ErrScan, err)
	}

	return out, nil
}

const (
	moneyPortfolioSelectColumnsSQL = `
	SELECT
		c.load_date,
		c.source_date,
		fx_rate_date=fx.rate_date,
		c.client_code,
		c.firm_code,
		c.firm_name,
		c.currency_code,
		fx.currency_name,
		balance=cast(isnull(c.balance,0) as decimal(18,2)),
		market_value_in_target_currency=cast(c.balance /fx.rate as decimal(18,2))
		FROM cte c
`

	moneyPortfolioLatestSettleCTEBaseSQL = `
		WITH cte AS (
			SELECT
				li.load_date,
				li.source_date,
				li.client_code,
				currency_code = case when UPPER(TRIM(li.currency_code)) in ('SUR', 'RUR') THEN 'RUB' ELSE UPPER(TRIM(li.currency_code)) END,
				li.settle_code,
				li.firm_code,
				li.firm_name,
				li.balance,
				settle_max = MAX(li.settle_code) OVER (
					PARTITION BY li.load_date,
					li.client_code,
					li.currency_code,
					li.position_code,
					li.firm_code
				)
			FROM
				quik.money_limits li
			WHERE
				li.load_date = cast(@p1 as date)
`
	moneyPortfolioLatestSettleCTEAllSQL = moneyPortfolioLatestSettleCTEBaseSQL + `
		)
`
	moneyPortfolioLatestSettleCTEByClientsSQL = moneyPortfolioLatestSettleCTEBaseSQL + `
				AND EXISTS (SELECT 1 FROM @codes c WHERE c.client_code = li.client_code)
		)
`

	moneyPortfolioFxRatesCTESQL = `
		, fx_keys AS (
			SELECT DISTINCT
				c.currency_code
			FROM cte c
			WHERE c.settle_code = c.settle_max
		)
		, fx AS (
			SELECT
				k.currency_code,
				currency_name = COALESCE(cv_ec.currency_name, cv_iso.currency_name),
				cr.rate,
				cr.rate_date
			FROM fx_keys k
		LEFT JOIN ref.external_codes ec_ccy 
			ON ec_ccy.ext_system_id = (select
				ext_system_id
			from
				ref.external_systems
			where
				ext_system = 'QUIK')
			AND ec_ccy.ext_code_type_id = 1
				AND ec_ccy.ext_code = k.currency_code
		LEFT JOIN ref.currencies cv_ec  ON cv_ec.iso_code  = ec_ccy.internal_id
			LEFT JOIN ref.currencies cv_iso ON cv_iso.iso_char_code = k.currency_code
		cross apply market.fnFxRateCross(ISNULL(COALESCE(cv_ec.iso_char_code,cv_iso.iso_char_code), ''),@p2,@p1) cr
		)
`
	moneyPortfolioFxJoinSQL = `
		JOIN fx ON fx.currency_code = c.currency_code
		WHERE c.settle_code = c.settle_max	
	`
	selectMoneyPortfolioRowsAllSQL       = moneyPortfolioLatestSettleCTEAllSQL + moneyPortfolioFxRatesCTESQL + moneyPortfolioSelectColumnsSQL + moneyPortfolioFxJoinSQL
	selectMoneyPortfolioRowsByClientsSQL = moneyPortfolioLatestSettleCTEByClientsSQL + moneyPortfolioFxRatesCTESQL + moneyPortfolioSelectColumnsSQL + moneyPortfolioFxJoinSQL
)

func (r *Repository) ListMoneyPortfolio(ctx context.Context, date time.Time, targetCcy string, clientCodes []string) (result []quik.Position, err error) {
	const opName = "ListMoneyPortfolio"
	start := time.Now()
	defer func() { err = r.observeSelectExit(opName, date, start, err) }()
	pos, err := selectPortfolioRows(
		r,
		ctx,
		selectMoneyPortfolioRowsAllSQL,
		selectMoneyPortfolioRowsByClientsSQL,
		scanMoneyPortfolioRow,
		portfolioQuery{date: date, targetCcy: targetCcy, clientCodes: clientCodes},
	)
	if err != nil {
		return nil, err
	}
	return mapRows(pos, moneyPortfolioRow.toQuikPosition), nil
}
