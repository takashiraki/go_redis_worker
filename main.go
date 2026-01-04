package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	print("Hello, Go Worker!\n")

	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
	})

	fmt.Print("Setting key 'foo' to 'bar'...\n")
	for {
		res, err := rdb.BLPop(ctx, 0, "order.queue").Result()
		if err != nil {
			fmt.Printf("Error popping from queue: %v\n", err)
			continue
		}
		fmt.Printf("Popped value: %s\n", res[1])
	}
}
