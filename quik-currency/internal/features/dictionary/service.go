package dictionary

import (
	"context"

	"github.com/JohannesJHN/iso4217"
	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-currency/pkg/currencies"
	"go.uber.org/zap"
)

type service struct {
	logger *zap.Logger
	repo   currencyDictionaryRepo
}
type currencyDictionaryRepo interface {
	SelectNewCurrenciesFromCrossrates(ctx context.Context) ([]currencies.Currency, error)
	MergeCurrencies(ctx context.Context, currencies []currencies.Currency) error
}

func NewService(logger *zap.Logger, repo currencyDictionaryRepo) *service {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &service{
		logger: logger,
		repo:   repo,
	}
}
func (s service) logError(msg string, err error) {
	if shutdown.IsExceeded(err) {
		return
	}
	if err != nil {
		s.logger.Error(msg, zap.Error(err))
	}
}

func (s service) UpdateCurrencyDictionary(ctx context.Context) (err error) {
	var msg string
	defer func() { s.logError(msg, err) }()

	raw, err := s.repo.SelectNewCurrenciesFromCrossrates(ctx)
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
			s.logger.Warn("не нашли в либе", zap.String("ISOCharCode", ccy.ISOCharCode.String()))

			continue
		}

		var miu int32 = int32(libCcy.MinorUnits)
		if miu < 0 {
			miu = 0
		}
		c := currencies.Currency{
			ISOCode:     int16(libCcy.Numeric),
			ISOCharCode: ccy.ISOCharCode,
			Name:        ccy.Name,
			LatName:     libCcy.Name,
			MinorUnits:  miu,
		}
		out = append(out, c)
	}

	if len(out) == 0 {
		return nil
	}

	err = s.repo.MergeCurrencies(ctx, out)
	if err != nil {
		msg = "ошибка при сохранении справочника валют"

		return err
	}
	return nil
}
