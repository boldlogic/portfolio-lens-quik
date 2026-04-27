package workers

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/boldlogic/packages/periodic"
	"go.uber.org/zap"
)

type RollForwardOtcRunner interface {
	DoRollForwardOtc(ctx context.Context) error
}

func NewRollForwardOtcWorker(svc RollForwardOtcRunner, logger *zap.Logger, interval time.Duration) periodic.Worker {
	var lastUpToDate atomic.Int64
	return periodic.NewPeriodicWorker(
		"roll_forward_otc",
		"ошибка переноса OTC-лимитов",
		interval,
		withRollForwardDateCache(&lastUpToDate, svc.DoRollForwardOtc),
		logger,
	)
}
