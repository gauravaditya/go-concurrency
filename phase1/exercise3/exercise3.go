package main

import (
	"context"
	"fmt"
	"math"
)

func generator(ctx context.Context, nums []int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, n := range nums {
			fmt.Println("generating: ", n)
			select {
			case <-ctx.Done():
				return
			case out <- n:
			}
		}
	}()

	return out
}

func squarer(ctx context.Context, in <-chan int) <-chan int {
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
				case out <- n * n:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out
}

func printer(ctx context.Context, in <-chan int) {
	for {
		// time.Sleep(1 * time.Second) // adding delay also delay the pipeline due to backpressure
		select {
		case <-ctx.Done():
			return
		case n, ok := <-in:
			if !ok {
				return
			}

			fmt.Println(math.Sqrt(float64(n)), "->", n)
		}
	}
}
