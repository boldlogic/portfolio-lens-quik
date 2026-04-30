package writefirms

import (
	"context"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	repo   firmsRepo
}

type firmsRepo interface {
	InsertFirm(ctx context.Context, code string, name string) (quik.Firm, error)
	UpdateFirm(ctx context.Context, id uint8, name string) (quik.Firm, error)
}

func NewService(repo firmsRepo, logger *zap.Logger) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}
