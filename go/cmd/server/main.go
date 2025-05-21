package main

import (
	"context"
	"log"
	"net"
	"os"

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
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051" // default port
	}
	// connect to redis before anything
	redis_server.InitRedis()

	// start a tcp connection for to listen for gRPC
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	// Initialize a new gRPC server
	s := grpc.NewServer()
	reflection.Register(s)

	h := handler.NewGRPCHandler()
	limiter.RegisterRateLimiterServiceServer(s, &server{handler: h})

	log.Printf("gRPC server listening on :%s\n", port)

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
