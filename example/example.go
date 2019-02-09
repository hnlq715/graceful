package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"
	"syscall"
	"time"

	"github.com/hnlq715/graceful"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/reflection"
)

type handler struct {
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, port: %v, %q", r.Host, html.EscapeString(r.URL.Path))
}

func main() {
	// graceful.ListenAndServe(":9222", &handler{})
	listenMultiAddrs()
}

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func listenMultiAddrs() {
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	reflection.Register(s)

	server := graceful.NewServer()
	// server.Register("0.0.0.0:9223", &handler{})
	server.RegisterGrpc("0.0.0.0:9224", s)
	server.RegisterHTTP("0.0.0.0:9225", &handler{})

	err := server.Run()
	fmt.Printf("error: %v\n", err)
}

func callReload() {
	server := graceful.NewServer()
	server.RegisterHTTP("0.0.0.0:9226", &handler{})
	go func() {
		time.Sleep(time.Second)
		server.Reload()
	}()

	err := server.Run()
	fmt.Printf("error: %v\n", err)
}

func setReloadSignal() {
	server := graceful.NewServer(
		graceful.WithReloadSignals([]syscall.Signal{syscall.SIGUSR2}),
		graceful.WithStopSignals([]syscall.Signal{syscall.SIGINT}),
		graceful.WithStopTimeout(time.Minute),
		graceful.WithWatchInterval(10*time.Second),
	)
	server.RegisterHTTP("0.0.0.0:9226", &handler{})
	server.Run()
}
