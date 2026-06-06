package v1

import (
	"errors"
	"net/http"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

type portfolioDTO struct {
	MarketValueTotal decimal.Decimal `json:"marketValueTotal"`
	TargetCurrency   string          `json:"targetCurrency"`
	Positions        []positionDTO   `json:"positions"`
}

func positionsToPortfolioDTO(positions []quik.Position, total decimal.Decimal, portfolioCCY string) portfolioDTO {

	out := portfolioDTO{
		MarketValueTotal: total,
		TargetCurrency:   portfolioCCY,
	}
	pos := make([]positionDTO, 0, len(positions))
	for _, p := range positions {
		pos = append(pos, positionToDTO(p))
	}
	out.Positions = pos
	return out
}

func (h *Handler) getMoneyPositions(r *http.Request) (any, string, error) {
	ctx := r.Context()
	targetCurrency := parseCurrencyQueryParam(r)

	positions, total, portfolioCCY, err := h.service.GetMoneyPositions(ctx, targetCurrency)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return positionsToPortfolioDTO(positions, total, portfolioCCY), "", nil
}
