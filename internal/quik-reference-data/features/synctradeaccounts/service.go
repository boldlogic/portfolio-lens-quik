package synctradeaccounts

import (
	"context"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	repo   tradeAccountsRepo
}

type tradeAccountsRepo interface {
	SyncTradeAccountsFromLimits(ctx context.Context) error
}

func NewService(repo tradeAccountsRepo, logger *zap.Logger) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}

func (s *Service) SyncTradeAccountsFromLimits(ctx context.Context) error {
	return s.repo.SyncTradeAccountsFromLimits(ctx)
}
