package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/hnlq715/graceful"
	"github.com/hnlq715/graceful/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/reflection"
)

type handler struct {
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, port: %v, %q\n", r.Host, html.EscapeString(r.URL.Path))
}

func main() {
	listenAndServe()
	// listenMultiAddrs()
}

func listenAndServe() {
	graceful.ListenAndServe("0.0.0.0:9223", &handler{})
}

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func listenMultiAddrs() {
	gs := grpc.NewService()
	pb.RegisterGreeterServer(gs.Server(), &server{})
	reflection.Register(gs.Server())

	hs := graceful.NewHTTPService()
	hs.Server().Handler = &handler{}

	server := graceful.NewServer()
	server.Register("0.0.0.0:9224", gs)
	server.Register("0.0.0.0:9225", hs)

	err := server.Run()
	fmt.Printf("error: %v\n", err)
}
