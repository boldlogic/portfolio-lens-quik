package quik

import (
	"time"

	"github.com/shopspring/decimal"
)

type LimitType string

const (
	LimitTypeSecurities    LimitType = "securities"     // ценные бумаги (биржевые)
	LimitTypeSecuritiesOtc LimitType = "securities_otc" // ценные бумаги OTC
	LimitTypeMoney         LimitType = "money"          // денежные средства
)

const (
	MaxClientCodeLen = 12
	MinClientCodeLen = 1
)

type Limit struct {
	LoadDate     time.Time
	SourceDate   time.Time
	ClientCode   string
	Currency     string
	PositionCode string
	SettleCode   SettleCode
	FirmCode     string
	FirmName     string
	Balance      decimal.Decimal
}

type MoneyLimit struct {
	LoadDate     time.Time
	SourceDate   time.Time
	ClientCode   string
	Currency     string
	PositionCode string
	SettleCode   SettleCode
	FirmCode     string
	FirmName     string
	Balance      decimal.Decimal
}

type SecurityLimit struct {
	LoadDate       time.Time
	SourceDate     time.Time
	ClientCode     string
	Ticker         string
	TradeAccount   string
	SettleCode     SettleCode
	FirmCode       string
	FirmName       string
	Balance        decimal.Decimal
	AcquisitionCcy string
	ISIN           string
	ShortName      string
}

type Position struct {
	LimitType    LimitType
	LoadDate     time.Time       //дата исходного лимита
	SourceDate   time.Time       // если не равен LoadDate то означает дату, с которой перенесен лимит
	ClientCode   string          //код клиента
	FirmCode     string          //код фирмы
	FirmName     string          //название фирмы
	Ticker       string          //код инструмента
	Name         string          //название инструмента
	Balance      decimal.Decimal //текущее количество на максимальную дату расчета
	Price        decimal.Decimal //цена позиции
	AccruedInt   decimal.Decimal //НКД в валюте инструмента
	MVInstr      decimal.Decimal //оценка позиции в валюте инструмента
	MVTotal      decimal.Decimal // оценка позиции в валюте запроса
	CurrencyCode string          //валюта инструмента

}
