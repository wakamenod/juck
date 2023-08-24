package juck

import (
	"sync"
	"testing"
	"time"

	"github.com/wakamenod/juck/timeunit"
)

func TestScheduleWithDelay(t *testing.T) {
	pool := NewScheduledThreadPool(3)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	var execTime time.Duration
	start := time.Now()
	pool.ScheduleWithDelay(func() {
		execTime = time.Since(start)
		wg.Done()
	}, 3000, timeunit.Millisecond)

	wg.Wait()

	AssertWithinRange(t, execTime, 3, 4, timeunit.Second)
}
