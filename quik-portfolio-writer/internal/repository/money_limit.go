package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	mssql "github.com/microsoft/go-mssqldb"
	"go.uber.org/zap"
)

const (
	insertMoneyLimit = `
		WITH src AS
		(
			SELECT  client_code = @p1
				,currency_code = @p2
				,position_code = @p3
				,settle_code = @p4
				,firm_code = @p5
				,balance = @p6
		)
		INSERT INTO quik.money_limits (client_code, currency_code, position_code, settle_code, firm_code, firm_name, balance)
		OUTPUT inserted.load_date, inserted.source_date, inserted.client_code, inserted.currency_code, inserted.position_code, inserted.settle_code, inserted.firm_code, inserted.firm_name, inserted.balance
		SELECT  src.client_code
			,src.currency_code
			,src.position_code
			,src.settle_code
			,f.code
			,f.name
			,src.balance
		FROM src
		join ref.firms f on code = src.firm_code	
		`
)

func (r *Repository) InsertMoneyLimit(ctx context.Context, s quik.MoneyLimit) (quik.MoneyLimit, error) {
	var out quik.MoneyLimit
	row := r.Db.QueryRowContext(ctx, insertMoneyLimit,
		s.ClientCode,
		s.Currency,
		s.PositionCode,
		s.SettleCode,
		s.FirmCode,
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
			return quik.MoneyLimit{}, err
		}

		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Warn("при создании лимитов не найдена фирма",
				zap.String("firm_code", s.FirmCode))
			return quik.MoneyLimit{}, models.ErrNotFound
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
			return quik.MoneyLimit{}, models.ErrConflict
		}
		r.Logger.Error("ошибка при создании лимита по деньгам",
			zap.Time("load_date", s.LoadDate),
			zap.String("client_code", s.ClientCode),
			zap.String("ccy", s.Currency),
			zap.String("position_code", s.PositionCode),
			zap.String("settle_code", string(s.SettleCode)),
			zap.String("firm_code", s.FirmCode),
			zap.Error(err))
		return quik.MoneyLimit{}, models.ErrSavingData
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
