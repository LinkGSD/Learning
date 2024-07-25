package redis

import (
	"github.com/go-redis/redis"
	"time"
)

var Cli *redis.Client

func init() {
	Cli = redis.NewClient(&redis.Options{
		Addr:     "192.168.182.132:6379",
		DB:       0,
		Password: "",
	})

	_, err := Cli.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func IsConsumed(key string) bool {
	result, err := Cli.SetNX(key, "", 2*time.Hour).Result()
	if err != nil {
		return false
	}
	return result
}
