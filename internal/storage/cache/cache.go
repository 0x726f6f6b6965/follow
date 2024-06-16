package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type cache struct {
	client *redis.Client
}

// Delete implements Cache.
func (c *cache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get implements Cache.
func (c *cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrKeyNotFound
		}
		return "", err
	}
	return val, nil
}

// Set implements Cache.
func (c *cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func New(client *redis.Client) Cache {
	return &cache{client: client}
}
