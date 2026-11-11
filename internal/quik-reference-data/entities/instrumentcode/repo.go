package instrumentcode

import (
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-reference-data/shared/db"
)

type InstrumentCodesRepo struct {
	repo *db.Repository
}

func NewInstrumentCodesRepo(r *db.Repository) *InstrumentCodesRepo {
	return &InstrumentCodesRepo{repo: r}
}
