// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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
