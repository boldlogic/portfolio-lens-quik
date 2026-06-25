package quik

import (
	"time"

	"github.com/shopspring/decimal"
)

// type InstrumentTypeId uint8
type InstrumentTypeTitle string

const (
	InstrumentTitleSecurities = "Акции"
	InstrumentTitleBonds      = "Облигации"
	InstrumentTitleCurrencies = "Валюта"
)

type InstrumentType struct {
	Id    uint8
	Title InstrumentTypeTitle
}

type Board struct {
	Id           uint8
	TradePointId *uint8
	Code         string
	Name         string
	TradePoint   *TradePoint
	IsTraded     bool
}

type TradePoint struct {
	Id   uint8
	Code string
	Name string
}

type Instrument struct {
	InstrumentId   int64
	SecCode        string
	TradePointId   uint8
	ISIN           *string
	FullName       *string
	ShortName      *string
	RegNumber      *string
	FaceValue      *decimal.Decimal
	MaturityDate   *time.Time
	CouponDuration *int32
	Boards         []Board
}
