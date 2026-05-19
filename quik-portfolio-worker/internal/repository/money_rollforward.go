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
	getMoneyLimitsMaxDate       = `SELECT max(load_date) FROM quik.money_limits`
	deleteMoneyLimitsBeforeDate = `DELETE FROM quik.money_limits WHERE load_date < CAST(@p1 AS date)`
	insertMoneyLimitsCopy       = `
		INSERT INTO quik.money_limits (load_date, client_code, ccy, position_code, settle_code, firm_code, firm_name, balance, source_date)
		SELECT CAST(@p1 AS date), client_code, ccy, position_code, settle_code, firm_code, firm_name, balance, source_date
		FROM quik.money_limits
		WHERE load_date = CAST(@p2 AS date) AND balance <> 0
		ORDER BY load_date, client_code, ccy, position_code, firm_code
		`
)

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
