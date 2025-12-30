package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// 1 tight CPU loop in its own goroutine with context-based exit
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("tight loop exiting:", ctx.Err())
				return
			default:
				// tight CPU work
				// fmt.Println("doing work...")
				_ = 1 + 1
			}
		}
	}(ctx)

	// many goroutines doing time.Sleep with context-based exit
	for i := 0; i < 10000; i++ {
		go func(id int, ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("sleeper exiting", id) // optional log
					return
				case <-time.After(1 * time.Second):
					// continue sleeping loop
				}
			}
		}(i, ctx)
	}

	// stop after some time (example: 10s)
	time.AfterFunc(10*time.Second, func() {
		cancel()
	})

	// keep main alive until context canceled
	<-ctx.Done()
	time.Sleep(200 * time.Millisecond) // allow goroutines to print/cleanup
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0), "NumGoroutine:", runtime.NumGoroutine())
}
