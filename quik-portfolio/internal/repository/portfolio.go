package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/boldlogic/packages/shutdown"
	qmodels "github.com/boldlogic/portfolio-lens-quik/quik-portfolio/internal/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

const selectSecuritiesPortfolio = `
WITH
    cte AS (
        SELECT
            li.load_date,
            li.source_date,
            li.client_code,
            li.ticker,
            li.trade_account,
            li.firm_code,
            li.settle_code,
            li.firm_name,
            li.balance,
            li.acquisition_ccy,
            li.isin,
            settle_max = MAX(li.settle_code) OVER (
                PARTITION BY li.load_date, li.client_code, li.ticker, li.trade_account, li.firm_code
            )
        FROM quik.security_limits li
        WHERE li.load_date = cast(@p1 as date)
    ),
    cte_filtered AS (
        SELECT load_date, source_date, client_code, ticker, trade_account, firm_code, firm_name, balance, acquisition_ccy, isin
        FROM cte
        WHERE settle_code = settle_max AND balance <> 0
    )
SELECT
    c.load_date,
    c.source_date,
    c.client_code,
    c.ticker,
    c.trade_account,
    c.firm_code,
    c.firm_name,
    c.balance,
    c.acquisition_ccy,
    c.isin,
    cv.iso_char_code,
    cv.minor_units,
    ROUND(
        (isnull(a.price_in_ccy, 0) * c.balance)
            + (isnull(a.accrued_int, 0) * c.balance) * f_acc.rate / f.rate,
        isnull(cv.minor_units, 2)
    ),
    ROUND(
        (isnull(a.price_in_ccy, 0) * c.balance) * f.rate / ft.rate,
        isnull(ct.minor_units, 2)
    ),
    ROUND(
        (isnull(a.accrued_int, 0) * c.balance) * f_acc.rate / ft.rate,
        isnull(ct.minor_units, 2)
    ),
    ROUND(
        ((isnull(a.price_in_ccy, 0) * c.balance) * f.rate
            + (isnull(a.accrued_int, 0) * c.balance) * f_acc.rate) / ft.rate,
        isnull(ct.minor_units, 2)
    ),
    ct.iso_char_code,
    ct.minor_units,
    ltrim(rtrim(a.short_name)),
    a.quote_date
FROM cte_filtered c
OUTER APPLY (
    SELECT TOP 1
        price_in_ccy = case when q.instrument_type = 'Облигации'
            then (isnull(q.face_value, 0) / 100.0) * (case when isnull(q.last_price, 0) <> 0 then q.last_price else q.close_price end)
            else (case when isnull(q.last_price, 0) <> 0 then q.last_price else q.close_price end)
        end,
        q.accrued_int,
        q.short_name,
        q.quote_date,
        mv_currency      = case when q.instrument_type = 'Облигации' then coalesce(nullif(ltrim(rtrim(q.base_currency)), ''), nullif(ltrim(rtrim(q.quote_currency)), ''), nullif(ltrim(rtrim(q.currency)), '')) else isnull(q.counter_currency, q.base_currency) end,
        accrued_currency = case when q.instrument_type = 'Облигации' then isnull(q.counter_currency, q.currency) else null end
    FROM quik.current_quotes q
    WHERE q.ticker = c.ticker
    ORDER BY case when c.acquisition_ccy = q.base_currency and c.acquisition_ccy = q.counter_currency then 0
        when c.acquisition_ccy = q.base_currency then 1
        when c.acquisition_ccy = q.counter_currency then 2
        else 3 end
) a
CROSS APPLY (VALUES (
    CASE WHEN UPPER(LTRIM(RTRIM(ISNULL(a.mv_currency,      '')))) IN ('SUR','RUR') THEN 'RUB' ELSE UPPER(LTRIM(RTRIM(ISNULL(a.mv_currency,      '')))) END,
    CASE WHEN UPPER(LTRIM(RTRIM(ISNULL(a.accrued_currency, '')))) IN ('SUR','RUR') THEN 'RUB' ELSE UPPER(LTRIM(RTRIM(ISNULL(a.accrued_currency, '')))) END,
    CASE WHEN UPPER(LTRIM(RTRIM(ISNULL(@p2, 'RUB'))))             IN ('SUR','RUR') THEN 'RUB' ELSE UPPER(LTRIM(RTRIM(ISNULL(@p2, 'RUB'))))             END
)) norm(norm_mv, norm_acc, norm_tgt)
LEFT JOIN dbo.external_codes ec_mv  ON ec_mv.ext_system_id  = 2 AND ec_mv.ext_code_type_id = 1 AND ec_mv.ext_code  = norm.norm_mv
LEFT JOIN dbo.external_codes ec_acc ON ec_acc.ext_system_id = 2 AND ec_acc.ext_code_type_id = 1 AND ec_acc.ext_code = norm.norm_acc
LEFT JOIN currencies cv_ec  ON cv_ec.iso_code  = ec_mv.internal_id
LEFT JOIN currencies ca_ec  ON ca_ec.iso_code  = ec_acc.internal_id
LEFT JOIN currencies cv_iso ON cv_iso.iso_char_code = norm.norm_mv
LEFT JOIN currencies ca_iso ON ca_iso.iso_char_code = norm.norm_acc
CROSS APPLY (VALUES (
    COALESCE(cv_ec.iso_char_code, cv_iso.iso_char_code),
    COALESCE(cv_ec.minor_units,   cv_iso.minor_units)
)) cv(iso_char_code, minor_units)
CROSS APPLY (VALUES (
    COALESCE(ca_ec.iso_char_code, ca_iso.iso_char_code),
    COALESCE(ca_ec.minor_units,   ca_iso.minor_units)
)) ca(iso_char_code, minor_units)
LEFT JOIN currencies ct ON ct.iso_char_code = norm.norm_tgt
CROSS APPLY dbo.fnFxRateToRub(ISNULL(cv.iso_char_code, ''), c.load_date) f
CROSS APPLY dbo.fnFxRateToRub(ISNULL(ca.iso_char_code, ''), c.load_date) f_acc
CROSS APPLY dbo.fnFxRateToRub(norm.norm_tgt, c.load_date) ft
ORDER BY c.load_date, c.client_code, c.ticker, c.trade_account, c.firm_code
`

