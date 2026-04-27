package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/boldlogic/quik-portfolio/pkg/transport/httpserver/httputils"

	"github.com/boldlogic/quik-portfolio/pkg/models"
	"github.com/go-chi/chi/v5"
)

type firmPatchReqDTO struct {
	Name string `json:"firmName" validate:"required,min=1,max=128"`
}

func (h *Handler) UpdateFirm(r *http.Request) (any, string, error) {
	ctx := r.Context()
	id64, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 8)
	if err != nil {
		return nil, "некорректный id фирмы", models.ErrValidation
	}

	req, err := httputils.DecodeAndValidate[firmPatchReqDTO](r)
	if err != nil {
		if errors.Is(err, httputils.ErrUnsupportedMediaType) || errors.Is(err, httputils.ErrRequestEntityTooLarge) {
			return nil, err.Error(), err
		}
		return nil, err.Error(), models.ErrValidation
	}

	firm, err := h.service.UpdateFirm(ctx, uint8(id64), req.Name)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, models.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return firmToResp(firm), "", nil
}
