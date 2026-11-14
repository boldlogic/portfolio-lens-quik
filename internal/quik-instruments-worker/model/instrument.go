package model

import (
	"bytes"
	"time"

	"github.com/shopspring/decimal"
)

type InstrumentKey struct {
	TypeID        uint8
	Ticker        string
	TradePointID  uint8
	HasTradePoint bool
}

type InstrumentDraft struct {
	Key                InstrumentKey
	RegistrationNumber *string
	FullName           *string
	ShortName          *string
	ISIN               *string
	FaceValue          *decimal.Decimal
	MaturityDate       *time.Time
	CouponDuration     *int32
}

type StoredInstrument struct {
	InstrumentID int64
	Draft        InstrumentDraft
	RowVersion   []byte
}

type BoardBindingKey struct {
	InstrumentID int64
	BoardID      int16
}

type BoardBindingDraft struct {
	Key               BoardBindingKey
	CurrencyID        *int16
	BaseCurrencyID    *int16
	QuoteCurrencyID   *int16
	CounterCurrencyID *int16
	IsPrimary         bool
}

type StoredBoardBinding struct {
	Draft      BoardBindingDraft
	RowVersion []byte
}

type CurrencySource struct {
	IsoCharCode  string
	CurrencyName *string
	LatName      *string
}

type SourceQuote struct {
	SourceRowVersion   []byte
	LinkedInstrumentID *int64
	SecCode            string
	TradePointID       uint8
	BoardID            int16
	RegistrationNumber *string
	FullName           *string
	ShortName          *string
	ISIN               *string
	FaceValue          *decimal.Decimal
	MaturityDate       *time.Time
	CouponDuration     *int32
	CurrencyID         *int16
	BaseCurrencyID     *int16
	QuoteCurrencyID    *int16
	CounterCurrencyID  *int16
}

func (d InstrumentDraft) EqualPayload(stored StoredInstrument) bool {
	return equalStringPtr(d.RegistrationNumber, stored.Draft.RegistrationNumber) &&
		equalStringPtr(d.FullName, stored.Draft.FullName) &&
		equalStringPtr(d.ShortName, stored.Draft.ShortName) &&
		equalStringPtr(d.ISIN, stored.Draft.ISIN) &&
		equalDecimalPtr(d.FaceValue, stored.Draft.FaceValue) &&
		equalTimePtr(d.MaturityDate, stored.Draft.MaturityDate) &&
		equalInt32Ptr(d.CouponDuration, stored.Draft.CouponDuration)
}

func (d BoardBindingDraft) EqualPayload(stored StoredBoardBinding) bool {
	return equalInt16Ptr(d.CurrencyID, stored.Draft.CurrencyID) &&
		equalInt16Ptr(d.BaseCurrencyID, stored.Draft.BaseCurrencyID) &&
		equalInt16Ptr(d.QuoteCurrencyID, stored.Draft.QuoteCurrencyID) &&
		equalInt16Ptr(d.CounterCurrencyID, stored.Draft.CounterCurrencyID) &&
		d.IsPrimary == stored.Draft.IsPrimary
}

func RowVersionAfter(left []byte, right []byte) bool {
	return bytes.Compare(left, right) > 0
}

func equalStringPtr(left *string, right *string) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}

func equalInt16Ptr(left *int16, right *int16) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}

func equalInt32Ptr(left *int32, right *int32) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}

func equalDecimalPtr(left *decimal.Decimal, right *decimal.Decimal) bool {
	if left == nil || right == nil {
		return left == right
	}
	return left.Equal(*right)
}

func equalTimePtr(left *time.Time, right *time.Time) bool {
	if left == nil || right == nil {
		return left == right
	}
	return left.Equal(*right)
}
