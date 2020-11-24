package main

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type wg chan int64

func (w wg) wait(n int64) {
	w <- n
}

func (w wg) done(limit int64) {
	var current int64

	for c := range w {
		current += c
		if current > limit {
			return
		}
	}
}

type limiter struct {
	limit        int64
	currentLimit int64

	mu sync.Mutex
	wg wg
}

func NewLimiter(ctx context.Context, limit int64) *limiter {
	l := &limiter{
		limit:        limit,
		currentLimit: limit,
		wg:           make(chan int64),
	}

	go func() {
		t := time.NewTicker(time.Second)

		for {
			select {
			case <-ctx.Done():
				t.Stop()
				return
			case <-t.C:
				atomic.StoreInt64(&l.currentLimit, l.limit)
				l.wg.done(l.limit)
			}
		}
	}()

	return l
}

func (w *limiter) decrease(n int64) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentLimit < n {
		return false
	}

	w.currentLimit -= n
	return true
}

func (w *limiter) Do(n int64, f func() error) {
	if ok := w.decrease(n); ok {
		_ = f()
		return
	}

	w.wg.wait(n)
	w.Do(n, f)
}
