package handler

import (
	"context"
	"log"

	"github.com/uppercasee/drls/pb"
)

type GRPCHandler struct{}

func NewGRPCHandler() *GRPCHandler {
	return &GRPCHandler{}
}

// HandleCheck processes a rate limit check request.
func (h *GRPCHandler) HandleCheck(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	log.Println("Handler received request for client:", req.ClientId)

	// Basic response, always allowed (MVP)
	return &pb.CheckResponse{
		Allowed:    true,
		RetryAfter: 0,
	}, nil
}
