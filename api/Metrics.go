package api

import (
	"sync"
)

type AtomicInt struct {
	mu sync.Mutex // A lock than can be held by just one goroutine at a time.
	n  int
}

// Add adds n to the AtomicInt as a single atomic operation.
func (a *AtomicInt) Add(n int) {
	a.mu.Lock() // Wait for the lock to be free and then take it.
	a.n += n
	a.mu.Unlock() // Release the lock.
}

// Value returns the value of a.
func (a *AtomicInt) Value() int {
	a.mu.Lock()
	n := a.n
	a.mu.Unlock()
	return n
}
