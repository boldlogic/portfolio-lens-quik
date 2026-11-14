package pipeline

import (
	"context"
	"fmt"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/stats"
	"go.uber.org/zap"
)

type Job interface {
	Name() string
	Run(ctx context.Context, batchSize int) (stats.Stats, error)
}

type JobResult struct {
	Name  string
	Stats stats.Stats
}

type Result struct {
	Jobs []JobResult
}

func (r Result) StatsByName(name string) (stats.Stats, bool) {
	for _, item := range r.Jobs {
		if item.Name == name {
			return item.Stats, true
		}
	}
	return stats.Stats{}, false
}

type Runner struct {
	jobs   []Job
	logger *zap.Logger
}

func NewRunner(logger *zap.Logger, jobs ...Job) *Runner {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Runner{
		jobs:   jobs,
		logger: logger,
	}
}

func (r *Runner) Run(ctx context.Context, batchSize int) error {
	_, err := r.RunResult(ctx, batchSize)
	return err
}

func (r *Runner) RunResult(ctx context.Context, batchSize int) (Result, error) {
	if batchSize <= 0 {
		return Result{}, fmt.Errorf("размер batch должен быть больше нуля")
	}

	out := Result{Jobs: make([]JobResult, 0, len(r.jobs))}
	for _, job := range r.jobs {
		itemStats, err := job.Run(ctx, batchSize)
		if err != nil {
			return Result{}, fmt.Errorf("%s: %w", job.Name(), err)
		}
		r.logResult(job.Name(), itemStats)
		out.Jobs = append(out.Jobs, JobResult{
			Name:  job.Name(),
			Stats: itemStats,
		})
	}
	return out, nil
}

func (r *Runner) logResult(name string, itemStats stats.Stats) {
	r.logger.Info(
		name,
		zap.Int("selected", itemStats.Selected),
		zap.Int("instruments_unlinked", itemStats.InstrumentsUnlinked),
		zap.Int("quotes_linked", itemStats.QuotesLinked),
		zap.Int("quotes_skipped", itemStats.QuotesSkipped),
	)
}
