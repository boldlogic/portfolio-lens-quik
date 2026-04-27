package router

import (
	"net/http"

	"github.com/boldlogic/packages/metrics"
	"github.com/boldlogic/quik-portfolio/pkg/transport/httpserver/handler"
	"github.com/boldlogic/quik-portfolio/pkg/transport/httpserver/httpmetrics"
	"github.com/boldlogic/quik-portfolio/pkg/transport/httpserver/middleware"
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
	r.Use(mw.Metrics)
	r.Use(mw.Wrap)

	r.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		handler.WriteResp(w, http.StatusOK, "ok")
	})

	return &Router{mux: r}
}

func (r *Router) Mount(pattern string, h http.Handler) {
	r.mux.Mount(pattern, h)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
