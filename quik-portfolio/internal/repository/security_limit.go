package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (s securityLimitRow) toQuik() (quik.SecurityLimit, error) {
	sl := quik.SecurityLimit{
		LoadDate:       s.LoadDate,
		SourceDate:     s.SourceDate,
		ClientCode:     s.ClientCode,
		Ticker:         s.SecCode,
		TradeAccount:   s.TradeAccount,
		FirmCode:       s.FirmCode,
		AcquisitionCcy: s.AcquisitionCurrencyCode,
	}
	if s.Balance != nil {
		sl.Balance = *s.Balance
	}
	if s.FirmName.Valid {
		sl.FirmName = s.FirmName.String
	}
	if s.ISIN.Valid {
		sl.ISIN = s.ISIN.String
	}
	sl.SettleCode = quik.SettleCode(s.SettleCode)

	err := sl.SettleCode.Validate()
	if err != nil {
		return sl, err
	}
	return sl, nil

}

type securityLimitRow struct {
	LoadDate                  time.Time
	SourceDate                time.Time
	ClientCode                string
	SecCode                   string
	TradeAccount              string
	SettleCode                string
	FirmCode                  string
	FirmName                  sql.NullString
	Balance                   *decimal.Decimal
	AcquisitionCurrencyCode   string
	ISIN                      sql.NullString
	ShortName                 sql.NullString
}

func scanSecurityLimitRow(row *sql.Rows) (securityLimitRow, error) {
	out := securityLimitRow{}
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
		return securityLimitRow{}, fmt.Errorf("%w: %w", ErrScan, err)
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
	securityLimitOtcTableSQL        = `FROM quik.security_limits_otc li`
	securityLimitPageClauseSQL      = `
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
	selectSecurityLimitsByClients     = securityLimitSelectColumnsSQL + securityLimitExchangeTableSQL + "join @codes c on c.client_code = li.client_code" + securityLimitPageClauseSQL
	selectSecurityLimitsOtcByClients  = securityLimitSelectColumnsSQL + securityLimitOtcTableSQL + "join @codes c on c.client_code = li.client_code" + securityLimitPageClauseSQL
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
)

func (r *Repository) SelectSecurityLimitsWithFilters(ctx context.Context, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.SecurityLimit, totalCount *uint64, err error) {
	start := time.Now()
	defer func() { r.metrics.ObserveRepository("SelectSecurityLimitsWithFilters", time.Since(start), err) }()

	limits, totalCount, err := selectLimitsWithFilters(r, ctx, "SelectSecurityLimitsWithFilters", date, limit, offset, clientCodes, includeTotalCount, limitFilterSQL{
		countByClients:  countSecurityLimitsByClients,
		countAll:        countSecurityLimitsAllClients,
		selectByClients: selectSecurityLimitsByClients,
		selectAll:       selectSecurityLimitsAllClients,
	}, scanSecurityLimitRow)
	if err != nil {
		return nil, nil, err
	}
	for _, raw := range limits {
		limit, validationErr := raw.toQuik()
		if validationErr != nil {
			r.Logger.Warn("нужно завести новый код", zap.Error(validationErr),
				zap.Time("load_date", raw.LoadDate),
				zap.String("client_code", raw.ClientCode),
				zap.String("sec_code", raw.SecCode),
				zap.String("trade_account", raw.TradeAccount),
				zap.String("settle_code", raw.SettleCode),
				zap.String("firm_code", raw.FirmCode),
			)
		}
		result = append(result, limit)
		if raw.LoadDate.Before(date) {
			r.Logger.Warn("устаревший лимит по деньгам", zap.Time("запрошена дата", date),
				zap.Time("load_date", raw.LoadDate),
				zap.String("client_code", raw.ClientCode),
				zap.String("sec_code", raw.SecCode),
				zap.String("trade_account", raw.TradeAccount),
				zap.String("settle_code", raw.SettleCode),
				zap.String("firm_code", raw.FirmCode),
			)

		}
	}
	return result, totalCount, nil
}
