package referencehttp

import (
	"github.com/boldlogic/packages/metrics"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/router"
	v1 "github.com/boldlogic/portfolio-lens-quik/internal/quik-reference-data/api/http/v1"
	"go.uber.org/zap"
)

type Router struct {
	*router.Router
}

func NewRouter(handler *v1.Handler, logger *zap.Logger, reg metrics.Registry) *Router {
	base := router.NewRouter(logger, reg)
	base.Mount("/api/v1", v1.NewRouter(handler, logger))
	return &Router{Router: base}
}
