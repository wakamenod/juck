package juck

import (
	"sync"
)

// BlockingQueue is a simple implementation of a blocking queue.
type BlockingQueue struct {
	queue []any
	mu    sync.Mutex
	cond  *sync.Cond
}

// NewBlockingQueue creates a new BlockingQueue.
func NewBlockingQueue() *BlockingQueue {
	bq := &BlockingQueue{
		queue: make([]any, 0),
	}
	bq.cond = sync.NewCond(&bq.mu)
	return bq
}

// Put adds an item to the queue.
func (bq *BlockingQueue) Put(item any) {
	bq.mu.Lock()
	bq.queue = append(bq.queue, item)
	bq.cond.Signal()
	bq.mu.Unlock()
}

// Take removes and returns an item from the queue.
func (bq *BlockingQueue) Take() any {
	bq.mu.Lock()
	for len(bq.queue) == 0 {
		bq.cond.Wait()
	}
	item := bq.queue[0]
	bq.queue = bq.queue[1:]
	bq.mu.Unlock()
	return item
}
