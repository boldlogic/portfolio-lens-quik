package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/internal/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func scanSecurityLimitRow(row *sql.Rows) (models.SecurityLimitRow, error) {
	out := models.SecurityLimitRow{}
	err := row.Scan(
		&out.LoadDate,
		&out.SourceDate,
		&out.ClientCode,
		&out.SecCode,
		&out.SettleCode,
		&out.TradeAccount,
		&out.FirmCode,
		&out.FirmName,
		&out.Balance,
		&out.AcquisitionCurrencyCode,
		&out.ISIN,
		&out.ShortName,
	)
	if err != nil {
		return models.SecurityLimitRow{}, fmt.Errorf("%w: %w", ErrScan, err)
	}

	return out, nil
}

const (
	securityLimitSelectColumnsSQL = `
        SELECT
            li.load_date,
            li.source_date,
            li.client_code,
            li.sec_code,
            li.settle_code,
            li.trade_account,
            li.firm_code,
            li.firm_name,
            li.balance,
            li.acquisition_currency_code,
            li.isin,
            li.sec_name
        
    `
	securityLimitExchangeTableSQL = `FROM quik.security_limits li`
	securityLimitOtcTableSQL      = `FROM quik.security_limits_otc li`
	securityLimitPageClauseSQL    = `
        WHERE 
            li.load_date = cast(@p1 as date)
        ORDER BY li.load_date, 
            li.client_code, 
            li.sec_code, 
            li.trade_account, 
            li.settle_code,
            li.firm_code
        OFFSET @p2 ROWS 
        FETCH NEXT @p3 ROWS ONLY
`
	selectSecurityLimitsByClients     = securityLimitSelectColumnsSQL + securityLimitExchangeTableSQL + " join @codes c on c.client_code = li.client_code" + securityLimitPageClauseSQL
	selectSecurityLimitsOtcByClients  = securityLimitSelectColumnsSQL + securityLimitOtcTableSQL + " join @codes c on c.client_code = li.client_code" + securityLimitPageClauseSQL
	selectSecurityLimitsAllClients    = securityLimitSelectColumnsSQL + securityLimitExchangeTableSQL + securityLimitPageClauseSQL
	selectSecurityLimitsOtcAllClients = securityLimitSelectColumnsSQL + securityLimitOtcTableSQL + securityLimitPageClauseSQL

	countSecurityLimitsByClients = `
        SELECT COUNT(*)
        FROM quik.security_limits li
        join @codes c on c.client_code = li.client_code
        WHERE li.load_date = cast(@p1 as date)`

	countSecurityLimitsAllClients = `
        SELECT COUNT(*)
        FROM quik.security_limits li
        WHERE li.load_date = cast(@p1 as date)`

	countSecurityLimitsOtcByClients = `SELECT COUNT(*)
FROM quik.security_limits_otc li
join @codes c on c.client_code = li.client_code
WHERE li.load_date = cast(@p1 as date)`

	countSecurityLimitsOtcAllClients = `SELECT COUNT(*)
FROM quik.security_limits_otc li
WHERE li.load_date = cast(@p1 as date)`
)

func (r *Repository) ListSecurityLimits(ctx context.Context, limitType quik.LimitType, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.SecurityLimit, totalCount *uint64, err error) {
	const opName = "ListSecurityLimits"
	start := time.Now()
	defer func() { err = r.observeSelectExit(opName, date, start, err) }()

	limits, totalCount, err := selectLimitRows(r, ctx, limitListQuery{
		date:              date,
		limit:             limit,
		offset:            offset,
		clientCodes:       clientCodes,
		includeTotalCount: includeTotalCount,
	}, limitType, scanSecurityLimitRow)
	if err != nil {
		return nil, nil, err
	}

	return mapRows(limits, models.SecurityLimitRow.ToQuik), totalCount, nil
}
