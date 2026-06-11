package dictionary

import (
	"context"
	"time"

	"github.com/boldlogic/packages/periodic"
	"go.uber.org/zap"
)

type UpdateCurrencyDictionaryRunner interface {
	UpdateCurrencyDictionary(ctx context.Context) error
}

func NewUpdateCurrencyDictionaryWorker(svc UpdateCurrencyDictionaryRunner, logger *zap.Logger, name string, interval time.Duration) periodic.Worker {
	if logger == nil {
		logger = zap.NewNop()
	}
	return periodic.NewPeriodicWorker(
		name,
		"ошибка обновления справочника валют",
		interval,
		func(ctx context.Context) error { return svc.UpdateCurrencyDictionary(ctx) },
		logger,
	)
}
