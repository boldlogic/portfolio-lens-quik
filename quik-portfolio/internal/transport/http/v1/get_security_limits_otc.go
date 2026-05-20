package v1

import (
	"errors"
	"net/http"

	"github.com/boldlogic/packages/transport/httputils"
	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"go.uber.org/zap"
)

func (h *Handler) GetSecurityLimitsOtc(r *http.Request) (any, string, error) {
	ctx := r.Context()
	date, err := h.extractDateQueryParam(r)
	if err != nil {
		return nil, err.Error(), md.ErrValidation
	}

	limit, offset, err := httputils.ParseListPagination(r)
	if err != nil {
		return nil, err.Error(), md.ErrValidation
	}
	clients, err := h.extractClientsQueryParam(r)
	if err != nil {
		return nil, err.Error(), md.ErrValidation
	}

	sls, totalCount, err := h.service.GetSecurityLimitsOtcWithFilters(ctx, date, limit, offset, clients)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		h.logger.Error("лимиты OTC по бумагам: чтение HTTP", zap.Error(err), zap.Time("date", date))
		return nil, "", err
	}
	return securityLimitsWithPaginationToResp(sls, totalCount, limit, offset), "", nil
}
