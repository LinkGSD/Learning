package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func main() {
	limiter := rate.NewLimiter(rate.Every(time.Millisecond*31), 2)
	//time.Sleep(time.Second)
	for i := 0; i < 10; i++ {
		var ok bool
		if limiter.Allow() {
			ok = true
		}
		time.Sleep(time.Millisecond * 20)
		fmt.Println(ok, limiter.Burst())
	}
}
