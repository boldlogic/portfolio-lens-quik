package firms

import (
	"context"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	repo   FirmsRepo
}

func NewService(repo FirmsRepo, logger *zap.Logger) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}

func (s *Service) SyncFirmsFromLimits(ctx context.Context) error {
	return s.repo.SyncFirmsFromLimits(ctx)
}
