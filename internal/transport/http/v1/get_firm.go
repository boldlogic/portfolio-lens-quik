package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/boldlogic/quik-portfolio/pkg/models"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetFirm(r *http.Request) (any, string, error) {
	ctx := r.Context()
	id64, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 8)
	if err != nil {
		return nil, "некорректный id фирмы", models.ErrValidation
	}

	firm, err := h.service.GetFirmByID(ctx, uint8(id64))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return firmToResp(firm), "", nil
}
