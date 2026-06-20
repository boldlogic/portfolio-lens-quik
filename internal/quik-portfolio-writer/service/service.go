package service

import (
	"context"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"go.uber.org/zap"
)

type Repository interface {
	HandleRequest(ctx context.Context, limits []quik.Limit) error
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
