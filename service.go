package graceful

import (
	"context"
	"net"
)

// Service interface
type Service interface {
	Serve(net.Listener) error
	GracefulStop(context.Context) error
}
