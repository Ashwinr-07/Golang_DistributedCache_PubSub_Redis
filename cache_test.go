package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

func TestCacheSetGet(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Error starting miniredis: %v", err)
	}
	defer mr.Close()

	cache, err := NewCache(mr.Addr())
	if err != nil {
		t.Fatalf("Error creating cache: %v", err)
	}

	ctx := context.Background()
	err = cache.Set(ctx, "testkey", []byte("testvalue"), 5*time.Minute)
	if err != nil {
		t.Fatalf("Error setting value: %v", err)
	}

	value, err := cache.Get(ctx, "testkey")
	if err != nil {
		t.Fatalf("Error getting value: %v", err)
	}

	if string(value) != "testvalue" {
		t.Errorf("Expected 'testvalue', got '%s'", string(value))
	}
}