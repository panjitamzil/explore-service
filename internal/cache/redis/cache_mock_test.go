package redis

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type mockRedisClient struct {
	store map[string]string
}

func newMockRedisClient() *mockRedisClient {
	return &mockRedisClient{
		store: make(map[string]string),
	}
}

func (m *mockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx)
	if val, ok := m.store[key]; ok {
		cmd.SetVal(val)
		return cmd
	}
	cmd.SetErr(redis.Nil)
	return cmd
}

func (m *mockRedisClient) Set(ctx context.Context, key string, value interface{}, _ time.Duration) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx)
	m.store[key] = value.(string)
	cmd.SetVal("OK")
	return cmd
}

func (m *mockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)
	var removed int64
	for _, k := range keys {
		if _, ok := m.store[k]; ok {
			delete(m.store, k)
			removed++
		}
	}
	cmd.SetVal(removed)
	return cmd
}

func TestLikedCountCacheWithMock(t *testing.T) {
	client := newMockRedisClient()
	cache := NewLikedCountCache(client, "test:")
	ctx := context.Background()

	_, found, err := cache.GetCount(ctx, "u1")
	if err != nil {
		t.Fatalf("GetCount err: %v", err)
	}
	if found {
		t.Fatalf("expected not found")
	}

	if err := cache.SetCount(ctx, "u1", 5, 2*time.Second); err != nil {
		t.Fatalf("SetCount err: %v", err)
	}

	val, found, err := cache.GetCount(ctx, "u1")
	if err != nil {
		t.Fatalf("GetCount err: %v", err)
	}
	if !found || val != 5 {
		t.Fatalf("expected found count 5, got found=%v val=%d", found, val)
	}

	if err := cache.InvalidateCount(ctx, "u1"); err != nil {
		t.Fatalf("InvalidateCount err: %v", err)
	}
	if _, found, _ := cache.GetCount(ctx, "u1"); found {
		t.Fatalf("expected not found after invalidate")
	}
}
