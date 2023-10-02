package juck

import (
	"testing"
	"time"
)

func TestSubmitInt(t *testing.T) {
	pool := NewFixedThreadPool(1)

	future := pool.SubmitInt(func() int {
		time.Sleep(100 * time.Millisecond)
		return 42
	})
	future2 := pool.SubmitInt(func() int {
		time.Sleep(100 * time.Millisecond)
		return 53
	})

	AssertEqualWithin[int](t, time.Millisecond*1000, future, 42)
	AssertEqualWithin[int](t, time.Millisecond*1000, future2, 53)
}
