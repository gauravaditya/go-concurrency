package main

import (
	"context"
	"time"
)

func producer(num int) <-chan int {
	out := make(chan int)

	for i := range 5 {
		out <- (i + 1) * num
	}

	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(10*time.Second, cancel)
	pool := NewWorkerPool(
		ctx,
		20,
		withQSize(5),
		withResultQSize(5),
		withRateLimit(2),
	)

	go func() {
		defer pool.Q.shutdown()
		for i := 0; i < 500; i++ {
			pool.Q.push(ctx, i)
		}
	}()

	for v := range pool.results.jobs {
		println(v)
	}
}
