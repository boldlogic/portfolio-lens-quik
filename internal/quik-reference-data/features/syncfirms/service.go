package syncfirms

import (
	"context"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	repo   firmsRepo
}

type firmsRepo interface {
	SyncFirmsFromLimits(ctx context.Context) error
}

func NewService(repo firmsRepo, logger *zap.Logger) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}

func (s *Service) SyncFirmsFromLimits(ctx context.Context) error {
	return s.repo.SyncFirmsFromLimits(ctx)
}
