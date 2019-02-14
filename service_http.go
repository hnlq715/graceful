package graceful

import (
	"context"
	"net"
	"net/http"
)

// HTTPService implements http service
type HTTPService struct {
	server *http.Server
}

// NewHTTPService creates http service
func NewHTTPService() *HTTPService {
	return &HTTPService{server: &http.Server{}}
}

// Serve serves http service
func (s *HTTPService) Serve(net net.Listener) error {
	return s.server.Serve(net)
}

// GracefulStop stops http service graceful
func (s *HTTPService) GracefulStop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Server returns http server
func (s *HTTPService) Server() *http.Server {
	return s.server
}

// ListenAndServe starts server with (addr, handler)
func ListenAndServe(addr string, handler http.Handler) error {
	service := NewHTTPService()
	service.server.Handler = handler
	s := NewServer()
	s.Register(addr, service)

	return s.Run()
}
