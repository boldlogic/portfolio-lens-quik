package service

import (
	"context"
	"errors"

	"github.com/boldlogic/portfolio-lens-quik/internal/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"go.uber.org/zap"
)

func (s *Service) UpsertLimits(ctx context.Context, limits []models.LimitLine) error {

	var errs []error

	var out []quik.Limit
	for _, row := range limits {
		limit, err := quik.NewLimit(
			row.Type,
			row.ClientCode,
			row.Ticker,
			row.PositionCode,
			row.SettleCode,
			row.TradeAccount,
			row.FirmCode,
			row.Balance,
			row.AcquisitionCurrencyCode,
			row.ISIN,
		)
		if err != nil {
			errs = append(errs, err)
			s.logger.Error(err.Error())
			continue
		}
		out = append(out, limit)
		s.logger.Debug("", zap.Any("", limit))
	}
	if len(errs) != 0 {
		return errors.Join(errs...)
	}

	return s.repo.HandleRequest(ctx, out)

}
