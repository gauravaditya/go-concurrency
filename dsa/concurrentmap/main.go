package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main() {
	input := make([]int, 0, 10)

	for i := range cap(input) {
		input = append(input, i)
	}

	fmt.Println(input)
	results := ParallelMap(input, 3, func(t int) string {
		fmt.Println("converting ->", t)
		return strconv.Itoa(t)
	})

	fmt.Println(results)
	for i, _ := range input {
		fmt.Println("conversion completed ->", input[i], results[i])
	}
}

func ParallelMap[T any, R any](
	in []T,
	workers int,
	fn func(T) R,
) []R {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	wg := sync.WaitGroup{}
	jobs := make(chan input[T, R], len(in))
	results := make(chan input[T, R])
	out := make([]R, len(in), len(in))

	for i := range workers {
		wg.Add(1)
		go worker(ctx, i, jobs, results, &wg, fn)
	}

	for idx, v := range in {
		jobs <- input[T, R]{v: v, idx: idx}
	}
	close(jobs)

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(results)
	}(&wg)

	for r := range results {
		out[r.idx] = r.r
	}

	return out
}

type input[T, R any] struct {
	v   T
	r   R
	idx int
}

func worker[T, R any](ctx context.Context, id int, jobs <-chan input[T, R], results chan<- input[T, R], wg *sync.WaitGroup, fn func(T) R) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-jobs:
			if !ok {
				return
			}

			data.r = fn(data.v)
			select {
			case <-ctx.Done():
				return
			case results <- data:
			}
		}
	}
}
