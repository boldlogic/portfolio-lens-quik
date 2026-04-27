package clientmetrics

import (
	"time"

	"github.com/boldlogic/packages/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type ClientMetrics struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewMetrics(reg metrics.Registry) *ClientMetrics {
	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_client_requests_total",
			Help: "Total number of outgoing HTTP requests",
		},
		[]string{"method", "target", "endpoint", "status"},
	)
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_client_request_duration_seconds",
			Help:    "Outgoing HTTP request latency in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "target", "endpoint"},
	)
	reg.MustRegister(requestsTotal, requestDuration)
	return &ClientMetrics{
		requestsTotal:   requestsTotal,
		requestDuration: requestDuration,
	}
}

func (m *ClientMetrics) RecordRequest(method, target, endpoint, status string, duration time.Duration) {
	m.requestsTotal.WithLabelValues(method, target, endpoint, status).Inc()
	m.requestDuration.WithLabelValues(method, target, endpoint).Observe(duration.Seconds())
}
