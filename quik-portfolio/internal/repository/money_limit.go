package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

const selectMoneyLimitsByClients = `
SELECT
    li.load_date,
    li.source_date,
    li.client_code,
    li.ccy,
    li.position_code,
    li.settle_code,
    li.firm_code,
    li.firm_name,
    li.balance
FROM quik.money_limits li
join @codes c on c.client_code = li.client_code
WHERE 1 = 1
    and li.load_date = cast(@p1 as date) 
ORDER BY li.load_date, 
	li.client_code, 
	li.ccy,
	li.position_code,
    li.settle_code,
    li.firm_code 
OFFSET @p2 ROWS 
FETCH NEXT @p3 ROWS ONLY`

const selectMoneyLimitsAllClients = `SELECT
    li.load_date,
    li.source_date,
    li.client_code,
    li.ccy,
    li.position_code,
    li.settle_code,
    li.firm_code,
    li.firm_name,
    li.balance
FROM quik.money_limits li
WHERE 1 = 1
    and li.load_date = cast(@p1 as date) 
ORDER BY li.load_date, 
	li.client_code, 
	li.ccy,
	li.position_code,
    li.settle_code,
    li.firm_code 
OFFSET @p2 ROWS 
FETCH NEXT @p3 ROWS ONLY`

const countMoneyLimitsByClients = `SELECT COUNT(*)
FROM quik.money_limits li
join @codes c on c.client_code = li.client_code
WHERE li.load_date =  cast(@p1 as date) 
`

const countMoneyLimitsAllClients = `SELECT COUNT(*)
FROM quik.money_limits li
WHERE li.load_date = cast(@p1 as date) `

func (r *Repository) scanMoneyLimit(row *sql.Rows) (quik.MoneyLimit, error) {
	out := quik.MoneyLimit{}
	err := row.Scan(
		&out.LoadDate,
		&out.SourceDate,
		&out.ClientCode,
		&out.Currency,
		&out.PositionCode,
		&out.SettleCode,
		&out.FirmCode,
		&out.FirmName,
		&out.Balance,
	)
	if err != nil {
		return quik.MoneyLimit{}, err
	}
	return out, nil
}

func (r *Repository) SelectMoneyLimitsWithFilters(ctx context.Context, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.MoneyLimit, totalCount *uint64, err error) {
	start := time.Now()
	defer func() { r.metrics.ObserveRepository("SelectMoneyLimitsWithFilters", time.Since(start), err) }()

	return selectLimitsWithFilters(r, ctx, "SelectMoneyLimitsWithFilters", date, limit, offset, clientCodes, includeTotalCount, limitFilterSQL{
		countByClients:  countMoneyLimitsByClients,
		countAll:        countMoneyLimitsAllClients,
		selectByClients: selectMoneyLimitsByClients,
		selectAll:       selectMoneyLimitsAllClients,
	}, r.scanMoneyLimit)
}
