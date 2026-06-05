package repository

import (
	"context"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"go.uber.org/zap"
)

const (
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

	limits, totalCount, err := selectLimitsWithFilters(r, ctx, "SelectSecurityLimitsOtcWithFilters", date, limit, offset, clientCodes, includeTotalCount, limitFilterSQL{
		countByClients:  countSecurityLimitsOtcByClients,
		countAll:        countSecurityLimitsOtcAllClients,
		selectByClients: selectSecurityLimitsOtcByClients,
		selectAll:       selectSecurityLimitsOtcAllClients,
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