func (r *Repository) SelectSecuritiesPortfolio(ctx context.Context, date time.Time, targetCcy string) ([]qmodels.PortfolioEntry, error) {
	var result []qmodels.PortfolioEntry
	r.Logger.Debug("запрос портфеля по бумагам", zap.Time("date", date), zap.String("target_ccy", targetCcy))

	rows, err := r.Db.QueryContext(ctx, selectSecuritiesPortfolio, date, targetCcy)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("ошибка запроса портфеля по бумагам", zap.Time("date", date), zap.String("target_ccy", targetCcy), zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	defer rows.Close()

	dateTrunc := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	for rows.Next() {
		var row qmodels.PortfolioEntry
		var mvISOCharCode, targetISOCharCode, shortName sql.NullString
		var mvMinorUnits, targetMinorUnits sql.NullInt32
		var quoteDate sql.NullTime

		err = rows.Scan(
			&row.LoadDate, &row.SourceDate, &row.ClientCode, &row.Instrument, &row.TradeAccount,
			&row.FirmCode, &row.FirmName, &row.Balance, &row.AcquisitionCcy,
			&row.ISIN,
			&mvISOCharCode, &mvMinorUnits,
			&row.MvInCcy, &row.MvPrice, &row.MvAccrued, &row.MvTotal,
			&targetISOCharCode, &targetMinorUnits,
			&shortName,
			&quoteDate,
		)
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}
			r.Logger.Error("ошибка при сканировании строки портфеля по бумагам", zap.Time("date", date), zap.Error(err))
			return nil, models.ErrRetrievingData
		}
		row.LimitType = qmodels.LimitTypeSecurities
		if quoteDate.Valid {
			row.QuoteDate = &quoteDate.Time
			qt := time.Date(quoteDate.Time.Year(), quoteDate.Time.Month(), quoteDate.Time.Day(), 0, 0, 0, 0, quoteDate.Time.Location())
			if qt.Before(dateTrunc) {
				r.Logger.Warn("устаревшая котировка для портфеля",
					zap.String("ticker", row.Instrument),
					zap.Time("quote_date", quoteDate.Time),
					zap.Time("load_date", date))
			}
		}
		row.MvCurrency = mvISOCharCode.String
		if targetISOCharCode.Valid {
			row.TargetCurrency = targetISOCharCode.String
		} else {
			row.TargetCurrency = targetCcy
		}
		if shortName.Valid {
			row.ShortName = &shortName.String
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		r.Logger.Error("ошибка при чтении портфеля по бумагам", zap.Time("date", date), zap.Error(rows.Err()))
		return nil, models.ErrRetrievingData
	}
	if len(result) == 0 {
		r.Logger.Warn("позиции портфеля по бумагам не найдены", zap.Time("date", date))
	} else {
		r.Logger.Debug("портфель по бумагам получен", zap.Time("date", date), zap.Int("count", len(result)))
	}
	return result, nil
}

