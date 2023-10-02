package juck

import (
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

func TestScheduleWithDelay2(t *testing.T) {
	t.Parallel()
	pool := NewScheduledThreadPool(2)

	wg := &sync.WaitGroup{}
	execTimes := NewBlockingQueue[time.Duration]()

	start := time.Now()
	wg.Add(1)
	pool.ScheduleWithDelay(func() {
		execTimes.Put(time.Since(start))
		wg.Done()
	}, 1000, timeunit.Millisecond)

	wg.Add(1)
	pool.ScheduleWithDelay(func() {
		execTimes.Put(time.Since(start))
		wg.Done()
	}, 2000, timeunit.Millisecond)

	wg.Add(1)
	pool.ScheduleWithDelay(func() {
		execTimes.Put(time.Since(start))
		wg.Done()
	}, 3000, timeunit.Millisecond)

	wg.Wait()

	require.Equal(t, 3, execTimes.Len())
	AssertWithinRange(t, execTimes.Take(), 1000, 1050, timeunit.Millisecond)
	AssertWithinRange(t, execTimes.Take(), 2000, 2050, timeunit.Millisecond)
	AssertWithinRange(t, execTimes.Take(), 3000, 3050, timeunit.Millisecond)
}

func TestScheduleAtFixedRate(t *testing.T) {
	t.Parallel()
	pool := NewScheduledThreadPool(3)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	start := time.Now()
	execTimes := NewBlockingQueue[time.Duration]()

	pool.ScheduleAtFixedRate(func() {
		execTimes.Put(time.Since(start))
		if execTimes.Len() == 3 {
			wg.Done()
		}
	}, 1400, 500, timeunit.Millisecond)

	wg.Wait()
	pool.Shutdown()

	require.Equal(t, 3, execTimes.Len())
	AssertWithinRange(t, execTimes.Take(), 1400, 1450, timeunit.Millisecond)
	AssertWithinRange(t, execTimes.Take(), 1900, 1950, timeunit.Millisecond)
	AssertWithinRange(t, execTimes.Take(), 2400, 2450, timeunit.Millisecond)
}

func TestScheduleAtFixedRate2(t *testing.T) {
	t.Parallel()
	pool := NewScheduledThreadPool(3)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	start := time.Now()
	execTimes := NewBlockingQueue[time.Duration]()

	pool.ScheduleAtFixedRate(func() {
		execTimes.Put(time.Since(start))

		len := execTimes.Len()
		if len == 1 || len == 3 {
			time.Sleep(time.Millisecond * 1200)
		}
		if len == 4 {
			wg.Done()
		}
	}, 1400, 500, timeunit.Millisecond)

	wg.Wait()
	pool.Shutdown()

	require.Equal(t, 4, execTimes.Len())
	AssertWithinRange(t, execTimes.Take(), 1400, 1450, timeunit.Millisecond)
	AssertWithinRange(t, execTimes.Take(), 2600, 2650, timeunit.Millisecond)
	AssertWithinRange(t, execTimes.Take(), 2600, 2650, timeunit.Millisecond)
}

func TestScheduleAtFixedRate3(t *testing.T) {
	t.Parallel()
	pool := NewSingleScheduledExecutor()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	start := time.Now()
	execTimes := NewBlockingQueue[time.Duration]()

	pool.ScheduleAtFixedRate(func() {
		execTimes.Put(time.Since(start))

		len := execTimes.Len()
		if len == 1 {
			time.Sleep(time.Millisecond * 1200)
		}
		if len == 4 {
			wg.Done()
		}
	}, 2400, 500, timeunit.Millisecond)

	wg.Wait()
	pool.Shutdown()

	// If the task is not delayed it should fire at the following intervals:
	// 2400(ms) -> 2900 -> 3400 -> 3900

	require.Equal(t, 4, execTimes.Len())
	AssertWithinRange(t, execTimes.Take(), 2400, 2450, timeunit.Millisecond)
	AssertWithinRange(t, execTimes.Take(), 3600, 3650, timeunit.Millisecond)
	AssertWithinRange(t, execTimes.Take(), 3600, 3650, timeunit.Millisecond)
	AssertWithinRange(t, execTimes.Take(), 3900, 3950, timeunit.Millisecond)
}
