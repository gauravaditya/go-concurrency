package main

import (
	"sync"
	"testing"
)

func benchMap(b *testing.B, readers, writers int, get func(int) (int, bool), set func(int, int)) {
	var wg sync.WaitGroup
	b.ResetTimer()

	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < b.N; j++ {
				get(j)
			}
		}(i)
	}

	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < b.N; j++ {
				set(j, j)
			}
		}(i)
	}

	wg.Wait()
}

func BenchmarkMutex_ReadHeavy(b *testing.B) {
	m := NewMutexMap()
	benchMap(b, 19, 1, m.Get, m.Set)
}

func BenchmarkRWMutex_ReadHeavy(b *testing.B) {
	m := NewRWMutexMap()
	benchMap(b, 19, 1, m.Get, m.Set)
}

// Repeat for 10/10 and 1/19