const selectSecuritiesOtcPortfolio = `
WITH
    cte AS (
        SELECT
            li.load_date,
            li.source_date,
            li.client_code,
            li.ticker,
            li.trade_account,
            li.firm_code,
            li.settle_code,
            li.firm_name,
            li.balance,
            li.acquisition_ccy,
            li.isin,
            settle_max = MAX(li.settle_code) OVER (
                PARTITION BY li.load_date, li.client_code, li.ticker, li.trade_account, li.firm_code
            )
        FROM quik.security_limits_otc li
        WHERE li.load_date = cast(@p1 as date)
    ),
    cte_filtered AS (
        SELECT load_date, source_date, client_code, ticker, trade_account, firm_code, firm_name, balance, acquisition_ccy, isin
        FROM cte
        WHERE settle_code = settle_max AND balance <> 0
    )
SELECT
    c.load_date,
    c.source_date,
    c.client_code,
    c.ticker,
    c.trade_account,
    c.firm_code,
    c.firm_name,
    c.balance,
    c.acquisition_ccy,
    c.isin,
    cv.iso_char_code,
    cv.minor_units,
    ROUND(
        (isnull(a.price_in_ccy, 0) * c.balance)
            + (isnull(a.accrued_int, 0) * c.balance) * f_acc.rate / f.rate,
        isnull(cv.minor_units, 2)
    ),
    ROUND(
        (isnull(a.price_in_ccy, 0) * c.balance) * f.rate / ft.rate,
        isnull(ct.minor_units, 2)
    ),
    ROUND(
        (isnull(a.accrued_int, 0) * c.balance) * f_acc.rate / ft.rate,
        isnull(ct.minor_units, 2)
    ),
    ROUND(
        ((isnull(a.price_in_ccy, 0) * c.balance) * f.rate
            + (isnull(a.accrued_int, 0) * c.balance) * f_acc.rate) / ft.rate,
        isnull(ct.minor_units, 2)
    ),
    ct.iso_char_code,
    ct.minor_units,
    ltrim(rtrim(a.short_name)),
    a.quote_date
FROM cte_filtered c
OUTER APPLY (
    SELECT TOP 1
        price_in_ccy = case when q.instrument_type = 'Облигации'
            then (isnull(q.face_value, 0) / 100.0) * (case when isnull(q.last_price, 0) <> 0 then q.last_price else q.close_price end)
            else (case when isnull(q.last_price, 0) <> 0 then q.last_price else q.close_price end)
        end,
        q.accrued_int,
        q.short_name,
        q.quote_date,
        mv_currency      = case when q.instrument_type = 'Облигации' then coalesce(nullif(ltrim(rtrim(q.base_currency)), ''), nullif(ltrim(rtrim(q.quote_currency)), ''), nullif(ltrim(rtrim(q.currency)), '')) else isnull(q.counter_currency, q.base_currency) end,
        accrued_currency = case when q.instrument_type = 'Облигации' then isnull(q.counter_currency, q.currency) else null end
    FROM quik.current_quotes q
    WHERE q.ticker = c.ticker
    ORDER BY case when c.acquisition_ccy = q.base_currency and c.acquisition_ccy = q.counter_currency then 0
        when c.acquisition_ccy = q.base_currency then 1
        when c.acquisition_ccy = q.counter_currency then 2
        else 3 end
) a
CROSS APPLY (VALUES (
    CASE WHEN UPPER(LTRIM(RTRIM(ISNULL(a.mv_currency,      '')))) IN ('SUR','RUR') THEN 'RUB' ELSE UPPER(LTRIM(RTRIM(ISNULL(a.mv_currency,      '')))) END,
    CASE WHEN UPPER(LTRIM(RTRIM(ISNULL(a.accrued_currency, '')))) IN ('SUR','RUR') THEN 'RUB' ELSE UPPER(LTRIM(RTRIM(ISNULL(a.accrued_currency, '')))) END,
    CASE WHEN UPPER(LTRIM(RTRIM(ISNULL(@p2, 'RUB'))))             IN ('SUR','RUR') THEN 'RUB' ELSE UPPER(LTRIM(RTRIM(ISNULL(@p2, 'RUB'))))             END
)) norm(norm_mv, norm_acc, norm_tgt)
LEFT JOIN dbo.external_codes ec_mv  ON ec_mv.ext_system_id  = 2 AND ec_mv.ext_code_type_id = 1 AND ec_mv.ext_code  = norm.norm_mv
LEFT JOIN dbo.external_codes ec_acc ON ec_acc.ext_system_id = 2 AND ec_acc.ext_code_type_id = 1 AND ec_acc.ext_code = norm.norm_acc
LEFT JOIN currencies cv_ec  ON cv_ec.iso_code  = ec_mv.internal_id
LEFT JOIN currencies ca_ec  ON ca_ec.iso_code  = ec_acc.internal_id
LEFT JOIN currencies cv_iso ON cv_iso.iso_char_code = norm.norm_mv
LEFT JOIN currencies ca_iso ON ca_iso.iso_char_code = norm.norm_acc
CROSS APPLY (VALUES (
    COALESCE(cv_ec.iso_char_code, cv_iso.iso_char_code),
    COALESCE(cv_ec.minor_units,   cv_iso.minor_units)
)) cv(iso_char_code, minor_units)
CROSS APPLY (VALUES (
    COALESCE(ca_ec.iso_char_code, ca_iso.iso_char_code),
    COALESCE(ca_ec.minor_units,   ca_iso.minor_units)
)) ca(iso_char_code, minor_units)
LEFT JOIN currencies ct ON ct.iso_char_code = norm.norm_tgt
CROSS APPLY dbo.fnFxRateToRub(ISNULL(cv.iso_char_code, ''), c.load_date) f
CROSS APPLY dbo.fnFxRateToRub(ISNULL(ca.iso_char_code, ''), c.load_date) f_acc
CROSS APPLY dbo.fnFxRateToRub(norm.norm_tgt, c.load_date) ft
ORDER BY c.load_date, c.client_code, c.ticker, c.trade_account, c.firm_code
`

