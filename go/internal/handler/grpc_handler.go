package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/uppercasee/drls/internal/redis"
	"github.com/uppercasee/drls/pb"
)

type GRPCHandler struct{}

func NewGRPCHandler() *GRPCHandler {
	return &GRPCHandler{}
}

// HandleCheck processes a rate limit check request.
func (h *GRPCHandler) HandleCheck(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	limit := 5

	log.Println("Handler received request for client:", req.ClientId)
	// create redis key
	key := fmt.Sprintf("rate_limit:%s", req.ClientId)
	log.Println("key: ", key)

	// remove outdated entries older than the window
	now := time.Now().UnixNano()
	windowSize := int64(60) * int64(time.Second)
	threshold := now - windowSize

	removed, err := redis_server.RDB.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%v", threshold)).Result()

	log.Println("number of removed entries: ", removed)
	if err != nil {
		log.Printf("failed to remove outdated entries: %v", err)
	}

	// count entries in the window
	count, err := redis_server.RDB.ZCard(ctx, key).Result()
	log.Println("number of entries: ", count)
	if err != nil {
		log.Printf("failed to count the entries: %v", err)
	}

	// decision
	if count < int64(limit) {
		_, err := redis_server.RDB.ZAdd(ctx, key, redis.Z{
			Score:  float64(now),
			Member: now,
		}).Result()

		if err != nil {
			log.Println("ZAdd error:", err)
		}

		// add expiry of limit + 10seconds
		ttl := time.Duration(windowSize+10) * time.Second
		err = redis_server.RDB.Expire(ctx, key, ttl).Err()
		if err != nil {
			log.Printf("failed to set expiry on redis key: %v", err)
		}

		return &pb.CheckResponse{
			Allowed:    true,
			RetryAfter: 0,
		}, nil
	}

	// get the oldest entry
	oldest, err := redis_server.RDB.ZRange(ctx, key, 0, 0).Result()

	log.Println("the oldest entry: ", oldest)
	if err != nil {
		log.Printf("failed to get the oldest entry: %v", err)
	}

	// calculate retry after
	retryAfter := int32(0)
	if len(oldest) > 0 {
		oldestTs, _ := strconv.ParseInt(oldest[0], 10, 64)
		retryAfter = max(int32((windowSize-(now-oldestTs))/int64(time.Second)), 0)
	}

	// return false with retryAfter in seconds
	return &pb.CheckResponse{
		Allowed:    false,
		RetryAfter: int64(retryAfter),
	}, nil
}
