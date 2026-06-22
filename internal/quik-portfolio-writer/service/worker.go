package service

import (
	"context"
	"fmt"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	"go.uber.org/zap"
)

type worker struct {
	repo      Repository
	logger    *zap.Logger
	batchSize int
	queueSize uint
	interval  time.Duration
	jobsQueue chan writeJob
}

func newWorker(repo Repository,
	logger *zap.Logger,
	batchSize int,
	queueSize uint,
	interval uint) *worker {
	return &worker{
		logger:    logger,
		batchSize: batchSize,
		queueSize: queueSize,
		interval:  time.Duration(interval) * time.Millisecond,
		jobsQueue: make(chan writeJob, queueSize),
		repo:      repo,
	}

}

type writeJob struct {
	limits []quik.Limit
	resCh  chan error
}

func (w *worker) publish(ctx context.Context, limits ...quik.Limit) (chan error, error) {

	ch := make(chan error, 1)

	job := writeJob{
		limits: limits,
		resCh:  ch,
	}

	select {
	case w.jobsQueue <- job:
		return ch, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (w *worker) Work(ctx context.Context) error {

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var pending = make(map[[32]byte]quik.Limit, w.batchSize)
	var signals = make(map[[32]byte][]chan error, w.batchSize)
	for {
		select {
		case job := <-w.jobsQueue:
			for _, limit := range job.limits {
				setKeys(limit, job.resCh, pending, signals)
			}
			if len(pending) >= w.batchSize {
				w.logger.Debug(fmt.Sprintf("по размеру %d", len(pending)))
				_ = w.flush(ctx, pending, signals)
			}
		case <-ticker.C:
			if len(pending) == 0 {
				continue
			}
			w.logger.Debug(fmt.Sprintf("по времени %d", len(pending)))
			_ = w.flush(ctx, pending, signals)
		case <-ctx.Done():
			lastCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			for job := range w.jobsQueue {
				for _, limit := range job.limits {
					setKeys(limit, job.resCh, pending, signals)
				}
			}
			err := w.flush(lastCtx, pending, signals)
			if err != nil {
				return err
			}

		}
	}

}
func setKeys(limit quik.Limit, resCh chan error, pending map[[32]byte]quik.Limit, signals map[[32]byte][]chan error) {
	key := limit.KeyHash()
	pending[key] = limit

	channels, ok := signals[key]
	if !ok {
		channels = make([]chan error, 0, 10)
	}
	channels = append(channels, resCh)
	signals[key] = channels
}

func (w *worker) flush(ctx context.Context, pending map[[32]byte]quik.Limit, signals map[[32]byte][]chan error) error {
	if len(pending) == 0 {
		return nil
	}

	err := w.repo.HandleRequest(ctx, prepareBatch(pending))
	if err != nil {
		w.logger.Error(err.Error())
	}
	sendResult(signals, err)
	clear(pending)
	clear(signals)

	return err
}

func prepareBatch(pending map[[32]byte]quik.Limit) []quik.Limit {
	out := make([]quik.Limit, 0, len(pending))
	for _, lim := range pending {
		out = append(out, lim)
	}
	return out
}

func sendResult(signals map[[32]byte][]chan error, err error) {
	for _, chGroup := range signals {
		for _, ch := range chGroup {
			ch <- err
		}
	}
}
