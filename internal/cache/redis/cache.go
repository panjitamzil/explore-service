package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type LikedCountCache struct {
	client *redis.Client
	prefix string
}

func NewLikedCountCache(client *redis.Client, prefix string) *LikedCountCache {
	if prefix == "" {
		prefix = "liked_count:"
	}

	return &LikedCountCache{
		client: client,
		prefix: prefix,
	}
}

func (c *LikedCountCache) key(recipientID string) string {
	return c.prefix + recipientID
}

func (c *LikedCountCache) GetCount(ctx context.Context, recipientID string) (count uint64, found bool, err error) {
	key := c.key(recipientID)

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("redis get liked_count: %w", err)
	}

	n, parseErr := strconv.ParseUint(val, 10, 64)
	if parseErr != nil {
		return 0, false, fmt.Errorf("invalid liked_count value in cache: %w", parseErr)
	}

	return n, true, nil
}

func (c *LikedCountCache) SetCount(ctx context.Context, recipientID string, count uint64, ttl time.Duration) error {
	key := c.key(recipientID)

	if err := c.client.Set(ctx, key, strconv.FormatUint(count, 10), ttl).Err(); err != nil {
		return fmt.Errorf("redis set liked_count: %w", err)
	}

	return nil
}

func (c *LikedCountCache) InvalidateCount(ctx context.Context, recipientID string) error {
	key := c.key(recipientID)

	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis del liked_count: %w", err)
	}

	return nil
}
