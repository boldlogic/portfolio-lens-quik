package tradeaccount

import (
	"context"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"go.uber.org/zap"
)

const (
	syncTradeAccountsFromLimits = `
		INSERT INTO ref.trade_accounts (firm_id, client_code, code)
		SELECT DISTINCT f.firm_id, src.client_code, src.trade_account
		FROM (
			SELECT firm_code, client_code, trade_account
			FROM quik.security_limits
			WHERE LTRIM(RTRIM(client_code)) <> ''
				AND LTRIM(RTRIM(trade_account)) <> ''
			UNION
			SELECT firm_code, client_code, trade_account
			FROM quik.security_limits_otc
			WHERE LTRIM(RTRIM(client_code)) <> ''
				AND LTRIM(RTRIM(trade_account)) <> ''
		) src
		JOIN ref.firms f ON f.code = src.firm_code
		WHERE NOT EXISTS (
			SELECT 1
			FROM ref.trade_accounts t
			WHERE t.firm_id = f.firm_id
				AND t.client_code = src.client_code
				AND t.code = src.trade_account
		);
	`
)

func (r *TradeAccountsRepo) SyncTradeAccountsFromLimits(ctx context.Context) error {
	_, err := r.repo.Db.ExecContext(ctx, syncTradeAccountsFromLimits)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.repo.Logger.Error("ошибка синхронизации торговых счетов из лимитов", zap.Error(err))
		return models.ErrSavingData
	}

	return nil
}
