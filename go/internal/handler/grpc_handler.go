package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/uppercasee/drls/internal/redis"
	"github.com/uppercasee/drls/pb"
)

type GRPCHandler struct{}

func NewGRPCHandler() *GRPCHandler {
	return &GRPCHandler{}
}

func (h *GRPCHandler) HandleCheck(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	limit := 200
	windowSize := int64(60) * int64(time.Second)
	ttl := int64(210) // window + 10s
	now := time.Now().UnixNano()
	threshold := now - windowSize

	key := fmt.Sprintf("rate_limit:%s", req.ClientId)
	log.Println("Handler received request for client:", req.ClientId)
	log.Println("Redis key:", key)

	// -- KEYS[1] = rate limit key
	// -- ARGV[1] = current timestamp (in nanoseconds)
	// -- ARGV[2] = threshold timestamp (oldest valid timestamp)
	// -- ARGV[3] = rate limit count
	// -- ARGV[4] = TTL for the key (in seconds)
	script := `
	redis.call('ZREMRANGEBYSCORE', KEYS[1], '-inf', ARGV[2])
	local count = redis.call('ZCARD', KEYS[1])
	if tonumber(count) < tonumber(ARGV[3]) then
		redis.call('ZADD', KEYS[1], ARGV[1], ARGV[1])
		redis.call('EXPIRE', KEYS[1], tonumber(ARGV[4]))
		return {1, 0}
	else
		local oldest = redis.call('ZRANGE', KEYS[1], 0, 0)[1]
		return {0, oldest}
	end
	`

	res, err := redis_server.RDB.Eval(ctx, script, []string{key},
		now,
		threshold,
		limit,
		ttl,
	).Result()

	if err != nil {
		log.Printf("Lua script execution error: %v", err)
		return nil, err
	}

	result, ok := res.([]interface{})
	if !ok || len(result) != 2 {
		log.Printf("Unexpected Lua script result: %+v", res)
		return nil, fmt.Errorf("invalid script result")
	}

	allowed, _ := result[0].(int64)

	if allowed == 1 {
		return &pb.CheckResponse{
			Allowed:    true,
			RetryAfter: 0,
		}, nil
	}

	oldestStr, _ := result[1].(string)
	oldestTs, _ := strconv.ParseInt(oldestStr, 10, 64)
	retryAfter := max(int32((windowSize-(now-oldestTs))/int64(time.Second)), 0)

	return &pb.CheckResponse{
		Allowed:    false,
		RetryAfter: int64(retryAfter),
	}, nil
}

func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
