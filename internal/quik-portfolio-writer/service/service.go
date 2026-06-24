package service

import (
	"context"

	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	"go.uber.org/zap"
)

type Repository interface {
	HandleRequest(ctx context.Context, limits []quik.Limit) error
}

type Service struct {
	logger *zap.Logger
	worker
}

func NewService(repo Repository, cfg WorkerConfig, logger *zap.Logger) *Service {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Service{
		logger: logger,
		worker: *newWorker(
			repo, logger, int(cfg.BatchSize), cfg.QueueSize, cfg.Interval,
		),
	}
}
