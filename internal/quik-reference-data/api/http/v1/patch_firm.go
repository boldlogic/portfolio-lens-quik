package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/boldlogic/packages/transport/httputils"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) UpdateFirm(r *http.Request) (any, string, error) {
	param := chi.URLParam(r, "id")
	id64, err := strconv.ParseUint(param, 10, 8)
	if err != nil {
		return nil, fmt.Sprintf("некорректный id фирмы %s", param), models.ErrValidation
	}

	req, err := httputils.DecodeRequest[firmPatchReqDTO](r)
	if err != nil {
		if errors.Is(err, httputils.ErrUnsupportedMediaType) || errors.Is(err, httputils.ErrRequestEntityTooLarge) {
			return nil, err.Error(), err
		}
		return nil, err.Error(), models.ErrValidation
	}

	firm, err := h.writeService.UpdateFirm(r.Context(), uint8(id64), req.Name)
	if err != nil {
		return nil, err.Error(), err
	}
	return firmToResp(firm), "", nil
}
