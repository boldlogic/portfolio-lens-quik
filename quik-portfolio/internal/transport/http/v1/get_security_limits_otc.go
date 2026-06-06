package v1

import (
	"errors"
	"net/http"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"go.uber.org/zap"
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
		h.logger.Error("лимиты OTC по бумагам: чтение HTTP", zap.Error(err), zap.Time("date", q.LoadDate))
		return nil, "", err
	}
	return securityLimitsToResponseDTO(sls, q.Limit, q.Offset, totalCount, q.IncludeTotalCount), "", nil
}
