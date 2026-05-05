package handler

import (
	"errors"
	"net/http"

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
			var resp HTTPErr
			switch {
			case errors.Is(err, httputils.ErrUnsupportedMediaType):
				resp = UnsupportedMediaType(detail)
			case errors.Is(err, httputils.ErrRequestEntityTooLarge):
				resp = RequestEntityTooLarge(detail)
			case errors.Is(err, models.ErrValidation) || errors.Is(err, httputils.ErrReadingBody) || errors.Is(err, converters.ErrWrongJSON):
				resp = BadRequest(detail)
			case errors.Is(err, models.ErrBusinessValidation):
				resp = UnprocessableEntity(detail)
			case errors.Is(err, models.ErrNotFound):
				resp = NotFound(detail)
			case errors.Is(err, models.ErrConflict):
				resp = Conflict(detail)
			default:
				resp = Internal()
			}
			WriteResp(w, resp.Status, resp)
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == http.MethodPost {
			WriteResp(w, http.StatusCreated, data)
		} else {
			WriteResp(w, http.StatusOK, data)
		}
	}
}
