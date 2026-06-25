package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type InstrumentWithBoard struct {
	InstrumentClass        string
	SecCode                string
	TradePointId           uint8
	ISIN                   *string
	RegistrationNumber     *string
	FullName               *string
	ShortName              *string
	FaceValue              *decimal.Decimal
	MaturityDate           *time.Time
	CouponDuration         *int32
	ClassCode              string
	BoardId                uint8
	Currency               *string
	CurrencyNumeric        *int16
	BaseCurrency           *string
	BaseCurrencyNumeric    *int16
	QuoteCurrency          *string
	QuoteCurrencyNumeric   *int16
	CounterCurrency        *string
	CounterCurrencyNumeric *int16
	InstrumentId           *int64
}
