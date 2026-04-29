package repository

import (
	"context"

	"github.com/boldlogic/packages/dbzap"
	"github.com/boldlogic/packages/shutdown"
	"go.uber.org/zap"
)

type Repository struct {
	*dbzap.Pool
}

func NewRepository(ctx context.Context, dsn string, logger *zap.Logger) (*Repository, error) {
	pool, err := dbzap.New(ctx, dsn, logger)
	if err != nil {
		return nil, err
	}
	return &Repository{Pool: pool}, nil
}

func (r *Repository) isShutdown(err error) bool {
	return shutdown.IsExceeded(err)
}
