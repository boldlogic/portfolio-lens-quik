package v1

import (
	"errors"
	"net/http"
	"strings"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

type portfolio struct {
	MVTotal   decimal.Decimal `json:"marketValueTotal"`
	Currency  string          `json:"currency"`
	Positions []position      `json:"positions"`
}

func portfolioToDTO(positions []quik.Position, total decimal.Decimal, portfolioCCY string) portfolio {

	out := portfolio{
		MVTotal:  total,
		Currency: portfolioCCY,
	}
	pos := make([]position, 0, len(positions))
	for _, p := range positions {
		pos = append(pos, toDTO(p))
	}
	out.Positions = pos
	return out
}

func (h *Handler) getMoneyPositions(r *http.Request) (any, string, error) {
	ctx := r.Context()
	var currency *string
	targetCcy := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get(currencyQuery)))
	if targetCcy != "" {
		currency = &targetCcy
	}

	positions, total, portfolioCCY, err := h.service.GetMoneyPositions(ctx, currency)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return portfolioToDTO(positions, total, portfolioCCY), "", nil
}
