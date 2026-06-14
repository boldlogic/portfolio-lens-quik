package worker

import (
	"context"
	"time"

	"github.com/boldlogic/packages/shutdown"
	"go.uber.org/zap"
)

type Job struct {
	name       string
	enabled    bool
	runOnStart bool
	interval   time.Duration
	timeout    time.Duration
	logger     *zap.Logger
	jobFunc    JobFunc
}

func NewJob(conf JobConfig, logger *zap.Logger, jobFunc JobFunc) *Job {

	if logger == nil {
		logger = zap.NewNop()
	}

	out := Job{
		jobFunc: jobFunc,
		logger:  logger,
		name:    conf.Name,
		enabled: conf.Enabled,
	}

	if out.enabled == false {
		out.logger.Warn("воркер отключен", zap.String("name", out.name))
		return &out
	}
	out.runOnStart = conf.RunOnStart

	interval := time.Duration(conf.Interval) * time.Second
	if conf.Interval == 0 {
		logger.Warn("установлен дефолтный интервал")
		interval = time.Duration(defaultWorkerInterval) * time.Second
	}
	out.interval = interval

	timeout := time.Duration(conf.Timeout) * time.Second
	if conf.Timeout == 0 {
		logger.Warn("установлен дефолтный таймаут")
		timeout = time.Duration(defaultWorkerTimeout) * time.Second
	}
	out.timeout = timeout

	return &out

}

type JobFunc func(ctx context.Context) error

func (j *Job) Run(ctx context.Context) {
	if !j.enabled {
		return
	}
	if j.runOnStart {
		runCTX, runCancelFunc := context.WithTimeout(ctx, j.timeout)

		err := j.jobFunc(runCTX)
		runCancelFunc()
		if err != nil && !shutdown.IsExceeded(err) {
			j.logger.Error(j.name, zap.Error(err))
		}
	}
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

			runCTX, runCancelFunc := context.WithTimeout(ctx, j.timeout)

			err := j.jobFunc(runCTX)
			runCancelFunc()
			if err != nil && !shutdown.IsExceeded(err) {
				j.logger.Error(j.name, zap.Error(err))
			}
		}
	}

}
