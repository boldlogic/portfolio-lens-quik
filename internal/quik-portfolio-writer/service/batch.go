package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/boldlogic/portfolio-lens-quik/internal/models"
	errmodel "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
)

func (s *Service) UpsertLimits(ctx context.Context, limits []models.LimitLine) error {

	var errs []error

	out := make([]quik.Limit, 0, len(limits))
	byHash := make(map[[32]byte][]uint, len(limits))
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
			continue
		}
		h := limit.KeyHash()
		byHash[h] = append(byHash[h], row.Line)
		out = append(out, limit)
	}
	if len(errs) != 0 {
		err := errors.Join(errs...)
		return fmt.Errorf("%w: %w", errmodel.ErrBusinessValidation, err)
	}

	if groups := duplicateLineGroupsFromHash(byHash); len(groups) > 0 {
		return errDuplicateLimits(groups)
	}

	return s.repo.HandleRequest(ctx, out)

}
