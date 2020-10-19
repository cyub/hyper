package cache

import (
	"time"
)

// Cache interface
type Cache interface {
	Get(string, interface{}) error
	Set(string, interface{}, time.Duration) error
	Forever(string, interface{}) (err error)
	Remember(string, interface{}, SetCallback, time.Duration) error
	RememberForever(string, interface{}, SetCallback) (err error)
	Delete(string) error
	Key(string) string
}

// SetCallback type
type SetCallback func() error

// CacheError func type
type cacheError string

func (e cacheError) Error() string {
	return string(e)
}

// Nil error
const Nil = cacheError("cache miss")
