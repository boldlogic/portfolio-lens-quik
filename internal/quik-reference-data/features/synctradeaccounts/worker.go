package synctradeaccounts

import (
	"context"
	"time"

	"github.com/boldlogic/packages/periodic"
	"go.uber.org/zap"
)

type SyncTradeAccountsRunner interface {
	SyncTradeAccountsFromLimits(ctx context.Context) error
}

func NewActualizeTradeAccountsWorker(svc SyncTradeAccountsRunner, logger *zap.Logger, interval time.Duration) periodic.Worker {
	return periodic.NewPeriodicWorker(
		"actualize_trade_accounts",
		"ошибка синхронизации торговых счетов из лимитов",
		interval,
		func(ctx context.Context) error { return svc.SyncTradeAccountsFromLimits(ctx) },
		logger,
	)
}
