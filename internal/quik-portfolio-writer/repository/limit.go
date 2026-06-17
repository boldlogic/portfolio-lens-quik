package repository

import (
	"context"
	"fmt"

	"github.com/boldlogic/portfolio-lens-quik/internal/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (r *Repository) InsertLimit(ctx context.Context, l quik.Limit) (quik.Limit, error) {
	var err error
	switch l.Type {
	case quik.LimitTypeMoney:
		var m models.MoneyLimitRow

		row := r.Db.QueryRowContext(ctx, insertMoneyLimit,
			l.ClientCode,
			l.CurrencyCode,
			l.PositionCode,
			l.FirmCode,
			l.Balance,
		)
		m, err = models.ScanMoneyLimitRow(row)
		if err != nil {
			return quik.Limit{}, err
		}
		return m.ToLimit(), nil
	case quik.LimitTypeSecurities:
		var s models.SecurityLimitRow
		row := r.Db.QueryRowContext(ctx, insertSecurityLimit,
			s.ClientCode,
			s.SecCode,
			s.TradeAccount,
			string(s.SettleCode),
			s.FirmCode,
			s.Balance,
			s.AcquisitionCurrencyCode,
			s.ISIN)
		s, err = models.ScanSecurityLimitRow(row)
		if err != nil {
			return quik.Limit{}, err
		}
		return s.ToLimit(), nil
	case quik.LimitTypeSecuritiesOtc:
		var s models.SecurityLimitRow
		row := r.Db.QueryRowContext(ctx, insertSecurityLimitOtc,
			s.ClientCode,
			s.SecCode,
			s.TradeAccount,
			string(s.SettleCode),
			s.FirmCode,
			s.Balance,
			s.AcquisitionCurrencyCode,
			s.ISIN)
		s, err = models.ScanSecurityLimitRow(row)
		if err != nil {
			return quik.Limit{}, err
		}
		return s.ToLimit(), nil
	default:
		return quik.Limit{}, fmt.Errorf("неподдерживаемый тип лимита")
	}

}
