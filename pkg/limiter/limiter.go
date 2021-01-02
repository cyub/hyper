package limiter

import (
	"fmt"
	"time"
)

// Limiter interface
type Limiter interface {
	Acquire(string, Limit) (Result, error)
}

// Limit struct define Limit
type Limit struct {
	Rate   int
	Burst  int
	Period time.Duration
}

// KeyPrefix the prefix of key
const KeyPrefix = "hyper:limiter:"

func (l Limit) String() string {
	return fmt.Sprintf("%d req/%s (burst %d)", l.Rate, l.Period, l.Burst)
}

// PerSecond return limit base per second
func PerSecond(rate int) Limit {
	return Limit{
		Rate:   rate,
		Period: time.Second,
		Burst:  rate,
	}
}

// PerMinute return limit base per minute
func PerMinute(rate int) Limit {
	return Limit{
		Rate:   rate,
		Period: time.Minute,
		Burst:  rate,
	}
}

// PerHour return limit base per hour
func PerHour(rate int) Limit {
	return Limit{
		Rate:   rate,
		Period: time.Hour,
		Burst:  rate,
	}
}

// Result define the result of acquire lock
type Result struct {
	Limit      Limit
	Acquired   bool
	Remaining  int
	RetryAfter int
	ResetAfer  int
}
