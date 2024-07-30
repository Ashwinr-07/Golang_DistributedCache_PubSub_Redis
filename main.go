package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Ashwinr-07/Golang_DistributedCache_PubSub_Redis/cache"
)

func main() {
	ctx := context.Background()

	// Create two cache instances
	cache1, err := cache.NewCache("localhost:6379")
	if err != nil {
		panic(err)
	}
	defer cache1.Close()

	cache2, err := cache.NewCache("localhost:6379")
	if err != nil {
		panic(err)
	}
	defer cache2.Close()

	// Set a value in cache1
	err = cache1.Set(ctx, "examplekey", []byte("examplevalue"), 5*time.Minute)
	if err != nil {
		panic(err)
	}
	fmt.Println("Value set in cache1")

	// Allow some time for the update to propagate
	time.Sleep(100 * time.Millisecond)

	// Get the value from cache2
	value, err := cache2.Get(ctx, "examplekey")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Retrieved value from cache2: %s\n", string(value))

	// Delete the key from cache1
	err = cache1.Delete(ctx, "examplekey")
	if err != nil {
		panic(err)
	}
	fmt.Println("Key deleted from cache1")

	// Allow some time for the update to propagate
	time.Sleep(100 * time.Millisecond)

	// Try to get the deleted key from cache2
	_, err = cache2.Get(ctx, "examplekey")
	if err != nil {
		fmt.Printf("Error getting deleted key from cache2: %v\n", err)
	}
}