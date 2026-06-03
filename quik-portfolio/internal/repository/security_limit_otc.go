package repository

import (
	"context"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

const (
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
    li.sec_name
FROM quik.security_limits_otc li
join @codes c on c.client_code = li.client_code
WHERE li.load_date = cast(@p1 as date)
ORDER BY li.load_date, 
li.client_code, 
li.ticker, 
li.trade_account, 
li.settle_code,
li.firm_code
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
    li.sec_name
FROM quik.security_limits_otc li
WHERE li.load_date = cast(@p1 as date)
ORDER BY li.load_date, 
li.client_code, 
li.ticker, 
li.trade_account, 
li.settle_code,
li.firm_code
OFFSET @p2 ROWS FETCH NEXT @p3 ROWS ONLY`

	countSecurityLimitsOtcByClients = `SELECT COUNT(*)
FROM quik.security_limits_otc li
join @codes c on c.client_code = li.client_code
WHERE li.load_date = cast(@p1 as date)`

	countSecurityLimitsOtcAllClients = `SELECT COUNT(*)
FROM quik.security_limits_otc li
WHERE li.load_date = cast(@p1 as date)`
)

func (r *Repository) SelectSecurityLimitsOtcWithFilters(ctx context.Context, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.SecurityLimit, totalCount *uint64, err error) {
	start := time.Now()
	defer func() { r.metrics.ObserveRepository("SelectSecurityLimitsOtcWithFilters", time.Since(start), err) }()

	return selectLimitsWithFilters(r, ctx, "SelectSecurityLimitsOtcWithFilters", date, limit, offset, clientCodes, includeTotalCount, limitFilterSQL{
		countByClients:  countSecurityLimitsOtcByClients,
		countAll:        countSecurityLimitsOtcAllClients,
		selectByClients: selectSecurityLimitsOtcByClients,
		selectAll:       selectSecurityLimitsOtcAllClients,
	}, r.scanSecurityLimit)
}
