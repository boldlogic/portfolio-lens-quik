package httpmetrics

import (
	"errors"
	"time"

	"github.com/boldlogic/packages/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type HTTPMetrics struct {
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
}

func NewMetrics(reg metrics.Registry) *HTTPMetrics {
	return &HTTPMetrics{
		httpRequestsTotal:   registerOrGetCounterVec(reg, newRequestsCounter()),
		httpRequestDuration: registerOrGetHistogramVec(reg, newRequestDuration()),
	}
}

func newRequestsCounter() *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route", "status"},
	)
}

func newRequestDuration() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: []float64{.001, .0025, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "route"},
	)
}

func registerOrGetCounterVec(reg prometheus.Registerer, c *prometheus.CounterVec) *prometheus.CounterVec {
	if err := reg.Register(c); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			return are.ExistingCollector.(*prometheus.CounterVec)
		}
		panic(err)
	}
	return c
}

func registerOrGetHistogramVec(reg prometheus.Registerer, h *prometheus.HistogramVec) *prometheus.HistogramVec {
	if err := reg.Register(h); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			return are.ExistingCollector.(*prometheus.HistogramVec)
		}
		panic(err)
	}
	return h
}

func (m *HTTPMetrics) RecordRequest(method, route, status string, duration time.Duration) {
	m.httpRequestsTotal.WithLabelValues(method, route, status).Inc()
	m.httpRequestDuration.WithLabelValues(method, route).Observe(duration.Seconds())
}
