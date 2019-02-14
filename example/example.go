package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/hnlq715/graceful"
	grpcs "github.com/hnlq715/graceful/grpc"
	https "github.com/hnlq715/graceful/http"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/reflection"
)

type handler struct {
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, port: %v, %q\n", r.Host, html.EscapeString(r.URL.Path))
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
	gs := grpcs.NewService()
	pb.RegisterGreeterServer(gs.Server(), &server{})
	reflection.Register(gs.Server())

	hs := https.NewService()
	hs.Server().Handler = &handler{}

	server := graceful.NewServer()
	server.Register("0.0.0.0:9224", gs)
	server.Register("0.0.0.0:9225", hs)

	err := server.Run()
	fmt.Printf("error: %v\n", err)
}
