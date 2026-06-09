package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID    int
	Value int
}

type Result struct {
	ID    int
	Value int
}

func worker(ctx context.Context, jobs <-chan Job, out chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

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
			case out <- Result{ID: v.ID, Value: v.Value}:
			}
		}
	}
}

func orderByID(ctx context.Context, results <-chan Result, indexOffset int) <-chan Result {
	orderedResults := make(chan Result)
	currentOffset := indexOffset
	futureResults := make(map[int]Result)

	go func() {
		defer close(orderedResults)
		for {
			select {
			case <-ctx.Done():
				return
			case res, ok := <-results:
				if !ok {
					return
				}

				futureResults[res.ID] = res

				for {
					storedRes, exists := futureResults[currentOffset]
					if !exists {
						break
					}

					select {
					case <-ctx.Done():
						return
					case orderedResults <- storedRes:
						delete(futureResults, currentOffset)
						currentOffset++
					}
				}
			}
		}
	}()

	return orderedResults
}

func workerPool(
	ctx context.Context,
	workers int,
	jobs <-chan Job,
) <-chan Result {
	var wg sync.WaitGroup
	result := make(chan Result)

	for range workers {
		wg.Add(1)
		go worker(ctx, jobs, result, &wg)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return orderByID(ctx, result, 1)
}

func main() {
	jobs := make(chan Job)

	go func() {
		defer close(jobs)

		for i := 1; i <= 5000000; i++ {
			jobs <- Job{ID: i, Value: i}
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())

	results := workerPool(ctx, 3, jobs)

	time.AfterFunc(1*time.Second, cancel)

	for r := range results {
		fmt.Println(r)
	}

	fmt.Println("cancelled context...")
}
