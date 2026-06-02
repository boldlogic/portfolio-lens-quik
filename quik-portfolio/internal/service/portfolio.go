package service

import (
	"context"
	"fmt"

	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-currency/pkg/currencies"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"

	"golang.org/x/sync/errgroup"
)

func (s *Service) GetPortfolio(ctx context.Context, targetCcy string) ([]quik.PortfolioEntry, error) {
	if targetCcy == "" {
		targetCcy = "RUB"
	}
	ccy, err := currencies.ParseCurrencyCode(targetCcy)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", models.ErrBusinessValidation, err)
	}
	targetCcy = ccy.String()

	var securities []quik.PortfolioEntry
	var otcEntries []quik.PortfolioEntry
	var moneyEntries []quik.PortfolioEntry

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		securities, err = s.repo.SelectSecuritiesPortfolio(gCtx, dates.Today(), targetCcy)
		return err
	})

	g.Go(func() error {
		var err error
		otcEntries, err = s.repo.SelectSecuritiesOtcPortfolio(gCtx, dates.Today(), targetCcy)
		return err
	})

	g.Go(func() error {
		var err error
		moneyEntries, err = s.repo.SelectMoneyLimitsPortfolio(gCtx, dates.Today(), targetCcy)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	entries := make([]quik.PortfolioEntry, 0, len(securities)+len(otcEntries)+len(moneyEntries))
	entries = append(entries, securities...)
	entries = append(entries, otcEntries...)
	entries = append(entries, moneyEntries...)

	return entries, nil
}
