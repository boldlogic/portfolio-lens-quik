package repository

import (
	"database/sql"
	"time"

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
