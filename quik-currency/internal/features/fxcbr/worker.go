package fxcbr

import (
	"context"
	"time"

	"github.com/boldlogic/packages/periodic"
	"go.uber.org/zap"
)

type MergeFxCBRRatesQuikRunner interface {
	MergeFxCBRRatesQuik(ctx context.Context) error
}

func NewMergeFxCBRRatesQuikWorker(svc MergeFxCBRRatesQuikRunner, logger *zap.Logger, name string, interval time.Duration) periodic.Worker {
	if logger == nil {
		logger = zap.NewNop()
	}
	return periodic.NewPeriodicWorker(
		name,
		"ошибка сохранения кросс-курсов валют из QUIK",
		interval,
		func(ctx context.Context) error { return svc.MergeFxCBRRatesQuik(ctx) },
		logger,
	)
}
