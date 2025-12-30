package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	q := NewBlockingQueue[int](5)

	var wg sync.WaitGroup

	// Consumers
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			v := q.Get()
			fmt.Printf("[consumer] got %d\n", v)
			time.Sleep(300 * time.Millisecond) // slow consumer
		}
	}()

	// Producers
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			fmt.Printf("[producer] putting %d\n", i)
			q.Put(i)
			fmt.Printf("[producer] put %d\n", i)
		}
	}()

	wg.Wait()
	fmt.Println("done")
}
