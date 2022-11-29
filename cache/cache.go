// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cache

import (
	"errors"

	"github.com/cyub/hyper/app"
	"github.com/cyub/hyper/pkg/cache"
	"github.com/cyub/hyper/redis"
)

var (
	c cache.Cache
)

// Nil mean cache don't exists
const Nil = cache.Nil

// Provider provide cache
func Provider() app.ComponentMount {
	return func(app *app.App) error {
		driver := app.Config.GetString("cache.driver", "redis")
		prefix := app.Config.GetString("cache.prefix", "hyper")
		if driver != "redis" {
			return errors.New("only support redis as cache store backend")
		}

		if redis.Instance() == nil {
			return errors.New("please boot with redis as cache store backend")
		}
		c = cache.NewRedis(redis.Instance(), prefix)
		app.Logger.Infof("cache-%s ping is ok", driver)
		return nil
	}
}

// Instance return instance of Cache
func Instance() cache.Cache {
	return c
}
