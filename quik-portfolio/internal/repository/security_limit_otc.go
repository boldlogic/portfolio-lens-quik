package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	mssql "github.com/microsoft/go-mssqldb"
	"go.uber.org/zap"
)

const (
	insertSecurityLimitOtc = `
		WITH src AS
		(
			SELECT  load_date = @p1
				,client_code = @p2
				,ticker = @p3
				,trade_account = @p4
				,settle_code = @p5
				,firm_name = @p6
				,balance = @p7
				,acquisition_ccy = @p8
				,isin = @p9
		)
		INSERT INTO quik.security_limits_otc (load_date, client_code, ticker, trade_account, settle_code, firm_code, firm_name, balance, acquisition_ccy, isin, source_date)
		OUTPUT inserted.load_date, inserted.source_date, inserted.client_code, inserted.ticker, inserted.trade_account, inserted.settle_code, inserted.firm_code, inserted.firm_name, inserted.balance, inserted.acquisition_ccy, inserted.isin
		SELECT  src.load_date
			,src.client_code
			,src.ticker
			,src.trade_account
			,src.settle_code
			,f.code
			,src.firm_name
			,src.balance
			,src.acquisition_ccy
			,src.isin
			,src.load_date
		FROM src
		CROSS APPLY (
			SELECT TOP 1 code
			FROM quik.firms
			WHERE name = src.firm_name
		) f
	`
	getSecurityLimitsOtc = `
			WITH cte AS (
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
		)
		SELECT
			cte.load_date,
			cte.source_date,
			cte.client_code,
			cte.ticker,
			cte.settle_code,
			cte.trade_account,
			cte.firm_code,
			cte.firm_name,
			cte.balance,
			cte.acquisition_ccy,
			cte.isin,
			short_disp.short_name
		FROM cte
		OUTER APPLY (
			SELECT TOP 1 ltrim(rtrim(q.short_name)) AS short_name
			FROM quik.current_quotes q
			WHERE q.ticker = cte.ticker
			ORDER BY CASE
				WHEN cte.acquisition_ccy = q.base_currency AND cte.acquisition_ccy = q.counter_currency THEN 0
				WHEN cte.acquisition_ccy = q.base_currency THEN 1
				WHEN cte.acquisition_ccy = q.counter_currency THEN 2
				ELSE 3 END
		) short_disp
		WHERE 1=1
		ORDER BY cte.load_date, cte.client_code, cte.ticker, cte.trade_account, cte.firm_code
	`
	getSecurityLimitsOtcMaxDate = `
	select max(load_date) FROM quik.security_limits_otc 
	`
	deleteSecurityLimitsOtcBeforeDate = `
		DELETE FROM
			quik.security_limits_otc
		WHERE
			load_date<CAST(@p1 AS date)	
	`
	rollSecurityLimitsOtcFromDateToDate = `
		INSERT INTO
			quik.security_limits_otc (
				load_date, client_code, ticker, trade_account, settle_code,
				firm_code, firm_name, balance, acquisition_ccy, isin, source_date
			)
		SELECT
			CAST(@p1 AS date), client_code, ticker, trade_account, settle_code,
			firm_code, firm_name, balance, acquisition_ccy, isin, source_date
		FROM
			quik.security_limits_otc
		WHERE
			load_date = CAST(@p2 AS date) AND balance <> 0
		ORDER BY load_date, client_code, ticker, trade_account, firm_code
	`
)

func (r *Repository) DeleteSecurityLimitsOtc(ctx context.Context, date time.Time) error {
	_, err := r.Db.ExecContext(ctx, deleteSecurityLimitsOtcBeforeDate, date)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.Logger.Error("ошибка при удалении лимитов OTC по бумагам", zap.Time("load_date", date), zap.Error(err))
		return models.ErrSavingData
	}
	r.Logger.Debug("лимиты OTC по бумаге успешно удалены", zap.Time("load_date", date))
	return nil
}

