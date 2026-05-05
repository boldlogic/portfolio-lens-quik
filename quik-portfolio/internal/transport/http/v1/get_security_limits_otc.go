package v1

import (
	"errors"
	"net/http"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
)

func (h *Handler) GetSecurityLimitsOtc(r *http.Request) (any, string, error) {
	ctx := r.Context()
	date, err := h.extractDateQueryParam(r)
	if err != nil {
		return nil, err.Error(), md.ErrValidation
	}

	sls, err := h.service.GetSecurityLimitsOtc(ctx, date)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return securityLimitsToResp(sls), "", nil
}
