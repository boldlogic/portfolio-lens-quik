package writeserver

import (
	"github.com/boldlogic/packages/metrics"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/router"
	v1 "github.com/boldlogic/portfolio-lens-quik/quik-portfolio-writer/internal/transport/http/v1"
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
