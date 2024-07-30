package cache

import (
	"sync"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	client        *redis.Client
	localCache    map[string][]byte
	mutex         sync.RWMutex
	subscribeConn *redis.PubSub
}


func NewCache(addr string) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	cache := &Cache{
		client:     client,
		localCache: make(map[string][]byte),
	}

	if err := cache.subscribeToUpdates(); err != nil {
		return nil, err
	}

	return cache, nil
}


func (c *Cache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	if err := c.client.Set(ctx, key, value, expiration).Err(); err != nil {
		return err
	}

	update, err := json.Marshal(struct {
		Key   string `json:"key"`
		Value []byte `json:"value"`
	}{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return err
	}

	if err := c.client.Publish(ctx, "cache_updates", update).Err(); err != nil {
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

	update, err := json.Marshal(struct {
		Key   string `json:"key"`
		Value []byte `json:"value"`
	}{
		Key:   key,
		Value: nil,
	})
	if err != nil {
		return err
	}

	if err := c.client.Publish(ctx, "cache_updates", update).Err(); err != nil {
		return err
	}

	c.mutex.Lock()
	delete(c.localCache, key)
	c.mutex.Unlock()

	return nil
}

func (c *Cache) subscribeToUpdates() error {
	pubsub := c.client.Subscribe(context.Background(), "cache_updates")
	c.subscribeConn = pubsub

	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			var update struct {
				Key   string `json:"key"`
				Value []byte `json:"value"`
			}
			if err := json.Unmarshal([]byte(msg.Payload), &update); err != nil {
				fmt.Printf("Error unmarshaling update: %v\n", err)
				continue
			}
			c.mutex.Lock()
			if update.Value == nil {
				delete(c.localCache, update.Key)
			} else {
				c.localCache[update.Key] = update.Value
			}
			c.mutex.Unlock()
		}
	}()

	return nil
}