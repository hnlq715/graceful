package http

import (
	"context"
	"net"
	"net/http"
)

// Service implements http service
type Service struct {
	server *http.Server
}

// NewService creates http service
func NewService() *Service {
	return &Service{server: &http.Server{}}
}

// Serve serves http service
func (s *Service) Serve(net net.Listener) error {
	return s.server.Serve(net)
}

// GracefulStop stops http service graceful
func (s *Service) GracefulStop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Server returns http server
func (s *Service) Server() *http.Server {
	return s.server
}
