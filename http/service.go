package http

import (
	"context"
	"net"
	"net/http"
)

type Service struct {
	server *http.Server
}

func NewService() *Service {
	return &Service{server: &http.Server{}}
}

func (s *Service) Serve(net net.Listener) error {
	return s.server.Serve(net)
}

func (s *Service) GracefulStop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Service) Server() *http.Server {
	return s.server
}
