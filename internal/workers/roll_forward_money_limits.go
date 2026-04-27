package workers

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/boldlogic/packages/periodic"
	"go.uber.org/zap"
)

type RollForwardMoneyLimitsRunner interface {
	DoRollForwardMoneyLimits(ctx context.Context) error
}

func NewRollForwardMoneyLimitsWorker(svc RollForwardMoneyLimitsRunner, logger *zap.Logger, interval time.Duration) periodic.Worker {
	var lastUpToDate atomic.Int64
	return periodic.NewPeriodicWorker(
		"roll_forward_money_limits",
		"ошибка переноса лимитов по деньгам",
		interval,
		withRollForwardDateCache(&lastUpToDate, svc.DoRollForwardMoneyLimits),
		logger,
	)
}
