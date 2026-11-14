package persist

import (
	"context"

	"github.com/boldlogic/portfolio-lens-quik/internal/quik-instruments-worker/model"
)

type Instruments struct {
	repo InstrumentRepository
}

func NewInstruments(repo InstrumentRepository) *Instruments {
	return &Instruments{repo: repo}
}

func (s *Instruments) Find(ctx context.Context, key model.InstrumentKey) (model.StoredInstrument, bool, error) {
	return s.repo.FindInstrument(ctx, key)
}
