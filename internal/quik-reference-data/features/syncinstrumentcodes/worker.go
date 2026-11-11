package syncinstrumentcodes

import (
	"context"
	"time"

	"github.com/boldlogic/packages/periodic"
	"go.uber.org/zap"
)

type SyncInstrumentCodesRunner interface {
	SyncInstrumentCodes(ctx context.Context) error
}

func NewActualizeInstrumentCodesWorker(svc SyncInstrumentCodesRunner, logger *zap.Logger, interval time.Duration) periodic.Worker {
	return periodic.NewPeriodicWorker(
		"actualize_instrument_codes",
		"ошибка синхронизации кодов инструментов",
		interval,
		func(ctx context.Context) error { return svc.SyncInstrumentCodes(ctx) },
		logger,
	)
}
