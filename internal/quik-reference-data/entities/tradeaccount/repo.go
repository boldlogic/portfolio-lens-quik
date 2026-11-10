package tradeaccount

import (
	"github.com/boldlogic/portfolio-lens-quik/internal/quik-reference-data/shared/db"
)

type TradeAccountsRepo struct {
	repo *db.Repository
}

func NewTradeAccountsRepo(r *db.Repository) *TradeAccountsRepo {
	return &TradeAccountsRepo{
		repo: r,
	}
}
