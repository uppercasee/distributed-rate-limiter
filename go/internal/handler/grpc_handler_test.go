package handler_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uppercasee/drls/internal/handler"
	"github.com/uppercasee/drls/internal/redis"
	"github.com/uppercasee/drls/pb"
)

func TestHandleCheck(t *testing.T) {
	ctx := context.Background()
	clientID := "test-client-123"
	key := "rate_limit:" + clientID

	redis_server.InitRedis()
	redis_server.RDB.Del(ctx, key)

	h := handler.NewGRPCHandler()

	// Send requests within the limit (200)
	for range 200 {
		resp, err := h.HandleCheck(ctx, &pb.CheckRequest{ClientId: clientID})
		assert.NoError(t, err)
		assert.True(t, resp.Allowed)
		assert.Equal(t, int64(0), resp.RetryAfter)
	}

	// This 6th request should be rejected
	resp, err := h.HandleCheck(ctx, &pb.CheckRequest{ClientId: clientID})
	assert.NoError(t, err)
	assert.False(t, resp.Allowed)
	assert.Greater(t, resp.RetryAfter, int64(0))

	// rate limit window expires
	time.Sleep(time.Duration(resp.RetryAfter+1) * time.Second)

	resp, err = h.HandleCheck(ctx, &pb.CheckRequest{ClientId: clientID})
	assert.NoError(t, err)
	assert.True(t, resp.Allowed)

	// Cleanup
	redis_server.RDB.Del(ctx, key)
}
