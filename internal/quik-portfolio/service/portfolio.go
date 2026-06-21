package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
)

func normalizePortfolioRequest(targetCurrency *string, clientCodes []string) (ccy string, clients []string, err error) {
	ccy = "RUB"
	if targetCurrency != nil {
		normCcy, parseErr := quik.ParseCurrencyCode(*targetCurrency)
		if parseErr != nil {
			return "", nil, fmt.Errorf("%w: %w", models.ErrBusinessValidation, parseErr)
		}
		ccy = normCcy.String()
	}

	clients, err = deduplicateClientCodes(clientCodes)
	if err != nil {
		return "", nil, err
	}
	return ccy, clients, nil
}

func sumTotalMarketValue(positions []quik.Position) decimal.Decimal {

	var total decimal.Decimal
	for _, pos := range positions {
		total = decimal.Sum(total, pos.MarketValueInTargetCurrency)

	}

	return total
}

func (s *Service) GetPositions(ctx context.Context, targetCurrency *string, clientCodes []string) ([]quik.Position, decimal.Decimal, string, error) {
	ccy, clients, err := normalizePortfolioRequest(targetCurrency, clientCodes)
	if err != nil {
		return nil, decimal.Decimal{}, "", err
	}

	date := dates.Today()

	g, gCTX := errgroup.WithContext(ctx)

	var money []quik.Position
	g.Go(func() error {
		var err error
		money, err = s.repo.ListMoneyPortfolio(gCTX, date, ccy, clients)
		return err
	})

	var sec []quik.Position
	g.Go(func() error {
		var err error
		sec, err = s.repo.ListSecurityPortfolio(gCTX, date, ccy, clients)
		return err
	})
	var otc []quik.Position
	g.Go(func() error {
		var err error
		otc, err = s.repo.ListSecurityPortfolioOtc(gCTX, date, ccy, clients)
		return err
	})
	err = g.Wait()
	if err != nil {
		return nil, decimal.Decimal{}, "", err
	}

	res := make([]quik.Position, 0, len(money)+len(sec)+len(otc))
	res = slices.Concat(money, sec, otc)

	return res, sumTotalMarketValue(res), ccy, nil
}
