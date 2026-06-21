package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/internal/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
)

func scanMoneyLimitRow(row *sql.Rows) (models.MoneyLimitRow, error) {
	out := models.MoneyLimitRow{}
	err := row.Scan(
		&out.LoadDate,
		&out.SourceDate,
		&out.ClientCode,
		&out.CurrencyCode,
		&out.PositionCode,
		&out.SettleCode,
		&out.FirmCode,
		&out.FirmName,
		&out.Balance,
	)
	if err != nil {
		return models.MoneyLimitRow{}, fmt.Errorf("%w: %w", ErrScan, err)
	}
	return out, nil
}

const (
	moneyLimitSelectColumnsSQL = `
        SELECT
            li.load_date,
            li.source_date,
            li.client_code,
            li.currency_code,
            li.position_code,
            li.settle_code,
            li.firm_code,
            li.firm_name,
            li.balance
        FROM quik.money_limits li
    `
	moneyLimitPageClauseSQL = `
        WHERE 1 = 1
            and li.load_date = cast(@p1 as date) 
        ORDER BY li.load_date, 
            li.client_code, 
            li.currency_code,
            li.position_code,
            li.settle_code,
            li.firm_code 
        OFFSET @p2 ROWS 
        FETCH NEXT @p3 ROWS ONLY
`
	selectMoneyLimitsAllClients = moneyLimitSelectColumnsSQL + moneyLimitPageClauseSQL
	selectMoneyLimitsByClients  = moneyLimitSelectColumnsSQL + " join @codes c on c.client_code = li.client_code" + moneyLimitPageClauseSQL
	countMoneyLimitsByClients   = `
        SELECT COUNT(*)
        FROM quik.money_limits li
        join @codes c on c.client_code = li.client_code
        WHERE li.load_date =  cast(@p1 as date) 
        `
	countMoneyLimitsAllClients = `
        SELECT COUNT(*)
        FROM quik.money_limits li
        WHERE li.load_date = cast(@p1 as date) `
)

func (r *Repository) ListMoneyLimits(ctx context.Context, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.MoneyLimit, totalCount *uint64, err error) {
	const opName = "ListMoneyLimits"
	start := time.Now()
	defer func() { err = r.observeSelectExit(opName, date, start, err) }()

	limits, totalCount, err := selectLimitRows(r, ctx, limitListQuery{
		date:              date,
		limit:             limit,
		offset:            offset,
		clientCodes:       clientCodes,
		includeTotalCount: includeTotalCount,
	}, quik.LimitTypeMoney, scanMoneyLimitRow)
	if err != nil {
		return nil, nil, err
	}

	return mapRows(limits, models.MoneyLimitRow.ToQuik), totalCount, nil
}
