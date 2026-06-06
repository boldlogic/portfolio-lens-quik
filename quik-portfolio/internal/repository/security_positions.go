package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

type securityPortfolioRow struct {
	LoadDate                    time.Time
	SourceDate                  time.Time
	QuoteDate                   sql.NullTime
	FxRateDate                  sql.NullTime
	ClientCode                  string
	FirmCode                    string
	FirmName                    sql.NullString
	SecCode                     string
	SecName                     sql.NullString
	Balance                     decimal.Decimal
	UnitPrice                   *decimal.Decimal
	AccruedInterest             *decimal.Decimal
	MarketValueInInstrCurrency  *decimal.Decimal
	InstrumentCurrencyCode      string
	MarketValueInTargetCurrency *decimal.Decimal
}

func scanSecurityPortfolioRow(row *sql.Rows) (securityPortfolioRow, error) {
	var out securityPortfolioRow
	err := row.Scan(
		&out.LoadDate,
		&out.SourceDate,
		&out.QuoteDate,
		&out.FxRateDate,
		&out.ClientCode,
		&out.FirmCode,
		&out.FirmName,
		&out.SecCode,
		&out.SecName,
		&out.Balance,
		&out.UnitPrice,
		&out.AccruedInterest,
		&out.MarketValueInInstrCurrency,
		&out.InstrumentCurrencyCode,
		&out.MarketValueInTargetCurrency,
	)
	if err != nil {
		return securityPortfolioRow{}, fmt.Errorf("%w: %w", ErrScan, err)
	}

	return out, nil
}

const (
	securityPortfolioSelectColumnsSQL = `
		select c.load_date,
			c.source_date,
			q.quote_date,
			fx_rate_date=cr.rate_date,
			c.client_code,
			c.firm_code,
			c.firm_name,
			c.sec_code,	
			sec_name=coalesce(c.sec_name,q.short_name),
			c.balance,
			unit_price=q.price,
			accrued_interest=q.ai,
			market_value_in_instr_currency=cast(q.market_value*c.balance as decimal(18,2)),
			instrument_currency_code=norm_ccy.currency,
			market_value_in_target_currency=cast (q.market_value*c.balance*cr.rate as decimal(18,2))
		from cte c
`
	securityPortfolioLatestSettleCTESQL = `
		WITH cte AS (
			SELECT
				li.load_date,
				li.source_date,
				li.client_code,
				li.sec_code,
				li.sec_name,
				li.settle_code,
				li.firm_code,
				li.firm_name,
				li.balance,
				li.acquisition_currency_code,
				li.isin,
				settle_max = MAX(li.settle_code) 
					OVER (
					PARTITION BY 
						li.load_date, 
						li.client_code, 
						li.sec_code, 
						li.trade_account, 
						li.firm_code
					)
`
	securityPortfolioSourceTableSQL = "FROM quik.security_limits li "

	securityPortfolioDateFilterSQL = `WHERE li.load_date = cast(@p1 as date))`

	securityPortfolioQuoteAndFxJoinSQL = `
		outer apply dbo.fnGetQuoteByAcquisitionCurrency(c.sec_code,c.acquisition_currency_code) q
		outer apply (select ccy=coalesce(q.currency,c.acquisition_currency_code)) cur
		outer apply (select currency=case when isnull(cur.ccy,'')  IN ('SUR','RUR') THEN 'RUB' ELSE cur.ccy END) norm_ccy
		LEFT JOIN dbo.external_codes ec_mv
			ON ec_mv.ext_system_id = (select
				ext_system_id
			from
				dbo.external_systems
			where
				ext_system = 'QUIK')
			AND ec_mv.ext_code_type_id = 1
			AND ec_mv.ext_code = norm_ccy.currency
		LEFT JOIN currencies c_mv_ec  ON c_mv_ec.iso_code  = ec_mv.internal_id
		LEFT JOIN currencies c_mv_iso ON c_mv_iso.iso_char_code = norm_ccy.currency
		cross apply dbo.fnFxRateCross(ISNULL(COALESCE(c_mv_ec.iso_char_code,c_mv_iso.iso_char_code), ''),@p2,@p1) cr
		where c.settle_code=c.settle_max`

	selectSecurityPortfolioRowsSQL = securityPortfolioLatestSettleCTESQL + securityPortfolioSourceTableSQL + securityPortfolioDateFilterSQL + securityPortfolioSelectColumnsSQL + securityPortfolioQuoteAndFxJoinSQL
)

func (r *Repository) ListSecurityPortfolio(ctx context.Context, date time.Time, targetCcy string) (result []quik.Position, err error) {
	pos, err := selectPortfolioRows(
		r,
		ctx,
		"ListSecurityPortfolio",
		selectSecurityPortfolioRowsSQL,
		scanSecurityPortfolioRow,
		portfolioQuery{date: date, targetCcy: targetCcy},
	)
	if err != nil {
		return nil, err
	}
	return mapRows(pos, securityPortfolioRow.toQuikPosition), nil
}
