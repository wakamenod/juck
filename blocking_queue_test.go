package juck

import (
	"sync"
	"testing"
)

func TestBlockingQueue(t *testing.T) {
	bq := NewBlockingQueue()
	var wg sync.WaitGroup

	// Consumer: Takes items from the queue
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			item := bq.Take()
			if item != i {
				t.Errorf("Expected %d, got %v", i, item)
			}
		}
	}()

	// Producer: Puts items into the queue
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			bq.Put(i)
		}
	}()

	wg.Wait()
}
