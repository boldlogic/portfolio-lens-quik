package service

import (
	"context"
	"fmt"

	"github.com/JohannesJHN/iso4217"
	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-currency/pkg/currencies"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (s *Service) GetSecurityPositions(ctx context.Context, currency *string) ([]quik.Position, decimal.Decimal, string, error) {
	if currency == nil {
		currency = new("RUB")
	} else {
		ccy, err := currencies.ParseCurrencyCode(*currency)
		if err != nil {
			return nil, decimal.Decimal{}, "", fmt.Errorf("%w: %w", models.ErrBusinessValidation, err)
		}
		currency = new(ccy.String())
	}
	c, ok := iso4217.LookupByAlpha3(*currency)
	var minorUnits int32 = 2
	if ok {
		minorUnits = int32(c.MinorUnits)
	}
	s.logger.Debug("mi", zap.Int32("minorUnits", minorUnits))
	date := dates.Today()

	positions, err := s.repo.SelectSecPortfolio(ctx, date, *currency)
	if err != nil {
		return nil, decimal.Decimal{}, "", err
	}

	return positions, sumTotalMarketValue(positions, minorUnits), *currency, nil
}
