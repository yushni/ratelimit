package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	l := NewLimiter(ctx, 5)

	s := time.Now()

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		go func(i int) {
			wg.Add(1)
			l.Do(1, func() error {
				fmt.Println(i)
				return nil
			})

			wg.Done()
		}(i)
		time.Sleep(time.Millisecond)
	}

	wg.Wait()
	cancel()

	fmt.Println(time.Since(s))
}
