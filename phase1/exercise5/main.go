package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	jobs := make(chan int)

	results := autoScalerPool(ctx, 3, 5, jobs)

	go func() {
		defer close(jobs)

		for i := 0; i < 10; {
			fmt.Println("sending data...")
			jobs <- i
			i++
		}
	}()

	fmt.Println("waiting for results...")
	// read results
	for r := range results {
		fmt.Println("result:", r)
	}

	fmt.Println("all workers stopped cleanly")
}
