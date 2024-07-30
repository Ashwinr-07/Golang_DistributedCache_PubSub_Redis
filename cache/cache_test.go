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


func TestDistributedCacheUpdates(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Error starting miniredis: %v", err)
	}
	defer mr.Close()

	cache1, err := NewCache(mr.Addr())
	if err != nil {
		t.Fatalf("Error creating cache1: %v", err)
	}

	cache2, err := NewCache(mr.Addr())
	if err != nil {
		t.Fatalf("Error creating cache2: %v", err)
	}

	ctx := context.Background()

	// Set a value in cache1
	err = cache1.Set(ctx, "sharedkey", []byte("sharedvalue"), 5*time.Minute)
	if err != nil {
		t.Fatalf("Error setting value in cache1: %v", err)
	}

	// Allow some time for the update to propagate
	time.Sleep(100 * time.Millisecond)

	// Try to get the value from cache2
	value, err := cache2.Get(ctx, "sharedkey")
	if err != nil {
		t.Fatalf("Error getting value from cache2: %v", err)
	}

	if string(value) != "sharedvalue" {
		t.Errorf("Expected 'sharedvalue', got '%s'", string(value))
	}

	// Delete the value from cache1
	err = cache1.Delete(ctx, "sharedkey")
	if err != nil {
		t.Fatalf("Error deleting value from cache1: %v", err)
	}

	// Allow some time for the update to propagate
	time.Sleep(100 * time.Millisecond)

	// Try to get the deleted value from cache2
	_, err = cache2.Get(ctx, "sharedkey")
	if err == nil {
		t.Errorf("Expected error for deleted key, got nil")
	}
}