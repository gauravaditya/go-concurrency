package main

import (
	"fmt"
	"sync"
)

func main() {
	jobs := make(chan int)

	go func() {
		defer close(jobs)

		for i := 1; i <= 5; i++ {
			jobs <- i
		}
	}()

	results := workerPool(3, jobs)

	for r := range results {
		fmt.Println(r)
	}
}

func workerPool(
	workers int,
	jobs <-chan int,
) <-chan int {
	resultChans := make([]<-chan int, 0)

	for range workers {
		resultChans = append(resultChans, worker(jobs))
	}

	return merge(resultChans...)
}

func worker(jobs <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for v := range jobs {
			out <- (v * v)
		}
		close(out)
	}()

	return out
}

func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	merged := make(chan int)

	for _, ch := range channels {
		wg.Add(1)
		go forward(ch, merged, &wg)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func forward(in <-chan int, out chan int, wg *sync.WaitGroup) {
	for v := range in {
		out <- v
	}
	wg.Done()
}
