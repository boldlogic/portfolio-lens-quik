package mssql

import (
	"context"
	"database/sql"

	"github.com/boldlogic/portfolio-lens-currency/pkg/currencies"
)

const (
	mergeCurrencies = `
WITH src AS (
	SELECT
		iso_code,
		iso_char_code,
		currency_name,
		lat_name,
		minor_units
	FROM @currencies

)
MERGE INTO ref.currencies AS tgt
USING src ON tgt.iso_code = src.iso_code
WHEN MATCHED
	AND (
		tgt.currency_name <> src.currency_name
		OR tgt.lat_name <> src.lat_name
		OR tgt.minor_units <> src.minor_units
	)
THEN UPDATE SET
	tgt.currency_name = src.currency_name,
	tgt.lat_name = src.lat_name,
	tgt.minor_units = src.minor_units,
	tgt.updated_at = SYSDATETIMEOFFSET()
WHEN NOT MATCHED BY TARGET
THEN INSERT (
	iso_code,
	iso_char_code,
	currency_name,
	lat_name,
	minor_units,
	updated_at
)
VALUES (
	src.iso_code,
	src.iso_char_code,
	src.currency_name,
	src.lat_name,
	src.minor_units,
	SYSDATETIMEOFFSET()
);`
)

func (repo CurrencyRepo) MergeCurrencies(ctx context.Context, currencies []currencies.Currency) error {
	crrs, ok := makeCurrencyList(currencies)
	if !ok {
		return nil
	}
	_, err := repo.runner.ExecContext(ctx, mergeCurrencies, sql.Named("currencies", crrs))
	if err != nil {
		return err
	}

	return nil
}
