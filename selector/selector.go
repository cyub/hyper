// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package selector

import (
	"errors"

	"github.com/cyub/hyper/app"
	"github.com/cyub/hyper/pkg/selector"
	"github.com/cyub/hyper/registry"
)

var defaultSelector selector.Selector

// Provider use for mount to app bootstrap
func Provider() app.ComponentMount {
	return func(app *app.App) (err error) {
		registry := registry.Instance()
		if registry == nil {
			return errors.New("please boot with registry as selector")
		}
		defaultSelector = selector.NewSelector(selector.WithRegistry(registry))
		return nil
	}
}

// Instance return the instance of selector
func Instance() selector.Selector {
	return defaultSelector
}
