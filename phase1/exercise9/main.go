package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go StartPeriodicWorker(ctx, 500*time.Millisecond, func() {
		fmt.Println("tick", time.Now())
	})

	time.Sleep(2 * time.Second)
	cancel()

	time.Sleep(1 * time.Second)
	fmt.Println("done")
}
