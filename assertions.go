package juck

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wakamenod/juck/timeunit"
)

// AssertEqualWithin は、指定した時間内に関数が期待される値を返すかどうかを検証します。
func AssertEqualWithin[T comparable](t *testing.T, duration time.Duration, f *Future[T], expected T) {
	t.Helper()

	select {
	case <-time.After(duration):
		t.Fatal("Function did not succeed within expected duration")
	case result := <-f.Ft():
		require.EqualValues(t, expected, result)
	}
}

func AssertWithinRange(t *testing.T, duration time.Duration, start, end int64, unit timeunit.TimeUnit) {
	t.Helper()

	// durationがstartとendの間にあることを確認
	require.True(t, duration >= unit.MultiplyWithInt(start) && duration <= unit.MultiplyWithInt(end),
		"Expected duration to be within %v and %v, but got %v", start, end, duration)
}
