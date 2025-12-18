package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	jobs := make(chan int)

	results, scale := autoScalerPool[int, int](ctx, 3, 5, jobs)

	go func() {
		defer close(jobs)
		// send jobs
		for i := 0; i < 10; {
			select {
			case jobs <- i:
				i++
			case <-time.After(time.Millisecond * 200):
				select {
				case <-ctx.Done():
					return
				case scale <- 1:
					fmt.Println("scale up request sent...")
				}
			}

			//			time.Sleep(500 * time.Millisecond)
		}
	}()

	// read results
	for r := range results {
		fmt.Println("result:", r)
	}

	fmt.Println("all workers stopped cleanly")
}
