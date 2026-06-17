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
		r.Get("/money-limits", handler.Adapt(handler.getMoneyLimits))
		r.Get("/security-limits", handler.Adapt(handler.getSecurityLimits))
		r.Get("/otc-security-limits", handler.Adapt(handler.getSecurityLimitsOtc))
		r.Get("/positions", handler.Adapt(handler.getPositions))
	})
	return &Router{
		mux:    r,
		logger: logger,
	}
}
