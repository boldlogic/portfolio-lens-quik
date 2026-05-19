package service

import (
	"context"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (s *Service) GetSecurityLimits(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error) {
	return s.repo.SelectSecurityLimits(ctx, date)
}

func (s *Service) GetSecurityLimitsOtc(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error) {
	return s.repo.SelectSecurityLimitsOtc(ctx, date)
}
