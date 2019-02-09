package graceful

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

// constants
const (
	EnvWorker       = "GRACEFUL_WORKER"
	EnvNumFD        = "GRACEFUL_NUMFD"
	EnvOldWorkerPid = "GRACEFUL_OLD_WORKER_PID"
	ValWorker       = "1"
)

var (
	defaultWatchInterval = time.Second
	defaultStopTimeout   = 20 * time.Second
	defaultReloadSignals = []syscall.Signal{syscall.SIGHUP, syscall.SIGUSR1}
	defaultStopSignals   = []syscall.Signal{syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT}

	StartedAt time.Time
)

type option struct {
	reloadSignals []syscall.Signal
	stopSignals   []syscall.Signal
	watchInterval time.Duration
	stopTimeout   time.Duration
}

type Option func(o *option)

// WithReloadSignals set reload signals, otherwise, default ones are used
func WithReloadSignals(sigs []syscall.Signal) Option {
	return func(o *option) {
		o.reloadSignals = sigs
	}
}

// WithStopSignals set stop signals, otherwise, default ones are used
func WithStopSignals(sigs []syscall.Signal) Option {
	return func(o *option) {
		o.stopSignals = sigs
	}
}

// WithStopTimeout set stop timeout for graceful shutdown
//  if timeout occurs, running connections will be discard violently.
func WithStopTimeout(timeout time.Duration) Option {
	return func(o *option) {
		o.stopTimeout = timeout
	}
}

// WithWatchInterval set watch interval for worker checking master process state
func WithWatchInterval(timeout time.Duration) Option {
	return func(o *option) {
		o.watchInterval = timeout
	}
}

type Server struct {
	opt   *option
	addrs []string
	grpc  map[string]*grpc.Server
	http  map[string]*http.Server
}

func NewServer(opts ...Option) *Server {
	option := &option{
		reloadSignals: defaultReloadSignals,
		stopSignals:   defaultStopSignals,
		watchInterval: defaultWatchInterval,
		stopTimeout:   defaultStopTimeout,
	}
	for _, opt := range opts {
		opt(option)
	}
	return &Server{
		addrs: make([]string, 0),
		opt:   option,
		grpc:  make(map[string]*grpc.Server),
		http:  make(map[string]*http.Server),
	}
}

// RegisterHTTP registers an addr and its corresponding handler
// all (addr, handler) pair will be started with server.Run
func (s *Server) RegisterHTTP(addr string, handler http.Handler) {
	s.addrs = append(s.addrs, addr)
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Println("invalid address:", addr)
		return
	}
	s.http[port] = &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}

// RegisterGrpc register grpc server
func (s *Server) RegisterGrpc(addr string, server *grpc.Server) {
	s.addrs = append(s.addrs, addr)
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Println("invalid address:", addr)
		return
	}
	s.grpc[port] = server
}

// Run runs all register servers
func (s *Server) Run() error {
	if len(s.addrs) == 0 {
		return ErrNoServers
	}
	StartedAt = time.Now()
	if IsWorker() {
		worker := &worker{grpc: s.grpc, http: s.http, opt: s.opt, stopCh: make(chan struct{})}
		return worker.run()
	}
	master := &master{addrs: s.addrs, opt: s.opt, workerExit: make(chan error)}
	return master.run()
}

// Reload reload server gracefully
func (s *Server) Reload() error {
	ppid := os.Getppid()
	if IsWorker() && ppid != 1 && len(s.opt.reloadSignals) > 0 {
		return syscall.Kill(ppid, s.opt.reloadSignals[0])
	}

	// Reload called by user from outside usally in user's handler
	// which works on worker, master don't need to handle this
	return nil
}

// ListenAndServe starts server with (addr, handler)
func ListenAndServe(addr string, handler http.Handler) error {
	server := NewServer()
	server.RegisterHTTP(addr, handler)
	return server.Run()
}

func IsWorker() bool {
	return os.Getenv(EnvWorker) == ValWorker
}

func IsMaster() bool {
	return !IsWorker()
}
