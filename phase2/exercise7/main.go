package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func producer(ctx context.Context, id int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for i := 1; i <= 5; i++ {
			select {
			case <-ctx.Done():
				return
			case out <- i * id:
			}
		}
	}()

	return out
}

func merge(ctx context.Context, channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	result := make(chan int)

	for _, ch := range channels {
		wg.Add(1)
		go forward(ctx, ch, result, &wg)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}

func forward(ctx context.Context, in <-chan int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-in:
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

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second, cancel)

	out := merge(
		ctx,
		producer(ctx, 1),
		producer(ctx, 10),
		producer(ctx, 100),
	)

	for v := range out {
		fmt.Println(v)
	}
}
