package mssql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-currency/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/dbrepo"
)

const (
	selectNewCurrenciesFromCrossrates = `
		WITH crossrate AS (
			SELECT
				UPPER(COALESCE(NULLIF(q.currency, ''), NULLIF(q.sec_code, ''))) AS iso_char_code,
				COALESCE(q.full_name, q.short_name) AS currency_name
			FROM quik.current_quotes q
			WHERE q.class_code = 'CROSSRATE'
		),
		cets_codes AS (
			SELECT DISTINCT v.code AS iso_char_code
			FROM quik.current_quotes q
			CROSS APPLY (VALUES (q.base_currency), (q.quote_currency)) AS v(code)
			WHERE q.class_code = 'CETS'
		),
		src AS (
			SELECT iso_char_code, currency_name
			FROM crossrate

			UNION

			SELECT c.iso_char_code, ''
			FROM cets_codes c
			WHERE NOT EXISTS (
				SELECT 1
				FROM crossrate cr
				WHERE cr.iso_char_code = c.iso_char_code
			)
		)
		SELECT s.iso_char_code, s.currency_name
		FROM src s
		WHERE s.iso_char_code <> ''
		AND NOT EXISTS (
			SELECT 1
			FROM ref.currencies r
			WHERE r.iso_char_code = s.iso_char_code
		)
		ORDER BY case when s.currency_name='' then 1 else 0 end;`
)

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
