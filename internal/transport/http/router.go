package portfolioserver

import (
	"github.com/boldlogic/packages/metrics"
	v1 "github.com/boldlogic/quik-portfolio/internal/transport/http/v1"
	"github.com/boldlogic/quik-portfolio/pkg/transport/httpserver/router"
	"go.uber.org/zap"
)

type Router struct {
	*router.Router
}

func NewRouter(handler *v1.Handler, log *zap.Logger, reg metrics.Registry) *Router {
	base := router.NewRouter(log, reg)
	base.Mount("/api/v1", v1.NewRouter(handler, log))
	return &Router{Router: base}
}
