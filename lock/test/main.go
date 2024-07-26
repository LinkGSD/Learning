package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func ttt(tries int) time.Duration { return time.Duration(rand.Intn(3)+1) * time.Second }
func main() {
	//a := 8
	//b := 0.05
	//ctx := context.Background()
	//var timer *time.Timer
	//for i := 0; i < 10; i++ {
	//	if i != 0 {
	//		duration := ttt(i)
	//		fmt.Println(duration)
	//		if timer == nil {
	//			timer = time.NewTimer(duration)
	//		} else {
	//			timer.Reset(duration)
	//		}
	//
	//		select {
	//		case <-ctx.Done():
	//			fmt.Println(333)
	//		case <-timer.C:
	//			fmt.Println(111)
	//		}
	//	}
	//	func() {
	//		_, cancel := context.WithTimeout(ctx, time.Duration(int64(float64(a)*b)))
	//		defer cancel()
	//	}()
	//	fmt.Println(222)
	//}
	TestTimeout()
}

func TestTimeout() {
	ctx := context.Background()
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, time.Second)
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("timeout")
			cancel()
		}
	}()
	time.Sleep(time.Second * 2)
}
