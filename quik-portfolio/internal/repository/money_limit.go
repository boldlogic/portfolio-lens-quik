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

type moneyLimitRow struct {
	LoadDate     time.Time
	SourceDate   time.Time
	ClientCode   string
	CurrencyCode string
	PositionCode string
	SettleCode   string
	FirmCode     string
	FirmName     sql.NullString
	Balance      *decimal.Decimal
}

func scanMoneyLimitRow(row *sql.Rows) (moneyLimitRow, error) {
	out := moneyLimitRow{}
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
		return moneyLimitRow{}, fmt.Errorf("%w: %w", ErrScan, err)
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
	selectMoneyLimitsByClients  = moneyLimitSelectColumnsSQL + "join @codes c on c.client_code = li.client_code" + moneyLimitPageClauseSQL
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

func (r *Repository) SelectMoneyLimitsWithFilters(ctx context.Context, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.MoneyLimit, totalCount *uint64, err error) {
	start := time.Now()
	defer func() { r.metrics.ObserveRepository("SelectMoneyLimitsWithFilters", time.Since(start), err) }()

	limits, totalCount, err := selectLimitsWithFilters(r, ctx, "SelectMoneyLimitsWithFilters", date, limit, offset, clientCodes, includeTotalCount, limitFilterSQL{
		countByClients:  countMoneyLimitsByClients,
		countAll:        countMoneyLimitsAllClients,
		selectByClients: selectMoneyLimitsByClients,
		selectAll:       selectMoneyLimitsAllClients,
	}, scanMoneyLimitRow)
	if err != nil {
		return nil, nil, err
	}
	for _, raw := range limits {
		limit, validationErr := raw.toQuik()
		if validationErr != nil {
			r.Logger.Warn("нужно завести новый код", zap.Error(validationErr),
				zap.Time("load_date", raw.LoadDate),
				zap.String("client_code", raw.ClientCode),
				zap.String("currency_code", raw.CurrencyCode),
				zap.String("position_code", raw.PositionCode),
				zap.String("settle_code", raw.SettleCode),
				zap.String("firm_code", raw.FirmCode),
			)
		}
		result = append(result, limit)
		if raw.LoadDate.Before(date) {
			r.Logger.Warn("устаревший лимит по деньгам", zap.Time("запрошена дата", date),
				zap.Time("load_date", raw.LoadDate),
				zap.String("client_code", raw.ClientCode),
				zap.String("currency_code", raw.CurrencyCode),
				zap.String("position_code", raw.PositionCode),
				zap.String("settle_code", raw.SettleCode),
				zap.String("firm_code", raw.FirmCode),
			)

		}
	}

	return result, totalCount, nil
}

func (m moneyLimitRow) toQuik() (quik.MoneyLimit, error) {
	ml := quik.MoneyLimit{
		LoadDate:     m.LoadDate,
		SourceDate:   m.SourceDate,
		ClientCode:   m.ClientCode,
		Currency:     m.CurrencyCode,
		PositionCode: m.PositionCode,
		FirmCode:     m.FirmCode,
	}

	if m.Balance != nil {
		ml.Balance = *m.Balance
	}
	if m.FirmName.Valid {
		ml.FirmName = m.FirmName.String
	}
	ml.SettleCode = quik.SettleCode(m.SettleCode)

	err := ml.SettleCode.Validate()
	if err != nil {
		return ml, err
	}

	return ml, nil
}
