package tradeaccount

import (
	"context"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"go.uber.org/zap"
)

const (
	syncTradeAccountsFromLimits = `
		INSERT INTO ref.trade_accounts (client_id, venue, code)
		SELECT DISTINCT c.client_id, src.venue, src.trade_account
		FROM (
			SELECT firm_code, client_code, trade_account, venue = 'exchange'
			FROM quik.security_limits
			WHERE LTRIM(RTRIM(client_code)) <> ''
				AND LTRIM(RTRIM(trade_account)) <> ''
			UNION
			SELECT firm_code, client_code, trade_account, venue = 'otc'
			FROM quik.security_limits_otc
			WHERE LTRIM(RTRIM(client_code)) <> ''
				AND LTRIM(RTRIM(trade_account)) <> ''
		) src
		JOIN ref.firms f ON f.code = src.firm_code
		JOIN ref.clients c ON c.firm_id = f.firm_id AND c.code = src.client_code
		WHERE NOT EXISTS (
			SELECT 1
			FROM ref.trade_accounts t
			WHERE t.client_id = c.client_id
				AND t.venue = src.venue
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
