package v1

import (
	"errors"
	"net/http"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
)

func (h *Handler) getSecurityLimitsOtc(r *http.Request) (any, string, error) {
	q, err := parseLimitsQueryParams(r)
	if err != nil {
		return nil, err.Error(), md.ErrValidation
	}

	sls, totalCount, err := h.service.GetSecurityLimitsOtcWithFilters(
		r.Context(), q.LoadDate, q.Limit, q.Offset, q.ClientCodes, q.IncludeTotalCount,
	)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return securityLimitsToResponseDTO(sls, q.Limit, q.Offset, totalCount, q.IncludeTotalCount), "", nil
}
