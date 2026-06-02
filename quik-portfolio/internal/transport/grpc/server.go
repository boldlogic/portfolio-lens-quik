package grpc

import (
	"context"
	"net"

	"github.com/boldlogic/packages/metrics"
	quikv1 "github.com/boldlogic/portfolio-lens-quik/proto/gen/go/quik-portfolio/v1"
	"github.com/boldlogic/portfolio-lens-quik/quik-portfolio/internal/observability"
	grpcv1 "github.com/boldlogic/portfolio-lens-quik/quik-portfolio/internal/transport/grpc/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	server   *grpc.Server
	listener net.Listener
	logger   *zap.Logger
}

func NewServer(addr string, limits *grpcv1.Handler, logger *zap.Logger, reg metrics.Registry) (*Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	grpcMetrics := observability.NewGRPCMetrics(reg)
	s := grpc.NewServer(grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()))
	quikv1.RegisterLimitsServiceServer(s, limits)

	return &Server{
		server:   s,
		listener: lis,
		logger:   logger,
	}, nil
}

func (s *Server) Start() error {
	s.logger.Info("gRPC: слушаем", zap.String("addr", s.listener.Addr().String()))
	return s.server.Serve(s.listener)
}

func (s *Server) StopWithTimeout(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("gRPC: graceful stop завершен")
		return nil
	case <-ctx.Done():
		s.logger.Warn("gRPC: graceful stop timeout, принудительная остановка")
		s.server.Stop()
		return ctx.Err()
	}
}

func (s *Server) Addr() net.Addr {
	return s.listener.Addr()
}
