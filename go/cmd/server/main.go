package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	limiter "github.com/uppercasee/drls/pb"
)

type server struct {
	limiter.UnimplementedRateLimiterServiceServer
}

func (s *server) Check(_ context.Context, req *limiter.CheckRequest) (*limiter.CheckResponse, error) {
	log.Println("Received Check request for client:", req.ClientId)

	return &limiter.CheckResponse{
		Allowed:    true,
		RetryAfter: 0,
	}, nil

}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	limiter.RegisterRateLimiterServiceServer(s, &server{})

	log.Println("gRPC server listening on :50051")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
