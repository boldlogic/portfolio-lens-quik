package repository

import (
	"database/sql"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

func mapRows[R, D any](rows []R, fn func(R) D) []D {
	out := make([]D, 0, len(rows))
	for _, row := range rows {
		out = append(out, fn(row))
	}
	return out
}

func stringFromNull(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func decimalFromPtr(p *decimal.Decimal) decimal.Decimal {
	if p != nil {
		return *p
	}
	return decimal.Decimal{}
}

func (m moneyLimitRow) toQuik() quik.MoneyLimit {
	ml := quik.MoneyLimit{
		LoadDate:     m.LoadDate,
		SourceDate:   m.SourceDate,
		ClientCode:   m.ClientCode,
		Currency:     m.CurrencyCode,
		PositionCode: m.PositionCode,
		FirmCode:     m.FirmCode,
		FirmName:     stringFromNull(m.FirmName),
		Balance:      decimalFromPtr(m.Balance),
		SettleCode:   quik.SettleCode(m.SettleCode),
	}
	return ml
}

func (s securityLimitRow) toQuik() quik.SecurityLimit {
	return quik.SecurityLimit{
		LoadDate:       s.LoadDate,
		SourceDate:     s.SourceDate,
		ClientCode:     s.ClientCode,
		Ticker:         s.SecCode,
		TradeAccount:   s.TradeAccount,
		FirmCode:       s.FirmCode,
		FirmName:       stringFromNull(s.FirmName),
		Balance:        decimalFromPtr(s.Balance),
		AcquisitionCcy: s.AcquisitionCurrencyCode,
		ISIN:           stringFromNull(s.ISIN),
		ShortName:      stringFromNull(s.ShortName),
		SettleCode:     quik.SettleCode(s.SettleCode),
	}
}

func (p moneyPortfolioRow) toQuikPosition() quik.Position {
	return quik.Position{
		LimitType:                   quik.LimitTypeMoney,
		LoadDate:                    p.LoadDate,
		ClientCode:                  p.ClientCode,
		FirmCode:                    p.FirmCode,
		FirmName:                    p.FirmName,
		Ticker:                      p.CurrencyCode,
		Name:                        stringFromNull(p.CurrencyName),
		Amount:                      p.Balance,
		MarketValueInInstrCurrency:  p.Balance,
		MarketValueInTargetCurrency: p.MarketValueInTargetCurrency,
	}
}

func (p securityPortfolioRow) toQuikPosition() quik.Position {
	return quik.Position{
		LimitType:                   quik.LimitTypeSecurities,
		LoadDate:                    p.LoadDate,
		SourceDate:                  p.SourceDate,
		ClientCode:                  p.ClientCode,
		FirmCode:                    p.FirmCode,
		Ticker:                      p.SecCode,
		Name:                        stringFromNull(p.SecName),
		Amount:                      p.Balance,
		UnitPrice:                   decimalFromPtr(p.UnitPrice),
		AccruedInterest:             decimalFromPtr(p.AccruedInterest),
		MarketValueInInstrCurrency:  decimalFromPtr(p.MarketValueInInstrCurrency),
		MarketValueInTargetCurrency: decimalFromPtr(p.MarketValueInTargetCurrency),
		InstrumentCurrencyCode:      p.InstrumentCurrencyCode,
	}
}
