package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type worker struct {
	repo         Repository
	logger       *zap.Logger
	batchSize    int
	queueSize    uint16
	interval     time.Duration
	jobsQueue    chan writeJob
	shutdownOnce sync.Once
}

func newWorker(repo Repository,
	logger *zap.Logger,
	batchSize int,
	queueSize uint16,
	interval uint16) *worker {
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
	uuid   uuid.UUID
	limits []quik.Limit
	resCh  chan error
}

type result struct {
	total      int
	queued     int
	succeed    int
	failed     int
	resultChan chan error
}

func (s *Service) Stop() {
	s.worker.stop()
}

func (w *worker) stop() {
	w.shutdownOnce.Do(func() {
		close(w.jobsQueue)
	})
}

func (w *worker) publish(ctx context.Context, limits ...quik.Limit) (chan error, error) {

	ch := make(chan error, 1)

	u, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	job := writeJob{
		uuid:   u,
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

func countFlushed(counter map[uuid.UUID]*result, err error) {

	for k, v := range counter {
		if err != nil {
			counter[k].failed += v.queued
		} else {
			counter[k].succeed += v.queued
		}
		counter[k].queued = 0
	}

}

func (s *Service) Run(ctx context.Context) error {
	return s.worker.work(ctx)
}

func (w *worker) work(ctx context.Context) error {

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	var pending = make(map[[32]byte]quik.Limit, w.batchSize)
	var counter = make(map[uuid.UUID]*result, w.batchSize*10)

	for {
		select {
		case job, ok := <-w.jobsQueue:
			if !ok {
				lastCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				err := w.shutdownDrain(lastCtx, pending, counter)
				return err
			}
			counter[job.uuid] = &result{
				total:      len(job.limits),
				resultChan: job.resCh,
			}

			for _, limit := range job.limits {

				pending[limit.KeyHash()] = limit
				counter[job.uuid].queued++

				if len(pending) >= w.batchSize {
					w.logger.Debug(fmt.Sprintf("по размеру %d", len(pending)))
					err := w.flush(ctx, pending)
					countFlushed(counter, err)
					sendResult(counter)
					clear(pending)
				}
			}

		case <-ticker.C:
			if len(pending) == 0 {
				continue
			}
			w.logger.Debug(fmt.Sprintf("по времени %d", len(pending)))
			err := w.flush(ctx, pending)
			countFlushed(counter, err)
			sendResult(counter)
			clear(pending)

		case <-ctx.Done():
			lastCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err := w.shutdownDrain(lastCtx, pending, counter)
			if err != nil {
				return err
			}
			return nil

		}
	}

}

func (w *worker) shutdownDrain(shutdownCTX context.Context, pending map[[32]byte]quik.Limit, counter map[uuid.UUID]*result) error {
	for job := range w.jobsQueue {
		counter[job.uuid] = &result{
			total:      len(job.limits),
			resultChan: job.resCh,
		}
		for _, limit := range job.limits {

			pending[limit.KeyHash()] = limit
			counter[job.uuid].queued++
		}

	}
	err := w.flush(shutdownCTX, pending)
	countFlushed(counter, err)
	sendResult(counter)
	clear(pending)

	for k, v := range counter {
		v.resultChan <- context.Canceled
		close(counter[k].resultChan)

	}

	return err
}

func (w *worker) flush(ctx context.Context, pending map[[32]byte]quik.Limit) error {
	if len(pending) == 0 {
		return nil
	}

	err := w.repo.HandleRequest(ctx, prepareBatch(pending))
	if err != nil {
		w.logger.Error(err.Error())
		return err
	}

	return nil
}

func prepareBatch(pending map[[32]byte]quik.Limit) []quik.Limit {
	out := make([]quik.Limit, 0, len(pending))
	for _, lim := range pending {
		out = append(out, lim)
	}
	return out
}

func sendResult(res map[uuid.UUID]*result) {

	for k, v := range res {

		switch {
		case res[k].total == res[k].succeed:
			v.resultChan <- nil
			close(res[k].resultChan)
			delete(res, k)
		case res[k].total == v.failed:
			v.resultChan <- fmt.Errorf("%w", models.ErrSavingData)
			close(res[k].resultChan)
			delete(res, k)
		case res[k].total == (res[k].failed + res[k].succeed):
			v.resultChan <- fmt.Errorf("%w: успешно=%d, с ошибкой=%d", models.ErrPartialSuccess, res[k].succeed, res[k].failed)
			close(res[k].resultChan)
			delete(res, k)
		default:
			continue
		}
	}
}
