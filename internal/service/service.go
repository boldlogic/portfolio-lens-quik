package service

import (
	"context"
	"time"

	qmodels "github.com/boldlogic/quik-portfolio/internal/models"
	"github.com/boldlogic/quik-portfolio/pkg/models/quik"
	"go.uber.org/zap"
)

type MoneyLimitsRepo interface {
	SelectMoneyLimits(ctx context.Context, date time.Time) ([]qmodels.MoneyLimit, error)
	InsertMoneyLimit(ctx context.Context, s qmodels.MoneyLimit) (qmodels.MoneyLimit, error)
	SelectMoneyLimitsMaxDate(ctx context.Context) (*time.Time, error)
	InsertMoneyLimitsCopy(ctx context.Context, dateFrom time.Time, dateTo time.Time) error
	DeleteMoneyLimits(ctx context.Context, date time.Time) error
}

type SecurityLimitsRepo interface {
	SelectSecurityLimits(ctx context.Context, date time.Time) ([]qmodels.SecurityLimit, error)
	InsertSecurityLimit(ctx context.Context, s qmodels.SecurityLimit) (qmodels.SecurityLimit, error)
	SelectSecurityLimitsMaxDate(ctx context.Context) (*time.Time, error)
	InsertSecurityLimitsCopy(ctx context.Context, dateFrom time.Time, dateTo time.Time) error
	DeleteSecurityLimits(ctx context.Context, date time.Time) error
}

type SecurityLimitsOtcRepo interface {
	SelectSecurityLimitsOtc(ctx context.Context, date time.Time) ([]qmodels.SecurityLimit, error)
	InsertSecurityLimitOtc(ctx context.Context, s qmodels.SecurityLimit) (qmodels.SecurityLimit, error)
	SelectSecurityLimitsOtcMaxDate(ctx context.Context) (*time.Time, error)
	InsertSecurityLimitsOtcCopy(ctx context.Context, dateFrom time.Time, dateTo time.Time) error
	DeleteSecurityLimitsOtc(ctx context.Context, date time.Time) error
}

type PortfolioRepo interface {
	SelectSecuritiesPortfolio(ctx context.Context, date time.Time, targetCcy string) ([]qmodels.PortfolioEntry, error)
	SelectSecuritiesOtcPortfolio(ctx context.Context, date time.Time, targetCcy string) ([]qmodels.PortfolioEntry, error)
	SelectMoneyLimitsPortfolio(ctx context.Context, date time.Time, targetCcy string) ([]qmodels.PortfolioEntry, error)
}

type FirmsRepo interface {
	SelectFirms(ctx context.Context) ([]quik.Firm, error)
	SelectFirmByID(ctx context.Context, id uint8) (quik.Firm, error)
	SelectFirmByName(ctx context.Context, name string) (quik.Firm, error)
	InsertFirm(ctx context.Context, code string, name string) (quik.Firm, error)
	UpdateFirm(ctx context.Context, id uint8, name string) (quik.Firm, error)
	SyncFirmsFromLimits(ctx context.Context) error
}

type CurrentQuotesRepo interface {
	SelectCurrentQuotes(ctx context.Context) ([]qmodels.CurrentQuote, error)
	SelectCurrentQuotesForKeys(ctx context.Context, keys []string) ([]qmodels.CurrentQuote, error)
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
