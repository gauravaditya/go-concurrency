package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	wg := sync.WaitGroup{}
	jobs := make(chan int)

	results := pool(ctx, 3, 10, jobs)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		defer close(jobs)

		for i := 0; i < 1000000; {
			jobs <- i
			i++
		}
	}(&wg)

	fmt.Println("waiting for results...")
	// read results
	for r := range results {
		fmt.Println("result:", r)
	}

	wg.Wait()
	fmt.Println("all workers stopped cleanly")
}
