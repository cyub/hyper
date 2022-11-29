// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// options.go
package registry

import (
	"context"
	"time"
)

// Options struct
type Options struct {
	Address []string
	Timeout time.Duration
	Debug   bool
	Context context.Context
}

// Option func
type Option func(*Options)

// WithAddr set address option
func WithAddr(address ...string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

// WithTimeout set timeout option
func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

// WithContext set context option
func WithContext(context context.Context) Option {
	return func(o *Options) {
		o.Context = context
	}
}

// WithDebug set debug option
func WithDebug(debug bool) Option {
	return func(o *Options) {
		o.Debug = debug
	}
}
