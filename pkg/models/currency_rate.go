package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type FxRate struct {
	Date             time.Time
	QuoteISOCode     int
	BaseISOCode      int
	RateQuotePerBase decimal.Decimal
	RateBasePerQuote decimal.Decimal
	CreatedAt        time.Time
	UpdatedAt        time.Time
	BaseCurrency     *Currency
	QuoteCurrency    *Currency
}
