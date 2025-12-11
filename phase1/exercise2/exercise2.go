package main

import (
	"context"
	"fmt"
	"sync"
)

func worker(ctx context.Context, id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d: closing (context cancelled)\n", id)

			return
		case data, ok := <-jobs:
			if !ok {
				fmt.Printf("worker %d: exiting (channel closed)\n", id)

				return
			}

			fmt.Printf("worker %d: received -> %d\n", id, data)
			select {
			case results <- data:
			case <-ctx.Done():
				fmt.Printf("worker %d: closed while sending (context cancelled)\n", id)

				return
			}
		}
	}
}

func startWorkerPool(ctx context.Context, numWorkers int) (chan<- int, <-chan int) {
	var wg sync.WaitGroup

	jobs := make(chan int)
	results := make(chan int)

	for id := range numWorkers {
		wg.Add(1)
		go worker(ctx, id, jobs, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return jobs, results
}
