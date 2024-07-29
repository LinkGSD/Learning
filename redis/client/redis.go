package client

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

func init() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "192.168.182.132:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := Redis.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}
