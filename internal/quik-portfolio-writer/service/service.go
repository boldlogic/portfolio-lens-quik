package service

import (
	"context"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"go.uber.org/zap"
)

type Repository interface {
	InsertMoneyLimit(ctx context.Context, s quik.MoneyLimit) (quik.MoneyLimit, error)
	InsertSecurityLimit(ctx context.Context, s quik.SecurityLimit) (quik.SecurityLimit, error)
	InsertSecurityLimitOtc(ctx context.Context, s quik.SecurityLimit) (quik.SecurityLimit, error)
	InsertLimit(ctx context.Context, l quik.Limit) (quik.Limit, error)
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
