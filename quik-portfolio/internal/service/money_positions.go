package service

import (
	"context"
	"fmt"

	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-currency/pkg/currencies"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

func (s *Service) GetMoneyPositions(ctx context.Context, currency *string) ([]quik.Position, decimal.Decimal, string, error) {
	if currency == nil {
		currency = new("RUB")
	} else {
		ccy, err := currencies.ParseCurrencyCode(*currency)
		if err != nil {
			return nil, decimal.Decimal{}, "", fmt.Errorf("%w: %w", models.ErrBusinessValidation, err)
		}
		currency = new(ccy.String())
	}

	date := dates.Today()

	positions, err := s.repo.SelectMoneyPortfolio(ctx, date, *currency)
	if err != nil {
		return nil, decimal.Decimal{}, "", err
	}

	var total decimal.Decimal
	for _, pos := range positions {
		total = decimal.Sum(total, pos.MVTotal)

	}

	return positions, total, *currency, nil
}
