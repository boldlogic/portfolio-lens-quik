package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Router struct {
	mux    *chi.Mux
	logger *zap.Logger
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func NewRouter(handler *Handler, logger *zap.Logger) *Router {
	r := chi.NewRouter()
	r.Post("/limits", handler.Adapt(handler.createLimit))
	r.Post("/limits/upload", handler.Adapt(handler.upload))

	return &Router{
		mux:    r,
		logger: logger,
	}
}
