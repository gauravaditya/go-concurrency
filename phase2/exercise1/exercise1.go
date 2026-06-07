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
	var wg sync.WaitGroup
	result := make(chan int)

	for range workers {
		wg.Add(1)
		go worker(jobs, result, &wg)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}

func worker(jobs <-chan int, result chan<- int, wg *sync.WaitGroup) {
	for v := range jobs {
		result <- (v * v)
	}
	wg.Done()
}
