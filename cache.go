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

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mutex.RLock()
	if value, ok := c.localCache[key]; ok {
		c.mutex.RUnlock()
		return value, nil
	}
	c.mutex.RUnlock()

	value, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	c.mutex.Lock()
	c.localCache[key] = value
	c.mutex.Unlock()

	return value, nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return err
	}

	c.mutex.Lock()
	delete(c.localCache, key)
	c.mutex.Unlock()

	return nil
}