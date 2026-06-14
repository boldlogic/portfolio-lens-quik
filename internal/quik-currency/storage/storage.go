package storage

import (
	"context"

	"github.com/boldlogic/packages/dbconfig"
	"github.com/boldlogic/packages/dbzap"
	"go.uber.org/zap"
)

type Storage struct {
	*dbzap.Pool
}

func NewStorage(ctx context.Context, cfg dbconfig.DBConfig, logger *zap.Logger) (*Storage, error) {
	pool, err := dbzap.New(ctx, cfg.GetDSN(), logger)
	if err != nil {
		return nil, err
	}
	return &Storage{Pool: pool}, nil
}
