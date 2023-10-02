package juck

import (
	"sync"
)

// BlockingQueue is a simple implementation of a blocking queue.
type BlockingQueue[T any] struct {
	queue []T
	mu    sync.Mutex
	cond  *sync.Cond
}

// NewBlockingQueue creates a new BlockingQueue.
func NewBlockingQueue[T any]() *BlockingQueue[T] {
	bq := &BlockingQueue[T]{
		queue: make([]T, 0),
	}
	bq.cond = sync.NewCond(&bq.mu)
	return bq
}

// Put adds an item to the queue.
func (bq *BlockingQueue[T]) Put(item T) {
	bq.mu.Lock()
	bq.queue = append(bq.queue, item)
	bq.cond.Signal()
	bq.mu.Unlock()
}

// Take removes and returns an item from the queue.
func (bq *BlockingQueue[T]) Take() T {
	bq.mu.Lock()
	for len(bq.queue) == 0 {
		bq.cond.Wait()
	}
	item := bq.queue[0]
	bq.queue = bq.queue[1:]
	bq.mu.Unlock()
	return item
}

func (bq *BlockingQueue[T]) Len() int {
	bq.mu.Lock()
	res := len(bq.queue)
	bq.mu.Unlock()
	return res
}
