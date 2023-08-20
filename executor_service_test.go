package juck

import (
	"testing"
	"time"
)

func TestSubmitInt(t *testing.T) {
	pool := NewFixedThreadPool(3)

	future := pool.SubmitInt(func() int {
		time.Sleep(100 * time.Millisecond)
		return 42
	})
	pool.Shutdown()

	AssertWithinDuration[int](t, time.Millisecond*1000, future, 42)
}
