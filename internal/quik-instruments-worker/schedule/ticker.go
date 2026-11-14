package schedule

import (
	"context"
	"time"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/config"
	"go.uber.org/zap"
)

type Pipeline interface {
	Run(ctx context.Context, batchSize int) error
}

type Ticker struct {
	cfg      config.WorkerConfig
	pipeline Pipeline
	logger   *zap.Logger
}

func NewTicker(cfg config.WorkerConfig, pipeline Pipeline, logger *zap.Logger) *Ticker {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Ticker{
		cfg:      cfg,
		pipeline: pipeline,
		logger:   logger,
	}
}

func (t *Ticker) Run(ctx context.Context) {
	if !t.cfg.Enabled {
		t.logger.Warn("воркер отключен", zap.String("name", t.cfg.Name))
		return
	}

	if t.cfg.RunOnStart {
		t.runOnce(ctx)
	}

	ticker := time.NewTicker(time.Duration(t.cfg.Interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.runOnce(ctx)
		}
	}
}

func (t *Ticker) runOnce(ctx context.Context) {
	runCtx, cancel := context.WithTimeout(ctx, time.Duration(t.cfg.Timeout)*time.Second)
	defer cancel()

	err := t.pipeline.Run(runCtx, t.cfg.BatchSize)
	if err != nil && !shutdown.IsExceeded(err) {
		t.logger.Error(t.cfg.Name, zap.Error(err))
	}
}
