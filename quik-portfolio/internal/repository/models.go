package repository

import (
	"database/sql"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)




func toPosition(positions []moneyPosition) []quik.Position {
	var out = make([]quik.Position, 0, len(positions))

	for _, p := range positions {
		var currencyName string
		if p.CurrencyName.Valid {
			currencyName = p.CurrencyName.String
		}
		out = append(out, quik.Position{
			LimitType:  quik.LimitTypeMoney,
			LoadDate:   p.LoadDate,
			ClientCode: p.ClientCode,
			FirmCode:   p.FirmCode,
			FirmName:   p.FirmName,
			Ticker:     p.CurrencyCode,
			Name:       currencyName,
			Balance:    p.Balance,
			MVInstr:    p.Balance,
			MVTotal:    p.MV,
		})
	}

	return out
}

func (p securityPosition) convertToPosition() quik.Position {
	out := quik.Position{
		LimitType:  quik.LimitTypeSecurities,
		LoadDate:   p.LoadDate,
		SourceDate: p.SourceDate,
		ClientCode: p.ClientCode,
		FirmCode:   p.FirmCode,
		Ticker:     p.SecCode,
	}

	out.Name = p.SecName.String

	out.Balance = p.Balance

	if p.Price != nil {
		out.Price = *p.Price
	}

	if p.AccruedInt != nil {
		out.AccruedInt = *p.AccruedInt
	}
	if p.MVInstr != nil {
		out.MVInstr = *p.MVInstr
	}

	if p.MV != nil {
		out.MVTotal = *p.MV
	}

	return out
}

func rawToPosition(positions []securityPosition) []quik.Position {
	var out = make([]quik.Position, 0, len(positions))

	for _, p := range positions {
		out = append(out, p.convertToPosition())
	}

	return out
}


type securityPosition struct {
	LoadDate     time.Time
	SourceDate   time.Time
	QuoteDate    sql.NullTime
	RateDate     sql.NullTime
	ClientCode   string
	FirmCode     string
	FirmName     sql.NullString
	SecCode      string
	SecName      sql.NullString
	Balance      decimal.Decimal
	Price        *decimal.Decimal
	AccruedInt   *decimal.Decimal
	MVInstr      *decimal.Decimal
	MV           *decimal.Decimal
	CurrencyCode string
}
