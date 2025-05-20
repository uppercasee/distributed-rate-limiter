package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/uppercasee/drls/internal/handler"
	"github.com/uppercasee/drls/internal/redis"
	limiter "github.com/uppercasee/drls/pb"
)

type server struct {
	limiter.UnimplementedRateLimiterServiceServer
	handler *handler.GRPCHandler
}

func (s *server) Check(ctx context.Context, req *limiter.CheckRequest) (*limiter.CheckResponse, error) {
	return s.handler.HandleCheck(ctx, req)
}

func main() {
	// connect to redis before anything
	redis_server.InitRedis()

	// start a tcp connection for to listen for gRPC
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	// Initialize a new gRPC server
	s := grpc.NewServer()
	reflection.Register(s)

	h := handler.NewGRPCHandler()
	limiter.RegisterRateLimiterServiceServer(s, &server{handler: h})

	log.Println("gRPC server listening on :50051")

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
