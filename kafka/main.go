package main

import (
	"Learning/kafka/consumer"
	"Learning/kafka/model"
	"Learning/kafka/producer"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	for i := 0; i < 30; i++ {
		user := &model.User{
			Id:       int64(i + 1),
			UserName: fmt.Sprintf("user:%d", i),
			Age:      220,
		}
		producer.SendMessage(ctx, user)
	}
	//producer.Producer.Close()
	//go consumer.ReadMessage(ctx)
	//listenSignal()
}

func listenSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	fmt.Printf("收到信号 %s ", sig.String())
	if consumer.Consumer != nil {
		consumer.Consumer.Close()
	}
	os.Exit(0)
}
