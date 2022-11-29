// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package registry

import (
	"context"
	"time"

	"github.com/cyub/hyper/app"
	"github.com/cyub/hyper/pkg/registry"
	"github.com/cyub/hyper/pkg/registry/consul"
)

var defaultRegistry registry.Registry

// Provider use for mount to app bootstrap
func Provider() app.ComponentMount {
	return func(app *app.App) (err error) {
		addr := app.Config.GetString("registry.addr", "localhost:8500")
		timeout := app.Config.GetDuration("registry.timeout", 1000)
		debug := app.Config.GetBool("registry.debug", false)

		defaultRegistry = consul.NewRegistry(
			registry.WithAddr(addr),
			registry.WithTimeout(timeout*time.Millisecond),
			registry.WithDebug(debug),
		)
		if _, err = defaultRegistry.GetService(context.Background(), "test"); err != nil {
			app.Logger.Errorf("registry ping failure %s", err.Error())
			return
		}
		app.Logger.Info("registry ping is ok")
		return nil
	}
}

// Instance return the instance of registry
func Instance() registry.Registry {
	return defaultRegistry
}
