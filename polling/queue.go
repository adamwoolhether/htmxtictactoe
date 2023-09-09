package pubsub

import "sync"

type Queue[T any] struct {
	mu    sync.RWMutex
	items []T
	cap   int
}

func NewQueue[T any](cap int) *Queue[T] {
	q := Queue[T]{
		items: make([]T, 0),
		cap:   cap,
	}

	return &q
}

func (q *Queue[T]) Enqueue(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Remove the oldest item if queue is at capacity.
	if len(q.items) == q.cap {
		q.items = q.items[1:]
	}

	q.items = append(q.items, item)
}

func (q *Queue[T]) Copy() []T {
	q.mu.RLock()
	defer q.mu.RUnlock()

	cpy := make([]T, len(q.items))
	copy(cpy, q.items)

	return cpy
}
