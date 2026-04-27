package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/boldlogic/packages/shutdown"
	qmodels "github.com/boldlogic/quik-portfolio/internal/models"
	"github.com/boldlogic/quik-portfolio/pkg/models"
	mssql "github.com/microsoft/go-mssqldb"
	"go.uber.org/zap"
)

const (
	getMoneyLimitsMaxDate       = `SELECT max(load_date) FROM quik.money_limits`
	deleteMoneyLimitsBeforeDate = `DELETE FROM quik.money_limits WHERE load_date < CAST(@p1 AS date)`
	insertMoneyLimitsCopy       = `
		INSERT INTO quik.money_limits (load_date, client_code, ccy, position_code, settle_code, firm_code, firm_name, balance, source_date)
		SELECT CAST(@p1 AS date), client_code, ccy, position_code, settle_code, firm_code, firm_name, balance, source_date
		FROM quik.money_limits
		WHERE load_date = CAST(@p2 AS date) AND balance <> 0
		ORDER BY load_date, client_code, ccy, position_code, firm_code
		`
	selectMoneyLimitsByDate = `
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
			c.settle_code,
			c.firm_code,
			c.firm_name,
			c.balance
		FROM cte c
		WHERE 1=1
		ORDER BY c.load_date, c.client_code, c.ccy, c.position_code, c.firm_code;
`
)

func (r *Repository) SelectMoneyLimits(ctx context.Context, date time.Time) ([]qmodels.MoneyLimit, error) {
	var result []qmodels.MoneyLimit

	rows, err := r.Db.QueryContext(ctx, selectMoneyLimitsByDate, date)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("ошибка выполнения запроса на получение позиций по деньгам", zap.Time("load_date", date), zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	defer rows.Close()

	for rows.Next() {
		row := qmodels.MoneyLimit{}
		err = rows.Scan(
			&row.LoadDate,
			&row.SourceDate,
			&row.ClientCode,
			&row.Currency,
			&row.PositionCode,
			&row.SettleCode,
			&row.FirmCode,
			&row.FirmName,
			&row.Balance,
		)
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}
			r.Logger.Error("ошибка чтения позиций по деньгам", zap.Time("load_date", date), zap.Error(err))
			return nil, models.ErrRetrievingData
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		r.Logger.Error("ошибка чтения позиций по деньгам", zap.Time("load_date", date), zap.Error(rows.Err()))
		return nil, models.ErrRetrievingData
	}

	r.Logger.Debug("результаты получения позиций по деньгам", zap.Time("load_date", date), zap.Int("count", len(result)))

	if len(result) == 0 {
		r.Logger.Warn("позиции по деньгам не найдены", zap.Time("load_date", date))
	}
	return result, nil
}

const (
	insertMoneyLimit = `
		WITH src AS
		(
			SELECT  load_date = @p1
				,client_code = @p2
				,ccy = @p3
				,position_code = @p4
				,settle_code = @p5
				,firm_name = @p6
				,balance = @p7
		)
		INSERT INTO quik.money_limits (load_date, client_code, ccy, position_code, settle_code, firm_code, firm_name, balance, source_date)
		OUTPUT inserted.load_date, inserted.source_date, inserted.client_code, inserted.ccy, inserted.position_code, inserted.settle_code, inserted.firm_code, inserted.firm_name, inserted.balance
		SELECT  src.load_date
			,src.client_code
			,src.ccy
			,src.position_code
			,src.settle_code
			,f.code
			,src.firm_name
			,src.balance
			,src.load_date
		FROM src
		CROSS APPLY (
			SELECT TOP 1 code
			FROM quik.firms
			WHERE name = src.firm_name
		) f
		`
)

func (r *Repository) InsertMoneyLimit(ctx context.Context, s qmodels.MoneyLimit) (qmodels.MoneyLimit, error) {
	var out qmodels.MoneyLimit
	row := r.Db.QueryRowContext(ctx, insertMoneyLimit,
		s.LoadDate,
		s.ClientCode,
		s.Currency,
		s.PositionCode,
		s.SettleCode,
		s.FirmName,
		s.Balance)
	err := row.Scan(
		&out.LoadDate,
		&out.SourceDate,
		&out.ClientCode,
		&out.Currency,
		&out.PositionCode,
		&out.SettleCode,
		&out.FirmCode,
		&out.FirmName,
		&out.Balance,
	)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return qmodels.MoneyLimit{}, err
		}

		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Warn("фирма не найдена", zap.String("firm_name", s.FirmName))
			return qmodels.MoneyLimit{}, models.ErrNotFound
		}

		var msErr mssql.Error
		if errors.As(err, &msErr) && (msErr.Number == 2627 || msErr.Number == 2601) {
			r.Logger.Warn("лимит по деньгам уже существует",
				zap.Time("load_date", s.LoadDate),
				zap.String("client_code", s.ClientCode),
				zap.String("ccy", s.Currency),
				zap.String("position_code", s.PositionCode),
				zap.String("settle_code", string(s.SettleCode)),
				zap.String("firm_code", s.FirmCode))
			return qmodels.MoneyLimit{}, models.ErrConflict
		}
		r.Logger.Error("ошибка при создании лимита по деньгам",
			zap.Time("load_date", s.LoadDate),
			zap.String("client_code", s.ClientCode),
			zap.String("ccy", s.Currency),
			zap.String("position_code", s.PositionCode),
			zap.String("settle_code", string(s.SettleCode)),
			zap.String("firm_code", s.FirmCode),
			zap.Error(err))
		return qmodels.MoneyLimit{}, models.ErrSavingData
	}
	r.Logger.Debug("лимит по деньгам успешно сохранен",
		zap.Time("load_date", s.LoadDate),
		zap.String("client_code", s.ClientCode),
		zap.String("ccy", s.Currency),
		zap.String("position_code", s.PositionCode),
		zap.String("settle_code", string(s.SettleCode)),
		zap.String("firm_code", s.FirmCode))
	return out, nil
}

func (r *Repository) SelectMoneyLimitsMaxDate(ctx context.Context) (*time.Time, error) {
	var date *time.Time
	row := r.Db.QueryRowContext(ctx, getMoneyLimitsMaxDate)
	err := row.Scan(&date)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("ошибка получения максимальной даты из quik.money_limits", zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	return date, nil
}

func (r *Repository) InsertMoneyLimitsCopy(ctx context.Context, dateFrom time.Time, dateTo time.Time) error {
	_, err := r.Db.ExecContext(ctx, insertMoneyLimitsCopy, dateTo, dateFrom)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		var msErr mssql.Error
		if errors.As(err, &msErr) && (msErr.Number == 2627 || msErr.Number == 2601) {
			r.Logger.Warn("лимит по деньгам уже существует", zap.Time("load_date", dateTo))
			return models.ErrConflict
		}
		r.Logger.Error("ошибка при создании копии лимитов по деньгам", zap.Time("load_date", dateTo), zap.Error(err))
		return models.ErrSavingData
	}
	r.Logger.Debug("копия лимитов по деньгам успешно создана", zap.Time("load_date", dateTo))
	return nil
}

func (r *Repository) DeleteMoneyLimits(ctx context.Context, date time.Time) error {
	_, err := r.Db.ExecContext(ctx, deleteMoneyLimitsBeforeDate, date)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.Logger.Error("ошибка при удалении лимитов по деньгам", zap.Time("load_date", date), zap.Error(err))
		return models.ErrSavingData
	}
	r.Logger.Debug("лимиты по деньгам успешно удалены", zap.Time("load_date", date))
	return nil
}
