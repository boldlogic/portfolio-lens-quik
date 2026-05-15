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
	r.Route("/quik", func(r chi.Router) {
		r.Route("/limits", func(r chi.Router) {
			r.Post("/money", handler.Adapt(handler.CreateMoneyLimit))
			r.Post("/securities", handler.Adapt(handler.CreateSecurityLimit))
			r.Post("/securities/otc", handler.Adapt(handler.CreateSecurityLimitOtc))
		})
	})
	return &Router{
		mux:    r,
		logger: logger,
	}
}
