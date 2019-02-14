package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

// Service implements grpc service
type Service struct {
	server *grpc.Server
}

// NewService creates grpc service
func NewService(opts ...grpc.ServerOption) *Service {
	return &Service{server: grpc.NewServer(opts...)}
}

// Serve serves grpc service
func (s *Service) Serve(net net.Listener) error {
	return s.server.Serve(net)
}

// GracefulStop stops grpc service graceful
func (s *Service) GracefulStop(ctx context.Context) error {
	s.server.GracefulStop()
	return nil
}

// Server returns grpc server
func (s *Service) Server() *grpc.Server {
	return s.server
}
