package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	// Simulated work function
	work := func(ctx context.Context, v int) (int, error) {
		select {
		case <-time.After(300 * time.Millisecond):
			return v * v, nil
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}

	pool := NewPool(3, 5, 500*time.Millisecond, work)

	// ---- submit jobs concurrently ----
	var submitWG sync.WaitGroup
	for i := 0; i < 10; i++ {
		i := i
		submitWG.Add(1)
		go func() {
			defer submitWG.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			if err := pool.Submit(ctx, i); err != nil {
				fmt.Println("[submit error]", err)
			} else {
				fmt.Println("[submitted]", i)
			}
		}()
	}

	// ---- read results concurrently ----
	go func() {
		for r := range pool.Results() {
			fmt.Println("[result]", r)
		}
		fmt.Println("[results closed]")
	}()

	// ---- shutdown while submits may still be happening ----
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("[shutdown]")
		pool.Shutdown()
	}()

	submitWG.Wait()

	// Give workers time to drain / exit
	time.Sleep(2 * time.Second)
	fmt.Println("main exit")
}