func (r *Repository) InsertSecurityLimitsOtcCopy(ctx context.Context, dateFrom time.Time, dateTo time.Time) error {
	_, err := r.Db.ExecContext(ctx, rollSecurityLimitsOtcFromDateToDate, dateTo, dateFrom)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		var msErr mssql.Error
		if errors.As(err, &msErr) && (msErr.Number == 2627 || msErr.Number == 2601) {
			r.Logger.Warn("лимит OTC по бумаге уже существует", zap.Time("load_date", dateTo))
			return models.ErrConflict
		}
		r.Logger.Error("ошибка при создании лимита OTC по бумаге", zap.Time("load_date", dateTo), zap.Error(err))
		return models.ErrSavingData
	}
	r.Logger.Debug("лимит OTC по бумаге успешно сохранен", zap.Time("load_date", dateTo))
	return nil
}

func (r *Repository) InsertSecurityLimitOtc(ctx context.Context, s quik.SecurityLimit) (quik.SecurityLimit, error) {
	var out quik.SecurityLimit
	row := r.Db.QueryRowContext(ctx, insertSecurityLimitOtc,
		s.LoadDate, s.ClientCode, s.Ticker, s.TradeAccount, string(s.SettleCode),
		s.FirmName, s.Balance, s.AcquisitionCcy, s.ISIN)
	err := row.Scan(&out.LoadDate, &out.SourceDate, &out.ClientCode, &out.Ticker, &out.TradeAccount, &out.SettleCode, &out.FirmCode, &out.FirmName, &out.Balance, &out.AcquisitionCcy, &out.ISIN)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.SecurityLimit{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Warn("фирма не найдена", zap.String("firm_name", s.FirmName))
			return quik.SecurityLimit{}, models.ErrNotFound
		}
		var msErr mssql.Error
		if errors.As(err, &msErr) && (msErr.Number == 2627 || msErr.Number == 2601) {
			r.Logger.Warn("лимит OTC по бумаге уже существует",
				zap.Time("load_date", s.LoadDate), zap.String("client_code", s.ClientCode), zap.String("ticker", s.Ticker))
			return quik.SecurityLimit{}, models.ErrConflict
		}
		r.Logger.Error("ошибка при создании лимита OTC по бумаге",
			zap.Time("load_date", s.LoadDate), zap.String("client_code", s.ClientCode), zap.String("ticker", s.Ticker), zap.Error(err))
		return quik.SecurityLimit{}, models.ErrSavingData
	}
	r.Logger.Debug("лимит OTC по бумаге успешно сохранен",
		zap.Time("load_date", s.LoadDate), zap.String("client_code", s.ClientCode), zap.String("ticker", s.Ticker))
	return out, nil
}

func (r *Repository) SelectSecurityLimitsOtcMaxDate(ctx context.Context) (*time.Time, error) {
	var date *time.Time
	row := r.Db.QueryRowContext(ctx, getSecurityLimitsOtcMaxDate)
	err := row.Scan(&date)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("ошибка получения максимальной даты из quik.security_limits_otc", zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	return date, nil
}

func (r *Repository) SelectSecurityLimitsOtc(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error) {
	var result []quik.SecurityLimit
	rows, err := r.Db.QueryContext(ctx, getSecurityLimitsOtc, date)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("ошибка запроса позиций OTC по бумагам", zap.Time("load_date", date), zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	defer rows.Close()

	for rows.Next() {
		row := quik.SecurityLimit{}
		var shortName sql.NullString
		err = rows.Scan(
			&row.LoadDate,
			&row.SourceDate,
			&row.ClientCode,
			&row.Ticker,
			&row.SettleCode,
			&row.TradeAccount,
			&row.FirmCode,
			&row.FirmName,
			&row.Balance,
			&row.AcquisitionCcy,
			&row.ISIN,
			&shortName,
		)
		if shortName.Valid {
			row.ShortName = shortName.String
		}
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}
			r.Logger.Error("ошибка при чтении позиции OTC по бумагам", zap.Time("load_date", date), zap.Error(err))
			return nil, models.ErrRetrievingData
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		r.Logger.Error("ошибка чтения позиций OTC по бумагам", zap.Time("load_date", date), zap.Error(rows.Err()))
		return nil, models.ErrRetrievingData
	}
	r.Logger.Debug("результаты получения позиций OTC по бумагам", zap.Time("load_date", date), zap.Int("count", len(result)))
	if len(result) == 0 {
		r.Logger.Warn("позиции OTC по бумагам не найдены", zap.Time("load_date", date))
	}
	return result, nil
}
