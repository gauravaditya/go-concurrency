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
	ticker := time.NewTicker(200 * time.Millisecond)

	go func() {
		defer close(jobs)

		for i := 0; i < 1000000; {
			select {
			case jobs <- i:
				i++
			case <-ticker.C:
				select {
				case <-ctx.Done():
					return
				case scale <- 1:
					fmt.Println("scale up request sent...")
				}
			}
		}
	}()

	// read results
	for r := range results {
		fmt.Println("result:", r)
	}

	fmt.Println("all workers stopped cleanly")
}
