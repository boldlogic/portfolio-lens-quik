package service

import (
	"context"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/boldlogic/portfolio-lens-quik/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

type MoneyLimitsRepo interface {
	SelectMoneyLimits(ctx context.Context, date time.Time) ([]quik.MoneyLimit, error)
	SelectMoneyLimitsMaxDate(ctx context.Context) (*time.Time, error)
	InsertMoneyLimitsCopy(ctx context.Context, dateFrom time.Time, dateTo time.Time) error
	DeleteMoneyLimits(ctx context.Context, date time.Time) error
}

type SecurityLimitsRepo interface {
	SelectSecurityLimits(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error)
	SelectSecurityLimitsMaxDate(ctx context.Context) (*time.Time, error)
	InsertSecurityLimitsCopy(ctx context.Context, dateFrom time.Time, dateTo time.Time) error
	DeleteSecurityLimits(ctx context.Context, date time.Time) error
}

type SecurityLimitsOtcRepo interface {
	SelectSecurityLimitsOtc(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error)
	SelectSecurityLimitsOtcMaxDate(ctx context.Context) (*time.Time, error)
	InsertSecurityLimitsOtcCopy(ctx context.Context, dateFrom time.Time, dateTo time.Time) error
	DeleteSecurityLimitsOtc(ctx context.Context, date time.Time) error
}

type PortfolioRepo interface {
	SelectSecuritiesPortfolio(ctx context.Context, date time.Time, targetCcy string) ([]quik.PortfolioEntry, error)
	SelectSecuritiesOtcPortfolio(ctx context.Context, date time.Time, targetCcy string) ([]quik.PortfolioEntry, error)
	SelectMoneyLimitsPortfolio(ctx context.Context, date time.Time, targetCcy string) ([]quik.PortfolioEntry, error)
}

type CurrentQuotesRepo interface {
	SelectCurrentQuotes(ctx context.Context) ([]models.CurrentQuote, error)
	SelectCurrentQuotesForKeys(ctx context.Context, keys []string) ([]models.CurrentQuote, error)
}

type Repository interface {
	MoneyLimitsRepo
	SecurityLimitsRepo
	SecurityLimitsOtcRepo
	PortfolioRepo
	CurrentQuotesRepo
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
