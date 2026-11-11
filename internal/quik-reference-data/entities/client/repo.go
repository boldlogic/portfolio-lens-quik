package client

import (
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-reference-data/shared/db"
)

type ClientsRepo struct {
	repo *db.Repository
}

func NewClientsRepo(r *db.Repository) *ClientsRepo {
	return &ClientsRepo{repo: r}
}
