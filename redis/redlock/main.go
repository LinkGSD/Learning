package main

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

type RedMutex struct {
	client        redis.UniversalClient
	key           string
	value         string
	tries         int
	expiry        time.Duration
	until         time.Time
	timeoutFactor float64
}

func NewRedMutex(client redis.UniversalClient, key string) *RedMutex {
	return &RedMutex{
		client:        client,
		key:           key,
		tries:         3,
		expiry:        8 * time.Second,
		timeoutFactor: 0.01,
	}
}

func (r *RedMutex) Lock() error {
	return r.LockContext(context.Background())
}

func (r *RedMutex) LockContext(ctx context.Context) error {
	return r.lockContext(ctx)
}
func genValue() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (r *RedMutex) lockContext(ctx context.Context) error {
	value, _ := genValue()
	for i := 0; i < r.tries; i++ {
		start := time.Now()
		result, err := func() (bool, error) {
			ctx, cancel := context.WithTimeout(ctx, time.Duration(int64(float64(r.expiry)*r.timeoutFactor)))
			defer cancel()
			return r.client.SetNX(ctx, r.key, value, r.expiry).Result()
		}()
		// 不处理err, 避免出现redis设置成功, 但是响应超时
		now := time.Now()
		until := now.Add(r.expiry - now.Sub(start))
		if result && now.Before(until) {
			r.value = value
			r.until = until
			return nil
		}
		_, _ = func() (bool, error) {
			ctx, cancel := context.WithTimeout(ctx, time.Duration(int64(float64(r.expiry)*r.timeoutFactor)))
			defer cancel()
			return r.release(ctx)
		}()
		if i == r.tries-1 && err != nil {
			return err
		}
	}
	return nil
}

//go:embed lua/unlock.lua
var deleteScript string

func (r *RedMutex) release(ctx context.Context) (bool, error) {
	script := redis.NewScript(deleteScript)
	result, err := script.Eval(ctx, r.client, []string{r.key}, r.value).Result()
	if err != nil {
		return false, err
	}
	if result == int64(-1) {
		return false, errors.New("failed to unlock, lock was already expired")
	}
	return result != int64(0), nil
}

func main() {
	client := redis.NewClient(&redis.Options{Addr: "192.168.182.132:6379"})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	//result, err := client.SetNX(context.Background(), "lock", "test", time.Hour).Result()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(result)
	bytes, _ := os.ReadFile("./lock/lua/unlock.lua")
	script := redis.NewScript(string(bytes))
	result, err := script.Eval(context.Background(), client, []string{"lock"}, "test").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: make([]string, 0),
	})
	mutex := NewRedMutex(clusterClient, "lock")
	err = mutex.Lock()
	if err != nil {
		fmt.Println(err)
	}
}
