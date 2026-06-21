package service

import (
	"context"
	"fmt"

	intmodel "github.com/boldlogic/portfolio-lens-quik/internal/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
)

func (s *Service) UpsertLimit(ctx context.Context, limit intmodel.LimitInput) error {

	lim, err := quik.NewLimit(limit.Type,
		limit.ClientCode,
		limit.Ticker,
		limit.PositionCode,
		limit.SettleCode,
		limit.TradeAccount,
		limit.FirmCode,
		limit.Balance,
		limit.AcquisitionCurrencyCode,
		limit.ISIN)
	if err != nil {
		return fmt.Errorf("%w: %w", models.ErrBusinessValidation, err)
	}
	out := []quik.Limit{lim}
	return s.repo.HandleRequest(ctx, out)
}
