package main

import (
	"fmt"
	"sync"
)

func main() {
}

// wg.state 高位表示当前任务数量,地位表示当前等待的协程数量
func t1(num int) {
	wg := &sync.WaitGroup{}
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("working")
		}()
	}
	wg.Wait()
}

func t2(num, Concurrency int) {
	wg := &sync.WaitGroup{}
	c := make(chan struct{}, Concurrency)
	defer close(c)
	wg.Add(num)
	for i := 0; i < num; i++ {
		c <- struct{}{}
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
			<-c
		}(i)
	}
	wg.Wait()
}
