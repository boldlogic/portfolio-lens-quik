package service

import (
	"context"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (s *Service) GetSecurityLimits(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error) {
	return s.repo.SelectSecurityLimits(ctx, date)
}

func (s *Service) GetSecurityLimitsWithFilters(ctx context.Context, date time.Time, limit, offset int, clientCodes []string) ([]quik.SecurityLimit, int, error) {
	return s.repo.SelectSecurityLimitsWithFilters(ctx, date, limit, offset, clientCodes)
}

func (s *Service) GetSecurityLimitsOtc(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error) {
	return s.repo.SelectSecurityLimitsOtc(ctx, date)
}

func (s *Service) GetSecurityLimitsOtcWithFilters(ctx context.Context, date time.Time, limit, offset int, clientCodes []string) ([]quik.SecurityLimit, int, error) {
	return s.repo.SelectSecurityLimitsOtcWithFilters(ctx, date, limit, offset, clientCodes)
}
