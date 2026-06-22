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
	repo   Repository

	worker
}

func NewService(repo Repository, logger *zap.Logger) *Service {

	return &Service{
		logger: logger,
		repo:   repo,
		worker: *newWorker(
			repo, logger, 100, 100, 2000,
		),
	}
}
