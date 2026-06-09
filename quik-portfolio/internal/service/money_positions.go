package service

import (
	"context"

	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

func (s *Service) GetMoneyPositions(ctx context.Context, targetCurrency *string, clientCodes []string) ([]quik.Position, decimal.Decimal, string, error) {
	ccy, clients, err := normalizePortfolioRequest(targetCurrency, clientCodes)
	if err != nil {
		return nil, decimal.Decimal{}, "", err
	}

	date := dates.Today()

	positions, err := s.repo.ListMoneyPortfolio(ctx, date, ccy, clients)
	if err != nil {
		return nil, decimal.Decimal{}, "", err
	}

	return positions, sumTotalMarketValue(positions), ccy, nil
}

func sumTotalMarketValue(positions []quik.Position) decimal.Decimal {

	var total decimal.Decimal
	for _, pos := range positions {
		total = decimal.Sum(total, pos.MarketValueInTargetCurrency)

	}

	return total
}
