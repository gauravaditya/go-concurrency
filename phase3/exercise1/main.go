package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func generate(ctx context.Context, n int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for i := range n {
			select {
			case <-ctx.Done():
				return
			case out <- i + 1:
			}
		}
	}()

	return out
}

func square(ctx context.Context, in <-chan int, workers int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for range workers {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case n, ok := <-in:
					if !ok {
						return
					}

					select {
					case <-ctx.Done():
						return
					case out <- n * n:
					}
				}
			}
		}(&wg)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func double(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case n, ok := <-in:
				if !ok {
					return
				}

				select {
				case <-ctx.Done():
					return
				case out <- n * 2:
				}
			}
		}
	}()

	return out
}

func filterEven(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case n, ok := <-in:
				if !ok {
					return
				}

				if n%2 == 0 {
					select {
					case <-ctx.Done():
						return
					case out <- n:
					}
				}
			}
		}
	}()

	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second, cancel)

	nums := generate(ctx, 10)
	squares := square(ctx, nums, 3)
	filtered := filterEven(ctx, squares)
	doubles := double(ctx, filtered)

	for n := range doubles {
		fmt.Println("number:", n)
	}
}
