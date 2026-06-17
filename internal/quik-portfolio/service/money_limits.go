package service

import (
	"context"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (s *Service) GetMoneyLimitsWithFilters(ctx context.Context, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.MoneyLimit, totalCount *uint64, err error) {
	dedublicated, err := validateLimitsContract(date, clientCodes)
	if err != nil {
		return nil, nil, err
	}
	return s.repo.ListMoneyLimits(ctx, date, limit, offset, dedublicated, includeTotalCount)
}
