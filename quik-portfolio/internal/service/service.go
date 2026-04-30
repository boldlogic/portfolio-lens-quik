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
	InsertMoneyLimit(ctx context.Context, s quik.MoneyLimit) (quik.MoneyLimit, error)
	SelectMoneyLimitsMaxDate(ctx context.Context) (*time.Time, error)
	InsertMoneyLimitsCopy(ctx context.Context, dateFrom time.Time, dateTo time.Time) error
	DeleteMoneyLimits(ctx context.Context, date time.Time) error
}

type SecurityLimitsRepo interface {
	SelectSecurityLimits(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error)
	InsertSecurityLimit(ctx context.Context, s quik.SecurityLimit) (quik.SecurityLimit, error)
	SelectSecurityLimitsMaxDate(ctx context.Context) (*time.Time, error)
	InsertSecurityLimitsCopy(ctx context.Context, dateFrom time.Time, dateTo time.Time) error
	DeleteSecurityLimits(ctx context.Context, date time.Time) error
}

type SecurityLimitsOtcRepo interface {
	SelectSecurityLimitsOtc(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error)
	InsertSecurityLimitOtc(ctx context.Context, s quik.SecurityLimit) (quik.SecurityLimit, error)
	SelectSecurityLimitsOtcMaxDate(ctx context.Context) (*time.Time, error)
	InsertSecurityLimitsOtcCopy(ctx context.Context, dateFrom time.Time, dateTo time.Time) error
	DeleteSecurityLimitsOtc(ctx context.Context, date time.Time) error
}

type PortfolioRepo interface {
	SelectSecuritiesPortfolio(ctx context.Context, date time.Time, targetCcy string) ([]quik.PortfolioEntry, error)
	SelectSecuritiesOtcPortfolio(ctx context.Context, date time.Time, targetCcy string) ([]quik.PortfolioEntry, error)
	SelectMoneyLimitsPortfolio(ctx context.Context, date time.Time, targetCcy string) ([]quik.PortfolioEntry, error)
}

type FirmsRepo interface {
	SelectFirms(ctx context.Context) ([]quik.Firm, error)
	SelectFirmByID(ctx context.Context, id uint8) (quik.Firm, error)
	SelectFirmByName(ctx context.Context, name string) (quik.Firm, error)
	InsertFirm(ctx context.Context, code string, name string) (quik.Firm, error)
	UpdateFirm(ctx context.Context, id uint8, name string) (quik.Firm, error)
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
	FirmsRepo
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
