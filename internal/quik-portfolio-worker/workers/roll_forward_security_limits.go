package workers

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/boldlogic/packages/periodic"
	"go.uber.org/zap"
)

type RollForwardSecurityLimitsRunner interface {
	DoRollForwardSecurityLimits(ctx context.Context) error
}

func NewRollForwardSecurityLimitsWorker(svc RollForwardSecurityLimitsRunner, logger *zap.Logger, interval time.Duration) periodic.Worker {
	var lastUpToDate atomic.Int64
	return periodic.NewPeriodicWorker(
		"roll_forward_security_limits",
		"ошибка переноса лимитов по бумагам",
		interval,
		withRollForwardDateCache(&lastUpToDate, svc.DoRollForwardSecurityLimits),
		logger,
	)
}
