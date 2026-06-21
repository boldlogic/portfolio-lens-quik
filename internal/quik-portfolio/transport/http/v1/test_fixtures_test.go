package v1

import (
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	"github.com/shopspring/decimal"
)

const (
	fixtureClientCode     = "AB12CD"
	fixtureFirmCode       = "NC0058900000"
	fixtureTradeAccount   = "L01-00000F00"
	fixtureISIN           = "RU000A0JX0J2"
	fixturePositionCode   = "EQTV"
	fixtureSecCodeSBER    = "SBER"
	fixtureSecCodeGAZP    = "GAZP"
	fixtureFirmName       = "Фирма брокера"
	fixtureLoadDate       = "2026-05-31"
	fixtureSourceDate     = "2026-05-30"
)

var (
	fixtureLoadDateTime   = time.Date(2026, 5, 31, 0, 0, 0, 0, time.Local)
	fixtureSourceDateTime = time.Date(2026, 5, 30, 0, 0, 0, 0, time.Local)
)

func fixtureMoneyLimit() quik.MoneyLimit {
	return quik.MoneyLimit{
		LoadDate:     fixtureLoadDateTime,
		SourceDate:   fixtureSourceDateTime,
		ClientCode:   fixtureClientCode,
		CurrencyCode: "RUB",
		PositionCode: fixturePositionCode,
		SettleCode:   quik.SettleCodeT2,
		FirmCode:     fixtureFirmCode,
		FirmName:     fixtureFirmName,
		Balance:      decimal.RequireFromString("10.25"),
	}
}

func fixtureMoneyLimitDTO() moneyLimitDTO {
	return moneyLimitDTO{
		LoadDate:     fixtureLoadDate,
		SourceDate:   fixtureSourceDate,
		ClientCode:   fixtureClientCode,
		Currency:     "RUB",
		PositionCode: fixturePositionCode,
		SettleCode:   "T2",
		FirmCode:     fixtureFirmCode,
		FirmName:     fixtureFirmName,
		Balance:      decimal.RequireFromString("10.25"),
	}
}

func fixtureSecurityLimit() quik.SecurityLimit {
	return quik.SecurityLimit{
		LoadDate:                fixtureLoadDateTime,
		SourceDate:              fixtureSourceDateTime,
		ClientCode:              fixtureClientCode,
		SecCode:                 fixtureSecCodeSBER,
		TradeAccount:            fixtureTradeAccount,
		SettleCode:              quik.SettleCodeT2,
		FirmCode:                fixtureFirmCode,
		FirmName:                fixtureFirmName,
		Balance:                 decimal.RequireFromString("10.25"),
		AcquisitionCurrencyCode: "RUB",
		ISIN:                    fixtureISIN,
		ShortName:               "Сбербанк",
	}
}

func fixtureSecurityLimitDTO() securityLimitDTO {
	return securityLimitDTO{
		LoadDate:                fixtureLoadDate,
		SourceDate:              fixtureSourceDate,
		ClientCode:              fixtureClientCode,
		SecCode:                 fixtureSecCodeSBER,
		TradeAccount:            fixtureTradeAccount,
		SettleCode:              "T2",
		FirmCode:                fixtureFirmCode,
		FirmName:                fixtureFirmName,
		Balance:                 decimal.RequireFromString("10.25"),
		AcquisitionCurrencyCode: "RUB",
		ISIN:                    fixtureISIN,
		ShortName:               "Сбербанк",
	}
}
