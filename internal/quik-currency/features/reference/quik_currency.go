package reference

import (
	"context"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-currency/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"go.uber.org/zap"
)

type CurrencyCatalog struct {
	logger *zap.Logger
	store  CurrencyCatalogStore
}

type CurrencyCatalogStore interface {
	SelectNewCurrenciesFromCrossrates(ctx context.Context) ([]models.CurrencyFromCrossrates, error)
	MergeCurrencies(ctx context.Context, currencies []quik.Currency) error
}

func NewCurrencyCatalog(logger *zap.Logger, store CurrencyCatalogStore) *CurrencyCatalog {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &CurrencyCatalog{
		logger: logger,
		store:  store,
	}
}

func (c *CurrencyCatalog) logError(msg string, err error) {
	if shutdown.IsExceeded(err) {
		return
	}
	if err != nil {
		c.logger.Error(msg, zap.Error(err))
	}
}

func (c *CurrencyCatalog) RefreshQuikCurrencies(ctx context.Context) (err error) {

	var msg string
	defer func() { c.logError(msg, err) }()

	raw, err := c.store.SelectNewCurrenciesFromCrossrates(ctx)
	if err != nil {
		msg = "ошибка при получении кросс-курсов"
		return err
	}
	if len(raw) == 0 {
		return nil
	}

	out := make([]quik.Currency, 0, len(raw))

	dedup := make(map[quik.CurrencyCode]struct{})

	for _, ccy := range raw {

		currency, err := quik.CurrencyFromQuik(ccy.IsoCharCode, &ccy.Name)
		if err != nil {
			c.logger.Warn("не нашли в либе", zap.String("ISOCharCode", ccy.IsoCharCode))
			continue
		}

		if _, ok := dedup[currency.Alpha()]; ok {
			continue
		}

		dedup[currency.Alpha()] = struct{}{}

		out = append(out, currency)
	}

	if len(out) == 0 {
		return nil
	}

	err = c.store.MergeCurrencies(ctx, out)
	if err != nil {
		msg = "ошибка при сохранении справочника валют"
		return err
	}
	return nil
}
