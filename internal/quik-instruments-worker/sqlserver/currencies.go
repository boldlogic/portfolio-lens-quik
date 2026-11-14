package sqlserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/model"
)

const selectCurrenciesSQL = `
SELECT
	iso_char_code,
	currency_name,
	lat_name
FROM ref.currencies
ORDER BY iso_char_code;`

func (r *Repository) SelectCurrencies(ctx context.Context) ([]model.CurrencySource, error) {
	rows, err := r.Db.QueryContext(ctx, selectCurrenciesSQL)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения валют: %w", err)
	}
	defer rows.Close()

	var out []model.CurrencySource
	for rows.Next() {
		item, err := scanCurrency(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка обхода валют: %w", err)
	}
	return out, nil
}

func scanCurrency(row *sql.Rows) (model.CurrencySource, error) {
	var isoCharCode string
	var currencyName sql.NullString
	var latName sql.NullString

	err := row.Scan(&isoCharCode, &currencyName, &latName)
	if err != nil {
		return model.CurrencySource{}, fmt.Errorf("ошибка чтения валюты: %w", err)
	}

	return model.CurrencySource{
		IsoCharCode:  isoCharCode,
		CurrencyName: stringPtr(currencyName),
		LatName:      stringPtr(latName),
	}, nil
}
