package models

import "github.com/shopspring/decimal"

type Limit struct {
	Type                    string
	ClientCode              string
	Ticker                  string
	PositionCode            *string
	SettleCode              string
	TradeAccount            *string
	FirmCode                string
	Balance                 decimal.Decimal
	AcquisitionCurrencyCode *string
	ISIN                    *string
}

type LimitLine struct {
	Limit
	Line uint
}
