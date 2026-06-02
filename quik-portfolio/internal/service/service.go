package service

import (
	"context"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"go.uber.org/zap"
)

type MoneyLimitsRepo interface {
	SelectMoneyLimitsWithFilters(ctx context.Context, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.MoneyLimit, totalCount *uint64, err error)
}

type SecurityLimitsRepo interface {
	SelectSecurityLimitsWithFilters(ctx context.Context, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.SecurityLimit, totalCount *uint64, err error)
}

type SecurityLimitsOtcRepo interface {
	SelectSecurityLimitsOtcWithFilters(ctx context.Context, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.SecurityLimit, totalCount *uint64, err error)
}

type PortfolioRepo interface {
	SelectSecuritiesPortfolio(ctx context.Context, date time.Time, targetCcy string) ([]quik.PortfolioEntry, error)
	SelectSecuritiesOtcPortfolio(ctx context.Context, date time.Time, targetCcy string) ([]quik.PortfolioEntry, error)
	SelectMoneyLimitsPortfolio(ctx context.Context, date time.Time, targetCcy string) ([]quik.PortfolioEntry, error)
}

type Repository interface {
	MoneyLimitsRepo
	SecurityLimitsRepo
	SecurityLimitsOtcRepo
	PortfolioRepo
}

type Service struct {
	logger *zap.Logger
	repo   Repository
}

func NewService(repo Repository, logger *zap.Logger) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}