func (r *Repository) SelectSecuritiesOtcPortfolio(ctx context.Context, date time.Time, targetCcy string) ([]qmodels.PortfolioEntry, error) {
	var result []qmodels.PortfolioEntry
	r.Logger.Debug("запрос портфеля по OTC-бумагам", zap.Time("date", date), zap.String("target_ccy", targetCcy))

	rows, err := r.Db.QueryContext(ctx, selectSecuritiesOtcPortfolio, date, targetCcy)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("ошибка запроса портфеля по OTC-бумагам", zap.Time("date", date), zap.String("target_ccy", targetCcy), zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	defer rows.Close()

	dateTrunc := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	for rows.Next() {
		var row qmodels.PortfolioEntry
		var mvISOCharCode, targetISOCharCode, shortName sql.NullString
		var mvMinorUnits, targetMinorUnits sql.NullInt32
		var quoteDate sql.NullTime

		err = rows.Scan(
			&row.LoadDate, &row.SourceDate, &row.ClientCode, &row.Instrument, &row.TradeAccount,
			&row.FirmCode, &row.FirmName, &row.Balance, &row.AcquisitionCcy,
			&row.ISIN,
			&mvISOCharCode, &mvMinorUnits,
			&row.MvInCcy, &row.MvPrice, &row.MvAccrued, &row.MvTotal,
			&targetISOCharCode, &targetMinorUnits,
			&shortName,
			&quoteDate,
		)
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}
			r.Logger.Error("ошибка при сканировании строки портфеля по OTC-бумагам", zap.Time("date", date), zap.Error(err))
			return nil, models.ErrRetrievingData
		}
		row.LimitType = qmodels.LimitTypeSecuritiesOtc
		if quoteDate.Valid {
			row.QuoteDate = &quoteDate.Time
			qt := time.Date(quoteDate.Time.Year(), quoteDate.Time.Month(), quoteDate.Time.Day(), 0, 0, 0, 0, quoteDate.Time.Location())
			if qt.Before(dateTrunc) {
				r.Logger.Warn("устаревшая котировка для OTC-портфеля",
					zap.String("ticker", row.Instrument),
					zap.Time("quote_date", quoteDate.Time),
					zap.Time("load_date", date))
			}
		}
		row.MvCurrency = mvISOCharCode.String
		if targetISOCharCode.Valid {
			row.TargetCurrency = targetISOCharCode.String
		} else {
			row.TargetCurrency = targetCcy
		}
		if shortName.Valid {
			row.ShortName = &shortName.String
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		r.Logger.Error("ошибка при чтении портфеля по OTC-бумагам", zap.Time("date", date), zap.Error(rows.Err()))
		return nil, models.ErrRetrievingData
	}
	if len(result) == 0 {
		r.Logger.Warn("позиции портфеля по OTC-бумагам не найдены", zap.Time("date", date))
	} else {
		r.Logger.Debug("портфель по OTC-бумагам получен", zap.Time("date", date), zap.Int("count", len(result)))
	}
	return result, nil
}

