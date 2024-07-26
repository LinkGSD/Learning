package main

import (
	"fmt"
	"sync/atomic"
	"time"
	_ "unsafe"
)

type Sema struct {
	sema  uint32
	count atomic.Int32
}

func (s *Sema) acquire() {
	s.count.Add(1)
	acquire(&s.sema)
}
func (s *Sema) release() {
	s.count.Add(-1)
	release(&s.sema, false, 0)
}

func main() {
	sema := &Sema{}
	go func() {
		sema.acquire()
		fmt.Println("t1")
	}()
	//time.Sleep(time.Millisecond)
	go func() {
		sema.acquire()
		fmt.Println("t2")
	}()
	time.Sleep(time.Millisecond)
	sema.release()
	time.Sleep(time.Millisecond)
	sema.release()
	time.Sleep(time.Millisecond)
	fmt.Println("finish")
}

//go:linkname acquire sync.runtime_Semacquire
func acquire(s *uint32)

//go:linkname release sync.runtime_Semrelease
func release(s *uint32, handoff bool, skipframes int)
