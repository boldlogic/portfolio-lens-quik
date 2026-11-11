package syncclients

import (
	"context"
	"time"

	"github.com/boldlogic/packages/periodic"
	"go.uber.org/zap"
)

type SyncClientsRunner interface {
	SyncClientsFromLimits(ctx context.Context) error
}

func NewActualizeClientsWorker(svc SyncClientsRunner, logger *zap.Logger, interval time.Duration) periodic.Worker {
	return periodic.NewPeriodicWorker(
		"actualize_clients",
		"ошибка синхронизации клиентов из лимитов",
		interval,
		func(ctx context.Context) error { return svc.SyncClientsFromLimits(ctx) },
		logger,
	)
}
