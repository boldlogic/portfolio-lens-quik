package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"go.uber.org/zap"
)

const (
	selectSecurityLimitsOtcByDate = `
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
			FROM quik.security_limits_otc li
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

	selectSecurityLimitsOtcByClients = `SELECT
    li.load_date,
    li.source_date,
    li.client_code,
    li.ticker,
    li.settle_code,
    li.trade_account,
    li.firm_code,
    li.firm_name,
    li.balance,
    li.acquisition_ccy,
    li.isin,
    short_disp.short_name
FROM quik.security_limits_otc li
join @codes c on c.client_code = li.client_code
` + securityLimitShortNameApply + `
WHERE li.load_date = cast(@p1 as date)
ORDER BY li.load_date, li.client_code, li.ticker, li.trade_account, li.firm_code
OFFSET @p2 ROWS FETCH NEXT @p3 ROWS ONLY`

	selectSecurityLimitsOtcAllClients = `SELECT
    li.load_date,
    li.source_date,
    li.client_code,
    li.ticker,
    li.settle_code,
    li.trade_account,
    li.firm_code,
    li.firm_name,
    li.balance,
    li.acquisition_ccy,
    li.isin,
    short_disp.short_name
FROM quik.security_limits_otc li
` + securityLimitShortNameApply + `
WHERE li.load_date = cast(@p1 as date)
ORDER BY li.load_date, li.client_code, li.ticker, li.trade_account, li.firm_code
OFFSET @p2 ROWS FETCH NEXT @p3 ROWS ONLY`

	countSecurityLimitsOtcByClients = `SELECT COUNT(*)
FROM quik.security_limits_otc li
join @codes c on c.client_code = li.client_code
WHERE li.load_date = cast(@p1 as date)`

	countSecurityLimitsOtcAllClients = `SELECT COUNT(*)
FROM quik.security_limits_otc li
WHERE li.load_date = cast(@p1 as date)`
)

func (r *Repository) SelectSecurityLimitsOtc(ctx context.Context, date time.Time) (result []quik.SecurityLimit, err error) {
	defer func() { err = r.finalizeSelectErr("SelectSecurityLimitsOtc", date, err) }()

	rows, err := r.Db.QueryContext(ctx, selectSecurityLimitsOtcByDate, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		row, err := r.scanSecurityLimit(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	r.Logger.Debug("результаты получения позиций OTC по бумагам", zap.Time("load_date", date), zap.Int("count", len(result)))
	if len(result) == 0 {
		r.Logger.Warn("позиции OTC по бумагам не найдены", zap.Time("load_date", date))
	}
	return result, nil
}

func (r *Repository) SelectSecurityLimitsOtcWithFilters(ctx context.Context, date time.Time, limit, offset int, clientCodes []string) (result []quik.SecurityLimit, totalCount int, err error) {
	defer func() { err = r.finalizeSelectErr("SelectSecurityLimitsOtcWithFilters", date, err) }()

	clients, hasClients := r.makeClientCodeList(clientCodes)

	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	var rows *sql.Rows
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()

	if hasClients {
		err = tx.QueryRowContext(ctx, countSecurityLimitsOtcByClients, date, sql.Named("codes", clients)).Scan(&totalCount)
		if err != nil {
			return nil, 0, err
		}
		rows, err = tx.QueryContext(ctx, selectSecurityLimitsOtcByClients, date, offset, limit, sql.Named("codes", clients))
		if err != nil {
			return nil, 0, err
		}
	} else {
		err = tx.QueryRowContext(ctx, countSecurityLimitsOtcAllClients, date).Scan(&totalCount)
		if err != nil {
			return nil, 0, err
		}
		rows, err = tx.QueryContext(ctx, selectSecurityLimitsOtcAllClients, date, offset, limit)
		if err != nil {
			return nil, 0, err
		}
	}

	for rows.Next() {
		var row quik.SecurityLimit
		row, err = r.scanSecurityLimit(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, row)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	return result, totalCount, nil
}
