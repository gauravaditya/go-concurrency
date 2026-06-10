package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-concurrency/phase2/types"
)

func worker(ctx context.Context, jobs <-chan types.Job, out chan<- types.Result, wg *sync.WaitGroup) {
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
			case out <- types.Result{ID: v.ID, Value: v.Value}:
			}
		}
	}
}

func orderByID(ctx context.Context, results <-chan types.Result, indexOffset int) <-chan types.Result {
	orderedResults := make(chan types.Result)
	currentOffset := indexOffset
	futureResults := make(map[int]types.Result)

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
	jobs <-chan types.Job,
) <-chan types.Result {
	var wg sync.WaitGroup
	result := make(chan types.Result)

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
	jobs := make(chan types.Job)

	go func() {
		defer close(jobs)

		for i := 1; i <= 5000000; i++ {
			jobs <- types.Job{ID: i, Value: i}
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
