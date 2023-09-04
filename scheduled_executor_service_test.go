package juck

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wakamenod/juck/timeunit"
)

func TestScheduleWithDelay(t *testing.T) {
	t.Parallel()
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

func TestScheduleAtFixedRate(t *testing.T) {
	t.Parallel()
	pool := NewScheduledThreadPool(3)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	start := time.Now()
	execTimes := make([]time.Duration, 0)

	pool.ScheduleAtFixedRate(func() {
		execTimes = append(execTimes, time.Since(start))
		if len(execTimes) == 3 {
			wg.Done()
		}
	}, 1400, 500, timeunit.Millisecond)

	wg.Wait()
	pool.Shutdown()

	require.Equal(t, 3, len(execTimes))
	AssertWithinRange(t, execTimes[0], 1400, 1450, timeunit.Millisecond)
	AssertWithinRange(t, execTimes[1], 1900, 1950, timeunit.Millisecond)
	AssertWithinRange(t, execTimes[2], 2400, 2450, timeunit.Millisecond)
}

func TestScheduleAtFixedRate2(t *testing.T) {
	t.Parallel()
	// TODO Poolの数を3にすると動かない。バグ調査必須 TODO
	pool := NewScheduledThreadPool(4)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	start := time.Now()
	execTimes := make([]time.Duration, 0)

	pool.ScheduleAtFixedRate(func() {
		execTimes = append(execTimes, time.Since(start))

		fmt.Printf("@@@@@@@@@@ task started(%d): %v\n", len(execTimes), time.Since(start))

		// TODO 非同期用Slice
		if len(execTimes) == 1 || len(execTimes) == 3 {
			time.Sleep(time.Millisecond * 1200)
		}

		fmt.Printf("@@@@@@@@@@ task ended: %v\n", time.Since(start))

		if len(execTimes) == 4 {
			wg.Done()
		}
	}, 2400, 500, timeunit.Millisecond)

	wg.Wait()
	fmt.Println("shutdown")
	pool.Shutdown()

	require.Equal(t, 4, len(execTimes))
	AssertWithinRange(t, execTimes[0], 1400, 1450, timeunit.Millisecond)
	AssertWithinRange(t, execTimes[1], 1900, 1950, timeunit.Millisecond)
	AssertWithinRange(t, execTimes[2], 2400, 2450, timeunit.Millisecond)
}
