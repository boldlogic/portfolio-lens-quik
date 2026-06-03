package repository

import (
	"database/sql"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

type portfolioEntry struct {
	LoadDate   time.Time
	SourceDate time.Time
	ClientCode string

	Instrument   string
	TradeAccount string
	FirmCode     string
	FirmName     string

	PositionCode string
	ISIN         *string
	MvCurrency   sql.NullString

	MinorUnits sql.NullInt32

	AcquisitionCcy string
	ShortName      sql.NullString
	QuoteDate      sql.NullTime
	Balance        decimal.Decimal

	MvInCcy          decimal.Decimal
	MvPrice          decimal.Decimal
	MvAccrued        decimal.Decimal
	MvTotal          decimal.Decimal
	TargetCurrency   sql.NullString
	TargetMinorUnits sql.NullInt32
}

func (r *Repository) scanSecurityToPortfolio(row *sql.Rows) (portfolioEntry, error) {
	out := portfolioEntry{}
	err := row.Scan(
		&out.LoadDate,
		&out.SourceDate,
		&out.ClientCode,
		&out.Instrument,
		&out.TradeAccount,
		&out.FirmCode,
		&out.FirmName,
		&out.Balance,
		&out.AcquisitionCcy,
		&out.ISIN,
		&out.MvCurrency,
		&out.MinorUnits,
		&out.MvInCcy,
		&out.MvPrice,
		&out.MvAccrued,
		&out.MvTotal,
		&out.TargetCurrency,
		&out.TargetMinorUnits,
		&out.ShortName,
		&out.QuoteDate,
	)
	if err != nil {
		return portfolioEntry{}, err
	}
	return out, nil
}

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

type moneyPosition struct {
	LoadDate     time.Time
	SourceDate   time.Time
	ClientCode   string
	FirmCode     string
	FirmName     string
	CurrencyCode string
	CurrencyName sql.NullString
	Balance      decimal.Decimal
	MV           decimal.Decimal

	// TradeAccount string

	// PositionCode string
	// ISIN         *string
	// MvCurrency   sql.NullString

	// MinorUnits sql.NullInt32

	// AcquisitionCcy string

	// QuoteDate      sql.NullTime

	// MvInCcy          decimal.Decimal
	// MvPrice          decimal.Decimal
	// MvAccrued        decimal.Decimal
	// MvTotal          decimal.Decimal
	// TargetCurrency   sql.NullString
	// TargetMinorUnits sql.NullInt32
}

type position struct {
	LoadDate   time.Time
	SourceDate time.Time
	ClientCode string
	FirmCode   string
	FirmName   string
	Instrument string
	ShortName  sql.NullString
	Balance    decimal.Decimal
	MvTotal    decimal.Decimal

	// TradeAccount string

	// PositionCode string
	// ISIN         *string
	// MvCurrency   sql.NullString

	// MinorUnits sql.NullInt32

	// AcquisitionCcy string

	// QuoteDate      sql.NullTime

	// MvInCcy          decimal.Decimal
	// MvPrice          decimal.Decimal
	// MvAccrued        decimal.Decimal
	// MvTotal          decimal.Decimal
	// TargetCurrency   sql.NullString
	// TargetMinorUnits sql.NullInt32
}
