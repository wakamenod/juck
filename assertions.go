package juck

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// AssertWithinDuration は、指定した時間内に関数が期待される値を返すかどうかを検証します。
func AssertWithinDuration[T comparable](t *testing.T, duration time.Duration, f *Future[T], expected T) {
	t.Helper()

	select {
	case <-time.After(duration):
		t.Fatal("Function did not succeed within expected duration")
	case result := <-f.Ft():
		require.EqualValues(t, expected, result)
	}
}
