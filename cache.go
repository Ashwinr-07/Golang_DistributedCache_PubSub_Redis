package cache

import (
	"sync"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	client     *redis.Client
	localCache map[string][]byte
	mutex      sync.RWMutex
}

func NewCache(addr string) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &Cache{
		client:     client,
		localCache: make(map[string][]byte),
	}, nil
}

func (c *Cache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	if err := c.client.Set(ctx, key, value, expiration).Err(); err != nil {
		return err
	}

	c.mutex.Lock()
	c.localCache[key] = value
	c.mutex.Unlock()

	return nil
}