package handler

import (
	"errors"
	"net/http"

	"github.com/boldlogic/packages/transport/httpserver/response"
	"github.com/boldlogic/packages/transport/httputils"
	"github.com/boldlogic/packages/utils/converters"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
)

type Adapter interface {
	Adapt(fn HandlerFunc) http.HandlerFunc
}

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

type HandlerFunc func(r *http.Request) (any, string, error)

func (h *Handler) Adapt(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, detail, err := fn(r)
		if err != nil {
			var status int
			switch {
			case errors.Is(err, httputils.ErrUnsupportedMediaType):
				status = http.StatusUnsupportedMediaType
			case errors.Is(err, httputils.ErrRequestEntityTooLarge):
				status = http.StatusRequestEntityTooLarge
			case errors.Is(err, models.ErrValidation) || errors.Is(err, httputils.ErrReadingBody) || errors.Is(err, converters.ErrWrongJSON):
				status = http.StatusBadRequest
			case errors.Is(err, models.ErrPartialSuccess):
				status = http.StatusMultiStatus
			case errors.Is(err, models.ErrBusinessValidation):
				status = http.StatusUnprocessableEntity
			case errors.Is(err, models.ErrNotFound):
				status = http.StatusNotFound
			case errors.Is(err, models.ErrConflict):
				status = http.StatusConflict
			default:
				status = http.StatusInternalServerError
			}
			response.WriteResp(w, status, response.Problem(status, "", detail))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == http.MethodPost {
			response.WriteResp(w, http.StatusCreated, data)
		} else {
			response.WriteResp(w, http.StatusOK, data)
		}
	}
}
