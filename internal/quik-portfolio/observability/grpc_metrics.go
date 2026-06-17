package observability

import (
	"context"
	"strings"
	"time"

	"github.com/boldlogic/packages/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	grpcstatus "google.golang.org/grpc/status"
)

const unknownGRPCMethod = "unknown"

type GRPCMetrics struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewGRPCMetrics(reg metrics.Registry) *GRPCMetrics {
	if reg == nil {
		return &GRPCMetrics{}
	}

	return &GRPCMetrics{
		requestsTotal:   registerOrGetCounterVec(reg, newGRPCRequestsTotal()),
		requestDuration: registerOrGetHistogramVec(reg, newGRPCRequestDuration()),
	}
}

func (m *GRPCMetrics) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if m == nil || m.requestsTotal == nil || m.requestDuration == nil {
			return handler(ctx, req)
		}

		service, method := splitGRPCFullMethod(info.FullMethod)
		start := time.Now()

		resp, err := handler(ctx, req)
		code := grpcstatus.Code(err).String()

		m.requestsTotal.WithLabelValues(service, method, code).Inc()
		m.requestDuration.WithLabelValues(service, method, code).Observe(time.Since(start).Seconds())

		return resp, err
	}
}

func newGRPCRequestsTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "quik_portfolio_grpc_requests_total",
		Help: "Total number of finished gRPC server requests.",
	}, []string{"grpc_service", "grpc_method", "grpc_code"})
}

func newGRPCRequestDuration() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "quik_portfolio_grpc_request_duration_seconds",
		Help:    "gRPC server request duration in seconds.",
		Buckets: durationBuckets,
	}, []string{"grpc_service", "grpc_method", "grpc_code"})
}

func splitGRPCFullMethod(fullMethod string) (string, string) {
	trimmed := strings.TrimPrefix(fullMethod, "/")
	idx := strings.LastIndex(trimmed, "/")
	if idx < 0 {
		return trimmed, unknownGRPCMethod
	}
	return trimmed[:idx], trimmed[idx+1:]
}
