package router

import (
	"net/http"

	"github.com/boldlogic/packages/metrics"
	"github.com/boldlogic/packages/transport/httpserver/httpmetrics"
	"github.com/boldlogic/packages/transport/httpserver/middleware"
	"github.com/boldlogic/packages/transport/httpserver/response"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Router struct {
	mux *chi.Mux
}

func NewRouter(logger *zap.Logger, reg metrics.Registry) *Router {
	r := chi.NewRouter()

	m := httpmetrics.NewMetrics(reg)
	mw := middleware.NewMiddleware(m, logger)
	r.Use(mw.Recover)
	r.Use(mw.Metrics)
	r.Use(mw.Wrap)

	r.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		response.WriteResp(w, http.StatusOK, "ok")
	})

	return &Router{mux: r}
}

func (r *Router) Mount(pattern string, h http.Handler) {
	r.mux.Mount(pattern, h)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