const selectMoneyLimitsPortfolio = `
WITH cte AS (
    SELECT
        li.load_date,
        li.source_date,
        li.client_code,
        li.ccy,
        li.position_code,
        li.firm_code,
        li.settle_code,
        li.firm_name,
        li.balance,
        settle_max = MAX(li.settle_code) OVER (
            PARTITION BY li.load_date, li.client_code, li.ccy, li.position_code, li.firm_code
        )
    FROM quik.money_limits li
    WHERE li.load_date = cast(@p1 as date)
)
SELECT
    c.load_date,
    c.source_date,
    c.client_code,
    c.ccy,
    c.position_code,
    c.firm_code,
    c.firm_name,
    c.balance,
    ROUND(c.balance * f.rate / ft.rate, isnull(ct.minor_units, 2)) as balance_target,
    f.rate_date        as quote_date,
    cv.iso_char_code   as mv_iso_char_code,
    cv.minor_units     as mv_minor_units,
    cv.currency_name   as short_name,
    ct.iso_char_code   as tgt_iso_char_code,
    ct.minor_units     as tgt_minor_units
FROM cte c
CROSS APPLY (VALUES (
    CASE WHEN UPPER(LTRIM(RTRIM(c.ccy)))              IN ('SUR','RUR') THEN 'RUB' ELSE UPPER(LTRIM(RTRIM(c.ccy)))              END,
    CASE WHEN UPPER(LTRIM(RTRIM(ISNULL(@p2, 'RUB')))) IN ('SUR','RUR') THEN 'RUB' ELSE UPPER(LTRIM(RTRIM(ISNULL(@p2, 'RUB')))) END
)) norm_ccy(code, tgt_code)
LEFT JOIN dbo.external_codes ec_ccy ON ec_ccy.ext_system_id = 2 AND ec_ccy.ext_code_type_id = 1 AND ec_ccy.ext_code = norm_ccy.code
LEFT JOIN currencies cv_ec  ON cv_ec.iso_code  = ec_ccy.internal_id
LEFT JOIN currencies cv_iso ON cv_iso.iso_char_code = norm_ccy.code
CROSS APPLY (VALUES (
    COALESCE(cv_ec.iso_char_code,   cv_iso.iso_char_code),
    COALESCE(cv_ec.minor_units,     cv_iso.minor_units),
    COALESCE(cv_ec.currency_name,   cv_iso.currency_name)
)) cv(iso_char_code, minor_units, currency_name)
LEFT JOIN currencies ct ON ct.iso_char_code = norm_ccy.tgt_code
CROSS APPLY dbo.fnFxRateToRub(ISNULL(cv.iso_char_code, ''), c.load_date) f
CROSS APPLY dbo.fnFxRateToRub(norm_ccy.tgt_code,            c.load_date) ft
WHERE c.settle_code = c.settle_max AND c.balance <> 0
ORDER BY c.load_date, c.client_code, c.ccy, c.position_code, c.firm_code;
`

