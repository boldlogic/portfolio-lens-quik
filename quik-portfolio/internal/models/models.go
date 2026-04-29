package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type CurrentQuote struct {
	Ticker          string
	ISIN            *string
	LastPrice       *decimal.Decimal
	ClosePrice      *decimal.Decimal
	AccruedInt      *decimal.Decimal
	FaceValue       *decimal.Decimal
	InstrumentType  string
	PriceCurrency   string
	AccruedCurrency string
	QuoteDate       time.Time
}
