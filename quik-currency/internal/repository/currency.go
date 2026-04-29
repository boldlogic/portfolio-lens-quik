package repository

import (
	"context"

	"github.com/boldlogic/portfolio-lens-currency/pkg/currencies"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"go.uber.org/zap"
)

const (
	mergeCurrencies = `
WITH src AS (
	SELECT
		iso_code = @p1,
		iso_char_code = @p2,
		currency_name = @p3,
		lat_name = @p4,
		minor_units = @p5

)
MERGE INTO quik.currencies AS tgt
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
	countCurrencies = `
		SELECT COUNT(1) FROM quik.currencies`
)

func (r *Repository) SelectCountCurrencies(ctx context.Context) (int, error) {
	var res int

	row := r.Db.QueryRowContext(ctx, countCurrencies)
	err := row.Scan(&res)

	if err != nil {
		if r.isShutdown(err) {
			return 0, err
		}
		r.Logger.Error("ошибка при получении количества валют в справочнике currencies", zap.Error(err))

		return 0, err
	}
	return res, nil
}
func (r *Repository) MergeCurrencies(ctx context.Context, currencies []currencies.Currency) error {

	if len(currencies) == 0 {
		return nil
	}

	for _, ccy := range currencies {
		_, err := r.Db.ExecContext(ctx, mergeCurrencies,
			ccy.ISOCode,
			ccy.ISOCharCode,
			ccy.Name,
			ccy.LatName,
			ccy.MinorUnits,
		)
		if err != nil {
			if r.isShutdown(err) {
				return err
			}
			r.Logger.Error("ошибка сохранения валюты", zap.Int16("iso_code", ccy.ISOCode), zap.Error(err))
			return models.ErrSavingData
		}
	}

	return nil

}
