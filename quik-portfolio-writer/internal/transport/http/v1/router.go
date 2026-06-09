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
	r.Post("/money-limits", handler.Adapt(handler.postMoneyLimit))
	r.Post("/security-limits", handler.Adapt(handler.CreateSecurityLimit))
	r.Post("/otc-security-limits", handler.Adapt(handler.CreateSecurityLimitOtc))

	return &Router{
		mux:    r,
		logger: logger,
	}
}
