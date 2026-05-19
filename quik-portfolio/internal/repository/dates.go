package repository

import (
	"context"
	"time"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"go.uber.org/zap"
)

const (
	getMoneyLimitsMaxDate = `SELECT max(load_date) FROM quik.money_limits`

	getSecurityLimitsOtcMaxDate = `
	select max(load_date) FROM quik.security_limits_otc 
	`

	selectSecurityLimitsMaxDate = `SELECT max(load_date) FROM quik.security_limits`
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
