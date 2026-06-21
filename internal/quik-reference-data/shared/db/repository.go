package db

import (
	"context"

	"github.com/boldlogic/packages/dbconfig"
	"github.com/boldlogic/packages/dbpool"
	"github.com/boldlogic/packages/dbzap"
	"go.uber.org/zap"
)

type Repository struct {
	*dbzap.Pool
}

func NewRepository(ctx context.Context, cfg dbconfig.DBConfig, logger *zap.Logger) (*Repository, error) {
	pool, err := dbzap.New(ctx, cfg.GetDSN(), logger)
	if err != nil {
		return nil, err
	}
	if err := dbpool.Apply(pool.Db, cfg.Pool); err != nil {
		pool.Close()
		return nil, err
	}
	return &Repository{Pool: pool}, nil
}
