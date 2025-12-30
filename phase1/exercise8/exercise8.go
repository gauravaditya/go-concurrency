package main

import "sync"

type MutexMap struct {
	mu sync.Mutex
	m  map[int]int
}

func NewMutexMap() *MutexMap {
	return &MutexMap{m: make(map[int]int)}
}

func (mm *MutexMap) Get(k int) (int, bool) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	v, ok := mm.m[k]
	return v, ok
}

func (mm *MutexMap) Set(k int, v int) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	mm.m[k] = v
}

// --------------------

type RWMutexMap struct {
	mu sync.RWMutex
	m  map[int]int
}

func NewRWMutexMap() *RWMutexMap {
	return &RWMutexMap{m: make(map[int]int)}
}

func (rm *RWMutexMap) Get(k int) (int, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	v, ok := rm.m[k]
	return v, ok
}

func (rm *RWMutexMap) Set(k int, v int) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.m[k] = v
}
