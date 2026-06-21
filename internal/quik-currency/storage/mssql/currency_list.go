package mssql

import (
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	sqlserver "github.com/microsoft/go-mssqldb"
)

type currencyTVP struct {
	ISOCode     int16
	ISOCharCode string
	Name        *string
	LatName     string
	MinorUnits  int32
}

func currencyToTVP(c quik.Currency) currencyTVP {
	return currencyTVP{
		ISOCode:     c.Numeric(),
		ISOCharCode: c.Alpha().String(),
		Name:        c.Name(),
		LatName:     c.LatName(),
		MinorUnits:  c.MinorUnits(),
	}
}

func makeCurrencyList(curr []quik.Currency) (sqlserver.TVP, bool) {
	if len(curr) == 0 {
		return sqlserver.TVP{}, false
	}

	currencies := make([]currencyTVP, 0, len(curr))
	for _, c := range curr {
		currencies = append(currencies, currencyToTVP(c))
	}

	return sqlserver.TVP{
		TypeName: "app.currency_list",
		Value:    currencies,
	}, true
}
