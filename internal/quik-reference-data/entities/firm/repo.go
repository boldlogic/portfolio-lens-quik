package firm

import (
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-reference-data/shared/db"
)

type FirmsRepo struct {
	repo *db.Repository
}

func NewFirmsRepo(r *db.Repository) *FirmsRepo {

	return &FirmsRepo{
		repo: r,
	}
}
