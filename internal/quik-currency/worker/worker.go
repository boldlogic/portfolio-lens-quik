package worker

import (
	"context"
	"sync"
)

type Worker interface {
	Run(ctx context.Context)
}

type Runner struct {
	workers []Worker
	wg      sync.WaitGroup
}

func NewRunner(workers ...Worker) *Runner {
	return &Runner{workers: workers}
}

func (r *Runner) Run(ctx context.Context) {
	for _, w := range r.workers {
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			w.Run(ctx)
		}()
	}
	<-ctx.Done()
	r.wg.Wait()
}
