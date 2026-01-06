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
		Addr: "redis:6379",
	})

	fmt.Print("Setting key 'foo' to 'bar'...\n")
	for {
		// res, err := rdb.BLPop(ctx, 0, "order.queue").Result()
		// if err != nil {
		// 	fmt.Printf("Error popping from queue: %v\n", err)
		// 	continue
		// }
		// fmt.Printf("Popped value: %s\n", res[1])

		streams, err := rdb.XReadGroup(
			ctx,
			&redis.XReadGroupArgs{
				Group:    "mail-group",
				Consumer: "worker-1",
				Streams:  []string{"order.events", ">"},
				Count:    1,
				Block:    5,
			},
		).Result()

		if err != nil && err != redis.Nil {
			fmt.Printf("Error reading from stream: %v\n", err)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				handleEvent(msg.Values)

				rdb.XAck(ctx, "order.events", "mail-group", msg.ID)
			}
		}
	}
}

func handleEvent(values map[string]interface{}) {
	fmt.Println("Received event: ")

	for k, v := range values {
		fmt.Printf("  %s: %v\n", k, v)
	}
}
