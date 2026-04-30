package db

import (
	"context"

	"github.com/boldlogic/packages/dbzap"
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
