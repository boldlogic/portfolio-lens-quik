package service

import (
	"context"

	"github.com/boldlogic/portfolio-lens-currency/pkg/currencies"
	"go.uber.org/zap"
)

type CurrencyRepository interface {
	MergeCurrencies(ctx context.Context, currencies []currencies.Currency) error

	SelectCountCurrencies(ctx context.Context) (int, error)

	SetEmptyCurrencyNamesFromQuik(ctx context.Context) error

	MergeFxCBRRatesQuik(ctx context.Context) error
}

type Service struct {
	logger *zap.Logger

	currencyRepo CurrencyRepository
}

func NewService(
	currencyRepo CurrencyRepository,
	logger *zap.Logger) *Service {

	s := &Service{
		logger:       logger,
		currencyRepo: currencyRepo,
	}
	return s
}
