package backoff

// Exponential Backoff with allow Jitter
// @see https://amazonaws-china.com/cn/blogs/architecture/exponential-backoff-and-jitter/

import (
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

// Backoff struct
type Backoff struct {
	Factor   float64
	Jitter   bool
	Max      time.Duration
	Min      time.Duration
	attempts uint64
}

const (
	// DefaultFactor factor
	DefaultFactor = math.E
	// DefaultJitter jitter
	DefaultJitter = true
	// DefaultMaxInterval max interval
	DefaultMaxInterval = 3 * time.Minute
	// DefaultMinInterval min interval
	DefaultMinInterval = 100 * time.Millisecond
	// MaxInt64 the max int64 to avoid overflow
	MaxInt64 = float64(math.MaxInt64 - 1024)
)

// New create backoff instance
func New(factor float64, jitter bool, max, min time.Duration) *Backoff {
	if factor <= 1 {
		factor = math.E
	}
	if min <= 0 {
		min = DefaultMinInterval
	}
	if max <= 0 {
		max = DefaultMaxInterval
	}

	if min >= max {
		min = DefaultMinInterval
	}

	if max < min {
		max = DefaultMaxInterval
	}

	return &Backoff{
		Factor:   factor,
		Jitter:   jitter,
		Max:      max,
		Min:      min,
		attempts: 0,
	}
}

// Default return backoff instance with default value
func Default() *Backoff {
	return New(
		DefaultFactor,
		DefaultJitter,
		DefaultMaxInterval,
		DefaultMinInterval,
	)
}

// WithoutJitter return backoff instance with default value except jitter
func WithoutJitter() *Backoff {
	return New(
		DefaultFactor,
		false,
		DefaultMaxInterval,
		DefaultMinInterval,
	)
}

// Next return next interval time to backoff
func (b *Backoff) Next() time.Duration {
	attempts := atomic.AddUint64(&b.attempts, 1) - 1
	minf := float64(b.Min)
	durf := math.Pow(b.Factor, float64(attempts)) * minf
	if b.Jitter {
		durf = minf + rand.Float64()*(durf-minf)
	}

	if durf > MaxInt64 {
		return b.Max
	}

	next := time.Duration(durf)
	if next > b.Max {
		return b.Max
	}

	if next < b.Min {
		return b.Min
	}

	return next
}

// Reset use for reset backoff attempt times
func (b *Backoff) Reset() {
	atomic.StoreUint64(&b.attempts, 0)
}

// Attempts return the attempt times
func (b *Backoff) Attempts() uint64 {
	return b.attempts
}
