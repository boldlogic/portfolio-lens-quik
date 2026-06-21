package repository

import (
	"github.com/boldlogic/portfolio-lens-quik/pkg/dbrepo"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func mapRows[R, D any](rows []R, fn func(R) D) []D {
	out := make([]D, 0, len(rows))
	for _, row := range rows {
		out = append(out, fn(row))
	}
	return out
}

func (p moneyPortfolioRow) toQuikPosition() quik.Position {
	return quik.Position{
		LimitType:                   quik.LimitTypeMoney,
		LoadDate:                    p.LoadDate,
		ClientCode:                  p.ClientCode,
		FirmCode:                    p.FirmCode,
		FirmName:                    p.FirmName,
		Ticker:                      p.CurrencyCode,
		Name:                        dbrepo.StringFromNull(p.CurrencyName),
		Amount:                      p.Balance,
		MarketValueInInstrCurrency:  p.Balance,
		MarketValueInTargetCurrency: p.MarketValueInTargetCurrency,
		InstrumentCurrencyCode:      p.CurrencyCode,
	}
}

func (p securityPortfolioRow) toQuikPositionWithType(limitType quik.LimitType) quik.Position {
	return quik.Position{
		LimitType:                   limitType,
		LoadDate:                    p.LoadDate,
		SourceDate:                  p.SourceDate,
		ClientCode:                  p.ClientCode,
		FirmCode:                    p.FirmCode,
		FirmName:                    dbrepo.StringFromNull(p.FirmName),
		Ticker:                      p.SecCode,
		Name:                        dbrepo.StringFromNull(p.SecName),
		Amount:                      p.Balance,
		UnitPrice:                   dbrepo.DecimalFromPtr(p.UnitPrice),
		AccruedInterest:             dbrepo.DecimalFromPtr(p.AccruedInterest),
		MarketValueInInstrCurrency:  dbrepo.DecimalFromPtr(p.MarketValueInInstrCurrency),
		MarketValueInTargetCurrency: dbrepo.DecimalFromPtr(p.MarketValueInTargetCurrency),
		InstrumentCurrencyCode:      p.InstrumentCurrencyCode,
	}
}
