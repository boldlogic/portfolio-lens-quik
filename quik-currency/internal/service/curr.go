package service

import (
	"context"

	"github.com/JohannesJHN/iso4217"
	"github.com/boldlogic/portfolio-lens-currency/pkg/currencies"
	"go.uber.org/zap"
)

func (s *Service) InitCurrencyDictionary(ctx context.Context) error {
	return s.getNewCurrenciesFromLib(ctx)
}

func (s *Service) getNewCurrenciesFromLib(ctx context.Context) error {

	count, err := s.currencyRepo.SelectCountCurrencies(ctx)
	if err != nil {
		return err
	}
	if count != 0 {
		s.logger.Debug("в справочнике уже есть валюты, библиотеку не используем", zap.Int("количество записей", count))
		return nil
	}

	lib := iso4217.AllActive()

	ccs := make([]currencies.Currency, 0, len(lib))

	for k, v := range lib {
		ccs = append(ccs, currencies.Currency{
			ISOCode:     int16(v.Numeric),
			ISOCharCode: currencies.CurrencyCode(k),
			LatName:     v.Name,
			MinorUnits:  int32(v.MinorUnits),
		})
	}

	err = s.currencyRepo.MergeCurrencies(ctx, ccs)
	if err != nil {
		s.logger.Error("произошла ошибка при добавлении валют из библиотеки", zap.Error(err))

		return err
	}
	s.logger.Info("справочник валют был пуст. добавлены валюты из библиотеки", zap.Int("количество записей", len(ccs)))

	return nil
}
