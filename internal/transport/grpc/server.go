package grpc

import (
	"net"

	grpcv1 "github.com/boldlogic/quik-portfolio/internal/transport/grpc/v1"
	quikv1 "github.com/boldlogic/quik-portfolio/proto/gen/go/quik/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	server   *grpc.Server
	listener net.Listener
	logger   *zap.Logger
}

func NewServer(addr string, limits *grpcv1.Handler, logger *zap.Logger) (*Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()
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

func (s *Server) Stop() {
	s.logger.Info("gRPC: остановка")
	s.server.GracefulStop()
}

func (s *Server) Addr() net.Addr {
	return s.listener.Addr()
}
