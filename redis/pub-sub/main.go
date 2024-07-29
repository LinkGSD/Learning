package main

import (
	"Learning/redis/client"
	"context"
	"fmt"
	"time"
)

func main() {
	subscribe := client.Redis.Subscribe(context.Background(), "test")
	_, err := subscribe.Receive(context.Background())
	if err != nil {
		panic(err)
	}
	channel := subscribe.Channel()
	err = client.Redis.Publish(context.Background(), "test", "hello").Err()
	if err != nil {
		panic(err)
	}

	time.AfterFunc(time.Second, func() {
		// When pubsub is closed channel is closed too.
		_ = subscribe.Close()
	})

	// Consume messages.
	for msg := range channel {
		fmt.Println(msg.Channel, msg.Payload)
	}
}
