package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan int)
	var wg sync.WaitGroup

	wg.Add(1)
	go worker(ctx, 1, jobs, &wg)

	// Send some jobs
	for i := 0; i < 5; i++ {
		jobs <- i
	}

	// Cancel early (worker should exit immediately)
	cancel()

	// No one closes jobs here; worker must still exit.
	wg.Wait()
	fmt.Println("worker exited cleanly")
}
