// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package selector

import "github.com/cyub/hyper/pkg/registry"

// Options struct
type Options struct {
	Registry registry.Registry
	Strategy Strategy
}

// Option struct
type Option func(*Options)

// WithRegistry set registry option
func WithRegistry(registry registry.Registry) Option {
	return func(o *Options) {
		o.Registry = registry
	}
}

// WithStrategy set strategy option
func WithStrategy(strategy Strategy) Option {
	return func(o *Options) {
		o.Strategy = strategy
	}
}
