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

func (r *Repository) InsertSecurityLimitOtc(ctx context.Context, s quik.SecurityLimit) (quik.SecurityLimit, error) {
	var out quik.SecurityLimit
	row := r.Db.QueryRowContext(ctx, insertSecurityLimitOtc,
		s.ClientCode, s.Ticker, s.TradeAccount, string(s.SettleCode),
		s.FirmCode, s.Balance, s.AcquisitionCcy, s.ISIN)
	err := row.Scan(&out.LoadDate, &out.SourceDate, &out.ClientCode, &out.Ticker, &out.TradeAccount, &out.SettleCode, &out.FirmCode, &out.FirmName, &out.Balance, &out.AcquisitionCcy, &out.ISIN)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.SecurityLimit{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Warn("при создании лимитов не найдена фирма",
				zap.String("firm_code", s.FirmCode))
			return quik.SecurityLimit{}, models.ErrNotFound
		}
		var msErr mssql.Error
		if errors.As(err, &msErr) && (msErr.Number == 2627 || msErr.Number == 2601) {
			r.Logger.Warn("лимит OTC по бумаге уже существует",
				zap.String("client_code", s.ClientCode),
				zap.String("ticker", s.Ticker), zap.String("trade_account", s.TradeAccount),
				zap.String("settle_code", string(s.SettleCode)), zap.String("firm_name", s.FirmName))
			return quik.SecurityLimit{}, models.ErrConflict
		}
		r.Logger.Error("ошибка при создании лимита OTC по бумаге",
			zap.String("client_code", s.ClientCode),
			zap.String("ticker", s.Ticker), zap.Error(err))
		return quik.SecurityLimit{}, models.ErrSavingData
	}
	r.Logger.Debug("лимит OTC по бумаге успешно сохранен",
		zap.Time("load_date", out.LoadDate), zap.String("client_code", out.ClientCode), zap.String("ticker", out.Ticker))
	return out, nil
}
