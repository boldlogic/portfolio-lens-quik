package workers

import (
	"context"
	"time"

	"github.com/boldlogic/packages/periodic"
	"go.uber.org/zap"
)

type MergeFxCBRRatesQuikRunner interface {
	MergeFxCBRRatesQuik(ctx context.Context) error
}

func NewMergeFxCBRRatesQuikWorker(svc MergeFxCBRRatesQuikRunner, logger *zap.Logger, interval time.Duration) periodic.Worker {
	return periodic.NewPeriodicWorker(
		"merge_fx_cbr_rates_quik",
		"ошибка сохранения кросс-курсов валют из QUIK",
		interval,
		func(ctx context.Context) error { return svc.MergeFxCBRRatesQuik(ctx) },
		logger,
	)
}
