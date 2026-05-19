package repository

import (
	"context"
	"errors"
	"time"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	mssql "github.com/microsoft/go-mssqldb"
	"go.uber.org/zap"
)

const (
	selectSecurityLimitsMaxDate    = `SELECT max(load_date) FROM quik.security_limits`
	deleteSecurityLimitsBeforeDate = `
		DELETE FROM quik.security_limits WHERE load_date < CAST(@p1 AS date)
	`
	insertSecurityLimitsCopy = `
		INSERT INTO quik.security_limits (load_date, client_code, ticker, trade_account, settle_code, firm_code, firm_name, balance, acquisition_ccy, isin, source_date)
		SELECT CAST(@p1 AS date), client_code, ticker, trade_account, settle_code, firm_code, firm_name, balance, acquisition_ccy, isin, source_date
		FROM quik.security_limits
		WHERE load_date = CAST(@p2 AS date) AND balance <> 0
		ORDER BY load_date, client_code, ticker, trade_account, firm_code
	`
)

func (r *Repository) SelectSecurityLimitsMaxDate(ctx context.Context) (*time.Time, error) {
	var date *time.Time
	row := r.Db.QueryRowContext(ctx, selectSecurityLimitsMaxDate)
	err := row.Scan(&date)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("ошибка получения максимальной даты из quik.security_limits", zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	return date, nil
}

func (r *Repository) DeleteSecurityLimits(ctx context.Context, date time.Time) error {
	_, err := r.Db.ExecContext(ctx, deleteSecurityLimitsBeforeDate, date)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.Logger.Error("ошибка при удалении лимитов по бумагам", zap.Time("load_date", date), zap.Error(err))
		return models.ErrSavingData
	}
	r.Logger.Debug("лимиты по бумагам успешно удалены", zap.Time("load_date", date))
	return nil
}

func (r *Repository) InsertSecurityLimitsCopy(ctx context.Context, dateFrom time.Time, dateTo time.Time) error {
	_, err := r.Db.ExecContext(ctx, insertSecurityLimitsCopy, dateTo, dateFrom)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		var msErr mssql.Error
		if errors.As(err, &msErr) && (msErr.Number == 2627 || msErr.Number == 2601) {
			r.Logger.Warn("лимит по бумаге уже существует", zap.Time("load_date", dateTo))
			return models.ErrConflict
		}
		r.Logger.Error("ошибка при создании лимита по бумаге", zap.Time("load_date", dateTo), zap.Error(err))
		return models.ErrSavingData
	}
	r.Logger.Debug("лимит по бумаге успешно сохранен", zap.Time("load_date", dateTo))
	return nil
}
