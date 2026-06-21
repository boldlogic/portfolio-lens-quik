package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
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
		select p.load_date,
			p.source_date,
			p.quote_date,
			fx_rate_date=fx.rate_date,
			p.client_code,
			p.firm_code,
			p.firm_name,
			p.sec_code,
			p.sec_name,
			p.balance,
			p.unit_price,
			p.accrued_interest,
			p.market_value_in_instr_currency,
			p.instrument_currency_code,
			market_value_in_target_currency=cast(p.market_value_in_instr_currency/fx.rate as decimal(18,2))
		from positions p
`
	securityPortfolioLatestSettleCTEBaseSQL = `
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

	securityPortfolioCteWhereAllSQL = `
			WHERE li.load_date = cast(@p1 as date)
		)`
	securityPortfolioCteWhereByClientsSQL = `
			WHERE li.load_date = cast(@p1 as date)
				AND EXISTS (SELECT 1 FROM @codes c WHERE c.client_code = li.client_code)
		)`

	securityPortfolioPositionsAndFxCTESQL = `
		, positions AS (
			SELECT
				c.load_date,
				c.source_date,
				q.quote_date,
				c.client_code,
				c.firm_code,
				c.firm_name,
				c.sec_code,
				sec_name = COALESCE(c.sec_name, q.short_name),
				c.balance,
				unit_price = q.price,
				accrued_interest = q.ai,
				market_value = q.market_value,
				market_value_in_instr_currency = CAST(q.market_value * c.balance AS DECIMAL(18, 2)),
				instrument_currency_code = norm_ccy.currency
			FROM cte c
		outer apply market.fnSecurityQuoteByAcquisitionCurrency(c.sec_code,c.acquisition_currency_code) q
		outer apply (select ccy=coalesce(q.currency,c.acquisition_currency_code)) cur
		outer apply (select currency=case when isnull(cur.ccy,'')  IN ('SUR','RUR') THEN 'RUB' ELSE cur.ccy END) norm_ccy
			WHERE c.settle_code = c.settle_max
		)
		, fx_keys AS (
			SELECT DISTINCT
				p.instrument_currency_code
			FROM positions p
		)
		, fx AS (
			SELECT
				k.instrument_currency_code,
				cr.rate,
				cr.rate_date
			FROM fx_keys k
		LEFT JOIN ref.external_codes ec_mv
			ON ec_mv.ext_system_id = (select
				ext_system_id
			from
				ref.external_systems
			where
				ext_system = 'QUIK')
			AND ec_mv.ext_code_type_id = 1
				AND ec_mv.ext_code = k.instrument_currency_code
		LEFT JOIN ref.currencies c_mv_ec  ON c_mv_ec.iso_code  = ec_mv.internal_id
			LEFT JOIN ref.currencies c_mv_iso ON c_mv_iso.iso_char_code = k.instrument_currency_code
		cross apply market.fnFxRateCross(ISNULL(COALESCE(c_mv_ec.iso_char_code,c_mv_iso.iso_char_code), ''),@p2,@p1) cr
		)
`
	securityPortfolioFxJoinSQL = `
		JOIN fx ON fx.instrument_currency_code = p.instrument_currency_code`
)

func buildSecurityPortfolioSQL(sourceTable string, cteWhereSQL string) string {
	return securityPortfolioLatestSettleCTEBaseSQL +
		sourceTable +
		cteWhereSQL +
		securityPortfolioPositionsAndFxCTESQL +
		securityPortfolioSelectColumnsSQL +
		securityPortfolioFxJoinSQL
}

var (
	selectSecurityPortfolioExchangeAllSQL       = buildSecurityPortfolioSQL(securityLimitExchangeTableSQL, securityPortfolioCteWhereAllSQL)
	selectSecurityPortfolioExchangeByClientsSQL = buildSecurityPortfolioSQL(securityLimitExchangeTableSQL, securityPortfolioCteWhereByClientsSQL)
	selectSecurityPortfolioOtcAllSQL            = buildSecurityPortfolioSQL(securityLimitOtcTableSQL, securityPortfolioCteWhereAllSQL)
	selectSecurityPortfolioOtcByClientsSQL      = buildSecurityPortfolioSQL(securityLimitOtcTableSQL, securityPortfolioCteWhereByClientsSQL)
)

func (r *Repository) listSecurityPortfolio(
	ctx context.Context,
	date time.Time,
	targetCcy string,
	clientCodes []string,
	limitType quik.LimitType,
	selectAllSQL string,
	selectByClientsSQL string,
) (result []quik.Position, err error) {
	pos, err := selectPortfolioRows(
		r,
		ctx,
		selectAllSQL,
		selectByClientsSQL,
		scanSecurityPortfolioRow,
		portfolioQuery{date: date, targetCcy: targetCcy, clientCodes: clientCodes},
	)
	if err != nil {
		return nil, err
	}
	return mapRows(pos, func(row securityPortfolioRow) quik.Position {
		return row.toQuikPositionWithType(limitType)
	}), nil
}

func (r *Repository) ListSecurityPortfolio(ctx context.Context, date time.Time, targetCcy string, clientCodes []string) (result []quik.Position, err error) {
	const opName = "ListSecurityPortfolio"
	start := time.Now()
	defer func() { err = r.observeSelectExit(opName, date, start, err) }()
	return r.listSecurityPortfolio(
		ctx,
		date,
		targetCcy,
		clientCodes,
		quik.LimitTypeSecurities,
		selectSecurityPortfolioExchangeAllSQL,
		selectSecurityPortfolioExchangeByClientsSQL,
	)
}

func (r *Repository) ListSecurityPortfolioOtc(ctx context.Context, date time.Time, targetCcy string, clientCodes []string) (result []quik.Position, err error) {
	const opName = "ListSecurityPortfolioOtc"
	start := time.Now()
	defer func() { err = r.observeSelectExit(opName, date, start, err) }()
	return r.listSecurityPortfolio(
		ctx,
		date,
		targetCcy,
		clientCodes,
		quik.LimitTypeSecuritiesOtc,
		selectSecurityPortfolioOtcAllSQL,
		selectSecurityPortfolioOtcByClientsSQL,
	)
}
