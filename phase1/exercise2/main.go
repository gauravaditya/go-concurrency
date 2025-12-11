package main

import (
	"context"
	"fmt"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs, results := startWorkerPool(ctx, 3)

	go func() {
		// send jobs
		for i := 0; i < 10; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// read results
	for r := range results {
		fmt.Println("result:", r)
	}

	fmt.Println("all workers stopped cleanly")
}
