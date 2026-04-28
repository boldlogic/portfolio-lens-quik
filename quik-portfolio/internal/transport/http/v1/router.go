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
			r.Get("/", handler.Adapt(handler.GetLimits))
			r.Get("/money", handler.Adapt(handler.GetMoneyLimits))
			r.Post("/money", handler.Adapt(handler.CreateMoneyLimit))

			r.Get("/securities", handler.Adapt(handler.GetSecurityLimits))
			r.Post("/securities", handler.Adapt(handler.CreateSecurityLimit))
			r.Get("/securities/otc", handler.Adapt(handler.GetSecurityLimitsOtc))
			r.Post("/securities/otc", handler.Adapt(handler.CreateSecurityLimitOtc))
		})
		r.Get("/portfolio", handler.Adapt(handler.GetPortfolio))
		r.Route("/firms", func(r chi.Router) {
			r.Get("/", handler.Adapt(handler.GetFirms))
			r.Post("/", handler.Adapt(handler.CreateFirm))
			r.Get("/{id}", handler.Adapt(handler.GetFirm))
			r.Patch("/{id}", handler.Adapt(handler.UpdateFirm))
		})
	})
	return &Router{
		mux:    r,
		logger: logger,
	}
}
