package reference

import (
	"context"

	"github.com/JohannesJHN/iso4217"
	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-currency/pkg/currencies"
	"go.uber.org/zap"
)

type CurrencyCatalog struct {
	logger *zap.Logger
	store  CurrencyCatalogStore
}

type CurrencyCatalogStore interface {
	SelectNewCurrenciesFromCrossrates(ctx context.Context) ([]currencies.Currency, error)
	MergeCurrencies(ctx context.Context, currencies []currencies.Currency) error
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

	out := make([]currencies.Currency, 0, len(raw))

	dedup := make(map[currencies.CurrencyCode]struct{})
	for _, ccy := range raw {
		if _, ok := dedup[ccy.ISOCharCode]; ok {
			continue
		}

		dedup[ccy.ISOCharCode] = struct{}{}
		libCcy, ok := iso4217.LookupByAlpha3(ccy.ISOCharCode.String())
		if !ok {
			c.logger.Warn("не нашли в либе", zap.String("ISOCharCode", ccy.ISOCharCode.String()))
			continue
		}

		var miu int32 = int32(libCcy.MinorUnits)
		if miu < 0 {
			miu = 0
		}
		curr := currencies.Currency{
			ISOCode:     int16(libCcy.Numeric),
			ISOCharCode: ccy.ISOCharCode,
			Name:        ccy.Name,
			LatName:     libCcy.Name,
			MinorUnits:  miu,
		}
		out = append(out, curr)
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
