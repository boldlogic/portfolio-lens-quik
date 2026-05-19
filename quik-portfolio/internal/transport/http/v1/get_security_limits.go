package v1

import (
	"errors"
	"net/http"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"go.uber.org/zap"
)

func (h *Handler) GetSecurityLimits(r *http.Request) (any, string, error) {
	ctx := r.Context()
	date, err := h.extractDateQueryParam(r)
	if err != nil {
		return nil, err.Error(), md.ErrValidation
	}

	sls, err := h.service.GetSecurityLimits(ctx, date)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		h.logger.Error("лимиты бумаги: чтение HTTP", zap.Error(err), zap.Time("date", date))

		return nil, "", err
	}
	return securityLimitsToResp(sls), "", nil
}
