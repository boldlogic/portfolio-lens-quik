package syncinstrumentcodes

import (
	"context"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	repo   instrumentCodesRepo
}

type instrumentCodesRepo interface {
	SyncInstrumentCodes(ctx context.Context) error
}

func NewService(repo instrumentCodesRepo, logger *zap.Logger) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}

func (s *Service) SyncInstrumentCodes(ctx context.Context) error {
	return s.repo.SyncInstrumentCodes(ctx)
}
