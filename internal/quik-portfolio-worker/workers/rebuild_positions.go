package workers

import (
	"context"
	"time"

	"github.com/boldlogic/packages/periodic"
	"go.uber.org/zap"
)

type RebuildPositionsRunner interface {
	RebuildPositions(ctx context.Context) error
}

func NewRebuildPositionsWorker(svc RebuildPositionsRunner, logger *zap.Logger, interval time.Duration) periodic.Worker {
	return periodic.NewPeriodicWorker(
		"rebuild_positions",
		"ошибка сборки позиций",
		interval,
		func(ctx context.Context) error { return svc.RebuildPositions(ctx) },
		logger,
	)
}