func (r *Repository) SelectMoneyLimitsPortfolio(ctx context.Context, date time.Time, targetCcy string) ([]qmodels.PortfolioEntry, error) {
	r.Logger.Debug("запрос денежных позиций для портфеля", zap.Time("date", date), zap.String("target_ccy", targetCcy))

	rows, err := r.Db.QueryContext(ctx, selectMoneyLimitsPortfolio, date, targetCcy)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("ошибка запроса денежных позиций для портфеля", zap.Time("date", date), zap.String("target_ccy", targetCcy), zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	defer rows.Close()

	var result []qmodels.PortfolioEntry
	for rows.Next() {
		var (
			ccy, positionCode, firmCode, firmName string
			balanceTarget                         decimal.Decimal
			quoteDate                             sql.NullTime
			mvISO, tgtISO, shortName              sql.NullString
			mvMinorUnits, tgtMinorUnits           sql.NullInt32
			row                                   qmodels.PortfolioEntry
		)
		err = rows.Scan(
			&row.LoadDate, &row.SourceDate, &row.ClientCode, &ccy, &positionCode,
			&firmCode, &firmName, &row.Balance,
			&balanceTarget,
			&quoteDate,
			&mvISO, &mvMinorUnits,
			&shortName,
			&tgtISO, &tgtMinorUnits,
		)
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}
			r.Logger.Error("ошибка чтения денежной позиции для портфеля", zap.Time("date", date), zap.Error(err))
			return nil, models.ErrRetrievingData
		}

		row.LimitType = qmodels.LimitTypeMoney
		row.Instrument = ccy
		row.PositionCode = positionCode
		row.FirmCode = firmCode
		row.FirmName = firmName
		row.MvInCcy = row.Balance
		row.MvTotal = balanceTarget
		if quoteDate.Valid {
			row.QuoteDate = &quoteDate.Time
		} else {
			row.QuoteDate = &row.LoadDate // fallback: курс не найден (rate=1), показываем дату лимита
		}
		if shortName.Valid {
			row.ShortName = &shortName.String
		}

		row.MvCurrency = mvISO.String
		if tgtISO.Valid {
			row.TargetCurrency = tgtISO.String
		} else {
			row.TargetCurrency = targetCcy
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		r.Logger.Error("ошибка чтения денежных позиций для портфеля", zap.Time("date", date), zap.Error(rows.Err()))
		return nil, models.ErrRetrievingData
	}
	if len(result) == 0 {
		r.Logger.Warn("денежные позиции для портфеля не найдены", zap.Time("date", date))
	} else {
		r.Logger.Debug("денежные позиции для портфеля получены", zap.Time("date", date), zap.Int("count", len(result)))
	}
	return result, nil
}
