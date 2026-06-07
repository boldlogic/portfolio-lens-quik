package repository

import (
	"context"

	"github.com/boldlogic/packages/shutdown"
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
MERGE INTO dbo.currencies AS tgt
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
		SELECT COUNT(1) FROM dbo.currencies`
	setEmptyNamesFromQuik = `
		WITH NAMES AS (
			SELECT
				iso_char_code = RTRIM(COALESCE(c.iso_char_code, norm.sec_code)),
				currency_name = MAX(RTRIM(COALESCE(q.full_name, q.short_name))),
				ext_system_id=max(es.ext_system_id)
			FROM (
				SELECT DISTINCT
					sec_code = CASE WHEN q.sec_code IN ('SUR', 'RUR', 'RUB') THEN 'RUB' ELSE q.sec_code END,
					full_name = q.full_name,
					short_name = q.short_name
				FROM quik.current_quotes q
				WHERE q.class_code = 'CROSSRATE'
					AND LEN(q.sec_code) = 3
			) q
			CROSS APPLY (SELECT sec_code = q.sec_code) norm
			LEFT JOIN dbo.external_systems es on es.ext_system='QUIK'
			LEFT JOIN dbo.external_codes ec ON ec.ext_code = norm.sec_code AND ec.ext_system_id = es.ext_system_id AND ec.ext_code_type_id = 1 
			LEFT JOIN dbo.currencies c ON c.iso_code = ec.internal_id
			GROUP BY COALESCE(c.iso_char_code, norm.sec_code)
		)
		UPDATE c
		SET
			c.currency_name = n.currency_name,
			updated_at = GETDATE(),
			ext_system_id = n.ext_system_id
		FROM dbo.currencies c
		INNER JOIN NAMES n ON c.iso_char_code = n.iso_char_code
		WHERE c.currency_name IS NULL;`
)

func (r *Repository) SelectCountCurrencies(ctx context.Context) (int, error) {
	var res int

	row := r.Db.QueryRowContext(ctx, countCurrencies)
	err := row.Scan(&res)

	if err != nil {
		if shutdown.IsExceeded(err) {
			return 0, err
		}
		r.Logger.Error("ошибка при получении количества валют в справочнике currencies", zap.Error(err))

		return 0, err
	}
	return res, nil
}

func (r *Repository) SetEmptyCurrencyNamesFromQuik(ctx context.Context) error {
	_, err := r.Db.ExecContext(ctx, setEmptyNamesFromQuik)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.Logger.Error("ошибка при обновлении currency_name в currencies", zap.Error(err))
		return models.ErrSavingData
	}

	return nil
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
			if shutdown.IsExceeded(err) {
				return err
			}
			r.Logger.Error("ошибка сохранения валюты", zap.Int16("iso_code", ccy.ISOCode), zap.Error(err))
			return models.ErrSavingData
		}
	}

	return nil

}
