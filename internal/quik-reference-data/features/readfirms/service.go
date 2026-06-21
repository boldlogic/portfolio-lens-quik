package readfirms

import (
	"context"

	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	"go.uber.org/zap"
)

type service struct {
	logger *zap.Logger
	repo   firmsRepo
}

type firmsRepo interface {
	SelectFirms(ctx context.Context) ([]quik.Firm, error)
	SelectFirmByID(ctx context.Context, id uint8) (quik.Firm, error)
}

func NewService(repo firmsRepo, logger *zap.Logger) *service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}
