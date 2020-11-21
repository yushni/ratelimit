package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	l := newLimiter(100)

	s := time.Now()

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		go func() {
			wg.Add(1)
			l.Take(50)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(time.Since(s))
}
