package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second*5, cancel)

	jobs := make(chan int)

	go func() {
		defer close(jobs)
		for i := range 5 {
			jobs <- i
		}
	}()

	results := workerPool(ctx, 3, jobs)

	for n := range results {
		fmt.Println(n)
	}
}

func workerPool(ctx context.Context, workers int, jobs <-chan int) <-chan int {
	var wg sync.WaitGroup
	results := make(chan int)

	for range workers {
		wg.Add(1)
		go worker(ctx, jobs, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func worker(ctx context.Context, in <-chan int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-in:
			if !ok {
				return
			}

			// do some work, fail ocassionally
			err := process(v)
			if err != nil {
				err = retry(
					ctx,
					func() error {
						return process(v)
					},
					3,
				)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			select {
			case <-ctx.Done():
				return
			case out <- v:
			}
		}
	}
}

func process(job int) error {
	if job%3 == 0 {
		return fmt.Errorf("error for val: %d", job)
	}

	return nil
}

func retry(ctx context.Context, work func() error, attempts int) error {
	delay := time.Millisecond * 500

	var err error
	for attempt := 1; attempt <= attempts; attempt++ {
		timer := time.NewTimer(delay * time.Duration(attempt))
		select {
		case <-ctx.Done():
			return fmt.Errorf("ctx cancelled during retry...")
		case <-timer.C:
			fmt.Println("retry attempt:", attempt)
			err = work()
			if err == nil {
				return nil
			}

			timer.Stop()
		}
	}

	return err
}
