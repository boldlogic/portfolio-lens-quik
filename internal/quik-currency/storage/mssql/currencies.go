package mssql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-currency/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/dbrepo"
)

const (
	selectExistingCurrencies = `
		SELECT iso_code
		,iso_char_code
		,currency_name
		,lat_name
		,minor_units
	FROM ref.currencies
  `
	selectNewCurrenciesFromCrossrates = `
		with
			b as (
				SELECT
					b.board_id,
					b.code,
					b.trade_point_id,
					t.title
				FROM
					ref.boards b
					JOIN ref.instrument_type_boards tb ON tb.board_id=b.board_id
					JOIN ref.instrument_types t ON t.type_id=tb.type_id
				WHERE
					t.title='Валюта'
			),
			crossrate AS (
				SELECT
					q.class_code,
					iso_char_code=UPPER(
						COALESCE(NULLIF(q.currency, ''), NULLIF(q.sec_code, ''))
					),
					currency_name=COALESCE(q.full_name, q.short_name),
					q.rw
				FROM
					b
					join quik.current_quotes q on q.class_code=b.code
			)
		SELECT
			s.iso_char_code,
			s.currency_name
		FROM
			crossrate s
		WHERE
			s.iso_char_code<>''
		ORDER BY
			case
				when s.currency_name='' then 1
				else 0
			end;
`
)

func scanCurrency(row *sql.Rows) (models.Currency, error) {

	var res models.Currency
	err := row.Scan(&res.Numeric, &res.CharCode, &res.Name, &res.LatName, &res.MinorUnits)
	if err != nil {
		return models.Currency{}, fmt.Errorf("ошибка чтения: %w", err)
	}
	return res, nil
}

func scanCurrencyFromCrossrates(row *sql.Rows) (models.CurrencyFromCrossrates, error) {
	var res models.CurrencyFromCrossrates
	err := row.Scan(&res.IsoCharCode, &res.Name)
	if err != nil {
		return models.CurrencyFromCrossrates{}, fmt.Errorf("ошибка чтения: %w", err)
	}
	return res, nil
}

func (repo CurrencyRepo) SelectNewCurrenciesFromCrossrates(ctx context.Context) ([]models.CurrencyFromCrossrates, error) {
	return dbrepo.SelectRows(ctx, repo.runner, selectNewCurrenciesFromCrossrates, scanCurrencyFromCrossrates)
}

func (repo CurrencyRepo) SelectCurrencies(ctx context.Context) ([]models.Currency, error) {
	return dbrepo.SelectRows(ctx, repo.runner, selectExistingCurrencies, scanCurrency)
}
