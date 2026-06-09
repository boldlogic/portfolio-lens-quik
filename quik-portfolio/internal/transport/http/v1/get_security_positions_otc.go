package v1

import (
	"errors"
	"net/http"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
)

func (h *Handler) getSecurityPositionsOtc(r *http.Request) (any, string, error) {
	ctx := r.Context()
	params := parsePortfolioQueryParams(r)

	positions, total, portfolioCCY, err := h.service.GetSecurityPositionsOtc(ctx, params.TargetCurrency, params.ClientCodes)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return positionsToPortfolioDTO(positions, total, portfolioCCY), "", nil
}
