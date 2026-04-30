package v1

import (
	"errors"
	"net/http"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/httputils"
)

func (h *Handler) CreateFirm(r *http.Request) (any, string, error) {
	req, err := httputils.DecodeAndValidate[firmCreateReqDTO](r)
	if err != nil {
		if errors.Is(err, httputils.ErrUnsupportedMediaType) || errors.Is(err, httputils.ErrRequestEntityTooLarge) {
			return nil, err.Error(), err
		}
		return nil, err.Error(), models.ErrValidation
	}

	firm, err := h.writeService.CreateFirm(r.Context(), req.Code, req.Name)
	if err != nil {
		if errors.Is(err, models.ErrConflict) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return firmToResp(firm), "", nil
}
