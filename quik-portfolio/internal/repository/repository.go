package repository

import (
	"context"

	"github.com/boldlogic/packages/dbconfig"
	"github.com/boldlogic/packages/dbpool"
	"github.com/boldlogic/packages/dbzap"
	"github.com/boldlogic/portfolio-lens-quik/quik-portfolio/internal/observability"
	"go.uber.org/zap"
)

type Repository struct {
	*dbzap.Pool
	metrics observability.Recorder
}

func NewRepository(ctx context.Context, cfg dbconfig.DBConfig, logger *zap.Logger, recorder ...observability.Recorder) (*Repository, error) {
	if logger == nil {
		logger = zap.NewNop()
	}
	pool, err := dbzap.New(ctx, cfg.GetDSN(), logger)
	if err != nil {
		return nil, err
	}
	if err := dbpool.Apply(pool.Db, cfg.Pool); err != nil {
		pool.Close()
		return nil, err
	}
	var rec observability.Recorder
	if len(recorder) > 0 {
		rec = recorder[0]
	}
	return &Repository{
		Pool:    pool,
		metrics: observability.EnsureRecorder(rec),
	}, nil
}
