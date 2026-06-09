package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	jobs := make(chan int)

	go func() {
		defer close(jobs)

		for i := 1; i <= 500000; i++ {
			jobs <- i
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())

	results := workerPool(ctx, 3, jobs)

	time.AfterFunc(time.Second, cancel)

	for r := range results {
		fmt.Println(r)
	}

	fmt.Println("cancelled context...")

}

func workerPool(
	ctx context.Context,
	workers int,
	jobs <-chan int,
) <-chan int {
	var wg sync.WaitGroup
	result := make(chan int)

	for range workers {
		wg.Add(1)
		go worker(ctx, jobs, result, &wg)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}

func worker(ctx context.Context, jobs <-chan int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		fmt.Println("worker: ctx cancelled..")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-jobs:
			if !ok {
				return
			}
			select {
			case <-ctx.Done():
				return
			case out <- v:
			}
		}
	}

}
