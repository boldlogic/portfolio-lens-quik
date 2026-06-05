package repository

import (
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func moneyPortfolioRowsToPositions(positions []moneyPortfolioRow) []quik.Position {
	var out = make([]quik.Position, 0, len(positions))

	for _, p := range positions {
		var currencyName string
		if p.CurrencyName.Valid {
			currencyName = p.CurrencyName.String
		}
		out = append(out, quik.Position{
			LimitType:  quik.LimitTypeMoney,
			LoadDate:   p.LoadDate,
			ClientCode: p.ClientCode,
			FirmCode:   p.FirmCode,
			FirmName:   p.FirmName,
			Ticker:     p.CurrencyCode,
			Name:       currencyName,
			Balance:    p.Balance,
			MVInstr:    p.Balance,
			MVTotal:    p.MarketValueInTargetCurrency,
		})
	}

	return out
}

func (p securityPortfolioRow) toQuikPosition() quik.Position {
	out := quik.Position{
		LimitType:  quik.LimitTypeSecurities,
		LoadDate:   p.LoadDate,
		SourceDate: p.SourceDate,
		ClientCode: p.ClientCode,
		FirmCode:   p.FirmCode,
		Ticker:     p.SecCode,
	}

	out.Name = p.SecName.String

	out.Balance = p.Balance

	if p.UnitPrice != nil {
		out.Price = *p.UnitPrice
	}

	if p.AccruedInterest != nil {
		out.AccruedInt = *p.AccruedInterest
	}
	if p.MarketValueInInstrCurrency != nil {
		out.MVInstr = *p.MarketValueInInstrCurrency
	}

	if p.MarketValueInTargetCurrency != nil {
		out.MVTotal = *p.MarketValueInTargetCurrency
	}

	return out
}

func securityPortfolioRowsToPositions(positions []securityPortfolioRow) []quik.Position {
	var out = make([]quik.Position, 0, len(positions))

	for _, p := range positions {
		out = append(out, p.toQuikPosition())
	}

	return out
}
