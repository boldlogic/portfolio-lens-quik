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
	logger      *zap.Logger
	repo        Repository
	limitsQueue chan quik.Limit
}

func NewService(repo Repository, logger *zap.Logger) *Service {
	return &Service{
		logger:      logger,
		repo:        repo,
		limitsQueue: make(chan quik.Limit, 100),
	}
}
