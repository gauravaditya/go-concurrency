package main

import (
	"context"
	"fmt"
	"sync"
)

func worker(ctx context.Context, id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("closing worker %d\n", id)

			return
		case i, ok := <-jobs:
			if !ok {
				// channel closed: no more jobs
				fmt.Printf("worker %d exiting (channel closed)\n", id)

				return
			}
			fmt.Printf("worker %d: %d\n", id, i)
		}
	}
}
