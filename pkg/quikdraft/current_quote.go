package quikdraft

import (
	"fmt"
	"strings"
	"time"
)

type CurrentQuote struct {
	QuoteDate          *time.Time
	InstrumentClass    string
	Ticker             string
	ISIN               *string
	RegistrationNumber *string
	FullName           *string
	ShortName          string
	FaceValue          *float64
	MaturityDate       *time.Time
	CouponDuration     *int
	ClassCode          string
	ClassName          string
	InstrumentType     string
	InstrumentSubtype  *string
	Currency           string
	BaseCurrency       string
	QuoteCurrency      *string
	CounterCurrency    *string
	InstrumentId       int

	LastPrice     *float64
	ClosePrice    *float64
	AccruedInt      *float64
	TradingStatus *string
}

func (q *CurrentQuote) Clear() {
	q.InstrumentClass = strings.TrimSpace(q.InstrumentClass)
	q.Ticker = strings.TrimSpace(q.Ticker)
	if q.RegistrationNumber != nil {
		trimmedRn := strings.TrimSpace(*q.RegistrationNumber)
		q.RegistrationNumber = &trimmedRn
	}
	if q.FullName != nil {
		trimmedFn := strings.TrimSpace(*q.FullName)
		q.FullName = &trimmedFn
	}
	q.ShortName = strings.TrimSpace(q.ShortName)
	q.ClassCode = strings.TrimSpace(q.ClassCode)
	q.ClassName = strings.TrimSpace(q.ClassName)
	q.InstrumentType = strings.TrimSpace(q.InstrumentType)
	if q.InstrumentSubtype != nil {
		trimmedSt := strings.TrimSpace(*q.InstrumentSubtype)
		q.InstrumentSubtype = &trimmedSt
	}
	if q.ISIN != nil {
		trimmedIsin := strings.TrimSpace(*q.ISIN)
		q.ISIN = &trimmedIsin
	}
	if q.QuoteCurrency != nil {
		trimmedQc := strings.TrimSpace(*q.QuoteCurrency)
		q.QuoteCurrency = &trimmedQc
	}
	if q.CounterCurrency != nil {
		trimmedCc := strings.TrimSpace(*q.CounterCurrency)
		q.CounterCurrency = &trimmedCc
	}
}

func (q CurrentQuote) String() string {
	faceVal := "nil"
	if q.FaceValue != nil {
		faceVal = fmt.Sprintf("%g", *q.FaceValue)
	}
	matDate := "nil"
	if q.MaturityDate != nil {
		matDate = q.MaturityDate.Format(time.DateOnly)
	}
	isin := "nil"
	if q.ISIN != nil {
		isin = fmt.Sprintf("%q", *q.ISIN)
	}
	return fmt.Sprintf("CurrentQuote{Ticker:%q ShortName:%q ClassCode:%q ISIN:%s FaceValue:%s BaseCurrency:%q MaturityDate:%s}",
		q.Ticker, q.ShortName, q.ClassCode, isin, faceVal, q.BaseCurrency, matDate)
}
