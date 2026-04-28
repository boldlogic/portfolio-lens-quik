package workers

import (
	"context"
	"time"

	"github.com/boldlogic/packages/periodic"
	"go.uber.org/zap"
)

type SyncFirmsRunner interface {
	SyncFirmsFromLimits(ctx context.Context) error
}

func NewActualizeFirmsWorker(svc SyncFirmsRunner, logger *zap.Logger, interval time.Duration) periodic.Worker {
	return periodic.NewPeriodicWorker(
		"actualize_firms",
		"ошибка синхронизации фирм из лимитов",
		interval,
		func(ctx context.Context) error { return svc.SyncFirmsFromLimits(ctx) },
		logger,
	)
}
