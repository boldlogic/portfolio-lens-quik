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
	ISIN           *string
	ShortName      string
}

type Limit struct {
	LimitType      LimitType
	LoadDate       time.Time
	SourceDate     time.Time
	ClientCode     string
	InstrumentCode string
	ISIN           *string
	SettleCode     SettleCode
	FirmCode       string
	FirmName       string
	Balance        decimal.Decimal
	AcquisitionCcy string
}

type PortfolioEntry struct {
	LimitType      LimitType
	LoadDate       time.Time
	SourceDate     time.Time
	ClientCode     string
	FirmCode       string
	FirmName       string
	Instrument     string
	TradeAccount   string
	PositionCode   string
	ISIN           *string
	AcquisitionCcy string
	ShortName      *string
	QuoteDate      *time.Time
	Balance        decimal.Decimal
	MvCurrency     string
	MvInCcy        decimal.Decimal
	MvPrice        decimal.Decimal
	MvAccrued      decimal.Decimal
	MvTotal        decimal.Decimal
	TargetCurrency string
}
