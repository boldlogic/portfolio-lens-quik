package service

import (
	"context"

	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

func (s *Service) GetSecurityPositions(ctx context.Context, targetCurrency *string, clientCodes []string) ([]quik.Position, decimal.Decimal, string, error) {
	return s.getSecurityPositions(ctx, targetCurrency, clientCodes, false)
}

func (s *Service) GetSecurityPositionsOtc(ctx context.Context, targetCurrency *string, clientCodes []string) ([]quik.Position, decimal.Decimal, string, error) {
	return s.getSecurityPositions(ctx, targetCurrency, clientCodes, true)
}

func (s *Service) getSecurityPositions(ctx context.Context, targetCurrency *string, clientCodes []string, otc bool) ([]quik.Position, decimal.Decimal, string, error) {
	ccy, clients, err := normalizePortfolioRequest(targetCurrency, clientCodes)
	if err != nil {
		return nil, decimal.Decimal{}, "", err
	}

	date := dates.Today()

	var positions []quik.Position
	if otc {
		positions, err = s.repo.ListSecurityPortfolioOtc(ctx, date, ccy, clients)
	} else {
		positions, err = s.repo.ListSecurityPortfolio(ctx, date, ccy, clients)
	}
	if err != nil {
		return nil, decimal.Decimal{}, "", err
	}

	return positions, sumTotalMarketValue(positions), ccy, nil
}
