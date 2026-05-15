package service

import (
	"context"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (s *Service) GetMoneyLimits(ctx context.Context, date time.Time) ([]quik.MoneyLimit, error) {
	maxDate, err := s.repo.SelectMoneyLimitsMaxDate(ctx)
	if err != nil {
		return nil, err
	}
	if maxDate == nil {
		return []quik.MoneyLimit{}, nil
	}
	if err := checkLimitDate(date, *maxDate); err != nil {
		return nil, err
	}
	return s.repo.SelectMoneyLimits(ctx, date)
}

func (s *Service) DoRollForwardMoneyLimits(ctx context.Context) error {
	return doRollForward(ctx,
		s.repo.SelectMoneyLimitsMaxDate,
		s.repo.InsertMoneyLimitsCopy,
		s.repo.DeleteMoneyLimits,
	)
}
