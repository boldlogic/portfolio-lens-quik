package syncclients

import (
	"context"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	repo   clientsRepo
}

type clientsRepo interface {
	SyncClientsFromLimits(ctx context.Context) error
}

func NewService(repo clientsRepo, logger *zap.Logger) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}

func (s *Service) SyncClientsFromLimits(ctx context.Context) error {
	return s.repo.SyncClientsFromLimits(ctx)
}
