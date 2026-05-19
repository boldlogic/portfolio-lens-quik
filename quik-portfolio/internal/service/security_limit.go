package service

import (
	"context"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (s *Service) GetSecurityLimits(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error) {
	maxDate, err := s.repo.SelectSecurityLimitsMaxDate(ctx)
	if err != nil {
		return nil, err
	}
	if maxDate == nil {
		return []quik.SecurityLimit{}, nil
	}
	if err := checkLimitDate(date, *maxDate); err != nil {
		return nil, err
	}
	return s.repo.SelectSecurityLimits(ctx, date)
}

func (s *Service) GetSecurityLimitsOtc(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error) {
	maxDate, err := s.repo.SelectSecurityLimitsOtcMaxDate(ctx)
	if err != nil {
		return nil, err
	}
	if maxDate == nil {
		return []quik.SecurityLimit{}, nil
	}
	if err := checkLimitDate(date, *maxDate); err != nil {
		return nil, err
	}
	return s.repo.SelectSecurityLimitsOtc(ctx, date)
}
