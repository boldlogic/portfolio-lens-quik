package v1

import (
	"errors"
	"net/http"

	"github.com/boldlogic/quik-portfolio/pkg/models"
	"github.com/boldlogic/quik-portfolio/pkg/transport/httpserver/httputils"
)

type firmCreateReqDTO struct {
	Code string `json:"firmCode" validate:"required,min=1,max=12"`
	Name string `json:"firmName" validate:"required,min=1,max=128"`
}

type firmRespDTO struct {
	Id   uint8  `json:"id"`
	Code string `json:"firmCode"`
	Name string `json:"firmName"`
}

func (h *Handler) CreateFirm(r *http.Request) (any, string, error) {
	ctx := r.Context()
	req, err := httputils.DecodeAndValidate[firmCreateReqDTO](r)
	if err != nil {
		if errors.Is(err, httputils.ErrUnsupportedMediaType) || errors.Is(err, httputils.ErrRequestEntityTooLarge) {
			return nil, err.Error(), err
		}
		return nil, err.Error(), models.ErrValidation
	}

	firm, err := h.service.CreateFirm(ctx, req.Code, req.Name)
	if err != nil {
		if errors.Is(err, models.ErrConflict) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return firmToResp(firm), "", nil
}
