package service

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type Repository interface {
	SelectMoneyLimitsMaxDate(ctx context.Context) (*time.Time, error)
	InsertMoneyLimitsCopy(ctx context.Context, dateFrom time.Time, dateTo time.Time) error
	DeleteMoneyLimits(ctx context.Context, date time.Time) error
	SelectSecurityLimitsMaxDate(ctx context.Context) (*time.Time, error)
	InsertSecurityLimitsCopy(ctx context.Context, dateFrom time.Time, dateTo time.Time) error
	DeleteSecurityLimits(ctx context.Context, date time.Time) error
	SelectSecurityLimitsOtcMaxDate(ctx context.Context) (*time.Time, error)
	InsertSecurityLimitsOtcCopy(ctx context.Context, dateFrom time.Time, dateTo time.Time) error
	DeleteSecurityLimitsOtc(ctx context.Context, date time.Time) error
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
