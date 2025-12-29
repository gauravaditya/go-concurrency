package main

import (
	"fmt"
	"sync"
)

func main() {
	input := make([]int, 0, 1000000)

	for i := range cap(input) {
		input = append(input, i)
	}

	result := ParallelReduce(input, 3, 0, func(a, b int) int {
		return a + b
	})

	fmt.Println("result ->", result)
}

func ParallelReduce[T any](
	in []T,
	workers int,
	zero T,
	fn func(T, T) T,
) T {
	var out T
	var wg sync.WaitGroup

	jobs := make(chan T)
	results := make(chan T)

	for range workers {
		wg.Add(1)
		go worker(jobs, results, &wg, zero, fn)
	}

	go func() {
		for _, v := range in {
			jobs <- v
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	out = zero
	for v := range results {
		out = fn(out, v)
	}

	return out
}

func worker[T any](jobs <-chan T, results chan<- T, wg *sync.WaitGroup, zero T, fn func(T, T) T) {
	defer wg.Done()

	acc := zero
	for data := range jobs {
		acc = fn(acc, data)
	}

	results <- acc
}
