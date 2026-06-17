package marketdata

import (
	"context"

	"github.com/boldlogic/packages/shutdown"
	"go.uber.org/zap"
)

type RateImporter struct {
	logger *zap.Logger
	store  RateStore
}

type RateStore interface {
	MergeFxCBRRates(ctx context.Context) error
}

func NewRateImporter(logger *zap.Logger, store RateStore) *RateImporter {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &RateImporter{
		logger: logger,
		store:  store,
	}
}

// TO:DO  разделяем на операции чтения и записи
// чтение по разным класс-кодам летит в два канала: цб пишет в два, биржа - в один
// курсы ЦБ летят в старую таблицу и новую
// биржевые курсы летят только в новую
func (i *RateImporter) ImportQuikCrossRates(ctx context.Context) error {
	if err := i.store.MergeFxCBRRates(ctx); err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		i.logger.Error("произошла ошибка при сохранении кросс-курсов валют из QUIK", zap.Error(err))
		return err
	}
	return nil
}
