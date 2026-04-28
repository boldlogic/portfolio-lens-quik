package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/boldlogic/packages/utils/dates"
	qmodels "github.com/boldlogic/portfolio-lens-quik/quik-portfolio/internal/models"
	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (s *Service) GetLimits(ctx context.Context, date time.Time) ([]qmodels.Limit, error) {
	var res []qmodels.Limit
	var ml []qmodels.MoneyLimit
	var sl []qmodels.SecurityLimit
	var otc []qmodels.SecurityLimit
	var maxDateMoney, maxDateSec, maxDateOtc *time.Time

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		maxDateMoney, err = s.repo.SelectMoneyLimitsMaxDate(gCtx)
		if errors.Is(err, md.ErrNotFound) {
			return nil
		}
		return err
	})
	g.Go(func() error {
		var err error
		maxDateSec, err = s.repo.SelectSecurityLimitsMaxDate(gCtx)
		if errors.Is(err, md.ErrNotFound) {
			return nil
		}
		return err
	})
	g.Go(func() error {
		var err error
		maxDateOtc, err = s.repo.SelectSecurityLimitsOtcMaxDate(gCtx)
		if errors.Is(err, md.ErrNotFound) {
			return nil
		}
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	minAvailable := dates.EarliestDate(maxDateMoney, maxDateSec, maxDateOtc)
	if minAvailable == nil {
		return []qmodels.Limit{}, nil
	}
	if err := checkLimitDate(date, *minAvailable); err != nil {
		return nil, err
	}

	g, gCtx = errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		ml, err = s.repo.SelectMoneyLimits(gCtx, date)
		return err
	})

	g.Go(func() error {
		var err error
		sl, err = s.repo.SelectSecurityLimits(gCtx, date)
		return err
	})

	g.Go(func() error {
		var err error
		otc, err = s.repo.SelectSecurityLimitsOtc(gCtx, date)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	for _, m := range ml {
		res = append(res, qmodels.Limit{
			LimitType:      qmodels.LimitTypeMoney,
			LoadDate:       m.LoadDate,
			SourceDate:     m.SourceDate,
			ClientCode:     m.ClientCode,
			InstrumentCode: m.Currency,
			SettleCode:     m.SettleCode,
			FirmCode:       m.FirmCode,
			FirmName:       m.FirmName,
			Balance:        m.Balance,
		})
	}

	for _, l := range sl {
		res = append(res, qmodels.Limit{
			LimitType:      qmodels.LimitTypeSecurities,
			LoadDate:       l.LoadDate,
			SourceDate:     l.SourceDate,
			ClientCode:     l.ClientCode,
			InstrumentCode: l.Ticker,
			SettleCode:     l.SettleCode,
			FirmCode:       l.FirmCode,
			FirmName:       l.FirmName,
			Balance:        l.Balance,
			ISIN:           l.ISIN,
			AcquisitionCcy: l.AcquisitionCcy,
		})
	}
	for _, o := range otc {
		res = append(res, qmodels.Limit{
			LimitType:      qmodels.LimitTypeSecuritiesOtc,
			LoadDate:       o.LoadDate,
			SourceDate:     o.SourceDate,
			ClientCode:     o.ClientCode,
			InstrumentCode: o.Ticker,
			SettleCode:     o.SettleCode,
			FirmCode:       o.FirmCode,
			FirmName:       o.FirmName,
			Balance:        o.Balance,
			ISIN:           o.ISIN,
			AcquisitionCcy: o.AcquisitionCcy,
		})
	}

	if len(res) == 0 {
		s.logger.Warn("позиции не найдены", zap.Time("load_date", date))
	}
	return res, nil
}

func (s *Service) GetPortfolio(ctx context.Context, targetCcy string) ([]qmodels.PortfolioEntry, error) {
	if targetCcy == "" {
		targetCcy = "RUB"
	}
	targetCcy = strings.ToUpper(strings.TrimSpace(targetCcy))
	if err := validateCurrencyCode(targetCcy); err != nil {
		return nil, err
	}

	var securities []qmodels.PortfolioEntry
	var otcEntries []qmodels.PortfolioEntry
	var moneyEntries []qmodels.PortfolioEntry

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

	entries := make([]qmodels.PortfolioEntry, 0, len(securities)+len(otcEntries)+len(moneyEntries))
	entries = append(entries, securities...)
	entries = append(entries, otcEntries...)
	entries = append(entries, moneyEntries...)

	return entries, nil
}
