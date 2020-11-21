package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type wg struct {
	n int64
	c chan struct{}
}

func (w *wg) wait() {
	atomic.AddInt64(&w.n, 1)
	<-w.c
}

func (w *wg) done() {
	wait := atomic.LoadInt64(&w.n)
	fmt.Println(wait)
	for i := 0; int64(i) < wait; i++ {
		w.c <- struct{}{}
		atomic.AddInt64(&w.n, -1)
	}
}

type limiter struct {
	limit        int64
	currentLimit int64

	mu sync.Mutex
	wg *wg
}

func newLimiter(limit int64) *limiter {
	l := &limiter{
		limit:        limit,
		currentLimit: limit,
		wg:           &wg{c: make(chan struct{})},
	}

	go func() {
		c := time.Tick(time.Second)

		for range c {
			atomic.StoreInt64(&l.currentLimit, l.limit)

			fmt.Println("update rate")

			l.wg.done()
		}
	}()

	return l
}

func (l *limiter) decrease(n int64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.currentLimit < n {
		return false
	}

	fmt.Printf("before decrease %d\n", l.currentLimit)
	l.currentLimit -= n
	fmt.Printf("after decrease %d\n", l.currentLimit)

	return true
}

func (l *limiter) Take(n int64) {
	if ok := l.decrease(n); ok {
		return
	}

	l.wg.wait()
	l.Take(n)
}
