package service

import (
	"context"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (s *Service) GetSecurityLimitsWithFilters(ctx context.Context, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.SecurityLimit, totalCount *uint64, err error) {
	dedublicated, err := validateLimitsContract(date, clientCodes)
	if err != nil {
		return nil, nil, err
	}
	return s.repo.ListSecurityLimits(ctx, quik.LimitTypeSecurities, date, limit, offset, dedublicated, includeTotalCount)
}

func (s *Service) GetSecurityLimitsOtcWithFilters(ctx context.Context, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.SecurityLimit, totalCount *uint64, err error) {
	dedublicated, err := validateLimitsContract(date, clientCodes)
	if err != nil {
		return nil, nil, err
	}
	return s.repo.ListSecurityLimits(ctx, quik.LimitTypeSecuritiesOtc, date, limit, offset, dedublicated, includeTotalCount)
}
