package mssql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-currency/models"
)

const (
	selectNewCurrenciesFromCrossrates = `
		with src as(
			SELECT
				norm.iso_char_code,
				n.currency_name,
				rn = ROW_NUMBER() over(
					partition by norm.iso_char_code
					order by
						n.currency_name
				)
			FROM
				quik.current_quotes q
				CROSS APPLY (
					SELECT
						iso_char_code = UPPER(
							COALESCE(NULLIF(q.currency, ''), NULLIF(q.sec_code, ''))
						)
				) raw
				CROSS APPLY (
					SELECT
						iso_char_code = 
							CASE 
								WHEN raw.iso_char_code IN ('SUR', 'RUR') THEN 'RUB' 
								WHEN raw.iso_char_code = 'USDX' THEN 'USD' 
								WHEN raw.iso_char_code = 'GLD' THEN 'XAU' 
								WHEN raw.iso_char_code = 'SLV' THEN 'XAG' 
								WHEN raw.iso_char_code = 'PLT' THEN 'XPT' 
								WHEN raw.iso_char_code = 'PLD' THEN 'XPD' 
								ELSE raw.iso_char_code 
							END
				) norm
				CROSS APPLY (
					select
						currency_name = (COALESCE(q.full_name, q.short_name))
				) n
			WHERE
				q.class_code = 'CROSSRATE'
		)
		select
			c.iso_char_code,
			c.currency_name
		from
			src c
		where
			c.rn = 1
			and c.iso_char_code<>''
			and not exists (
				select
					1
				from
					ref.currencies r
				where
					r.iso_char_code = c.iso_char_code
			)`
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
	rows, err := repo.runner.QueryContext(ctx, selectNewCurrenciesFromCrossrates)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]models.CurrencyFromCrossrates, 0, 50)
	for rows.Next() {
		row, err := scanCurrencyFromCrossrates(rows)
		if err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
