package service

import (
	"context"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"go.uber.org/zap"
)

type MoneyLimitsRepo interface {
	ListMoneyLimits(ctx context.Context, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.MoneyLimit, totalCount *uint64, err error)
}

type SecurityLimitsRepo interface {
	ListSecurityLimits(ctx context.Context, limitType quik.LimitType, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.SecurityLimit, totalCount *uint64, err error)
}

type PortfolioRepo interface {
	ListMoneyPortfolio(ctx context.Context, date time.Time, targetCcy string, clientCodes []string) (result []quik.Position, err error)
	ListSecurityPortfolio(ctx context.Context, date time.Time, targetCcy string, clientCodes []string) (result []quik.Position, err error)
	ListSecurityPortfolioOtc(ctx context.Context, date time.Time, targetCcy string, clientCodes []string) (result []quik.Position, err error)
}

type Repository interface {
	MoneyLimitsRepo
	SecurityLimitsRepo
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
