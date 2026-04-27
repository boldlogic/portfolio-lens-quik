package middleware

import (
	"github.com/boldlogic/quik-portfolio/pkg/transport/httpserver/httpmetrics"
	"go.uber.org/zap"
)

type Middleware struct {
	logger  *zap.Logger
	metrics *httpmetrics.HTTPMetrics
}

func NewMiddleware(metrics *httpmetrics.HTTPMetrics, logger *zap.Logger) *Middleware {
	return &Middleware{
		logger:  logger,
		metrics: metrics,
	}
}
