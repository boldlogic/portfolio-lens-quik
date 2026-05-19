package service

import (
	"context"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (s *Service) GetMoneyLimits(ctx context.Context, date time.Time) ([]quik.MoneyLimit, error) {
	return s.repo.SelectMoneyLimits(ctx, date)
}

func (s *Service) GetMoneyLimitsWithFilters(ctx context.Context, date time.Time, limit, offset int, clientCodes []string) ([]quik.MoneyLimit, int, error) {
	return s.repo.SelectMoneyLimitsWithFilters(ctx, date, limit, offset, clientCodes)
}
