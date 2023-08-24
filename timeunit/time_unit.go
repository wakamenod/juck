package timeunit

import "time"

// TimeUnit represents a duration with various units of time.
type TimeUnit time.Duration

const (
	Nanosecond  TimeUnit = TimeUnit(time.Nanosecond)
	Microsecond          = TimeUnit(time.Microsecond)
	Millisecond          = TimeUnit(time.Millisecond)
	Second               = TimeUnit(time.Second)
	Minute               = TimeUnit(time.Minute)
	Hour                 = TimeUnit(time.Hour)
)

// ToDuration converts TimeUnit to time.Duration
func (t TimeUnit) ToDuration() time.Duration {
	return time.Duration(t)
}

// MultiplyWithInt multiplies TimeUnit with an integer and returns the resulting time.Duration
func (t TimeUnit) MultiplyWithInt(factor int64) time.Duration {
	return t.ToDuration() * time.Duration(factor)
}
