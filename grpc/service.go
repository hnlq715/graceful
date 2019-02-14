package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

type Service struct {
	server *grpc.Server
}

func NewService(opts ...grpc.ServerOption) *Service {
	return &Service{server: grpc.NewServer(opts...)}
}

func (s *Service) Serve(net net.Listener) error {
	return s.server.Serve(net)
}

func (s *Service) GracefulStop(ctx context.Context) error {
	s.server.GracefulStop()
	return nil
}

func (s *Service) Server() *grpc.Server {
	return s.server
}
