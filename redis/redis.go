package redis

import (
	"time"

	"github.com/cyub/hyper/app"
	"github.com/go-redis/redis/v7"
)

var client *redis.Client

// Provider use for mount to app bootstrap
func Provider() app.ComponentMount {
	return func(app *app.App) error {
		return Init(app)
	}
}

// Instance return the instance of redis.Client
func Instance() *redis.Client {
	return client
}

// Init use for init redis
func Init(app *app.App) error {
	readTimeout := app.Config.GetInt("redis.read_timeout", 1)
	writeTimeout := app.Config.GetInt("redis.write_timeout", 3)
	poolTimeout := app.Config.GetInt("redis.pool_timeout", 5)
	opts := &redis.Options{
		Addr:         app.Config.GetString("redis.addr", "localhost:6379"),
		Password:     app.Config.GetString("redis.password", ""),
		DB:           app.Config.GetInt("redis.db", 0),
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		PoolSize:     app.Config.GetInt("redis.pool_size", 10),
		PoolTimeout:  time.Duration(poolTimeout) * time.Second,
	}

	client = redis.NewClient(opts)
	if _, err := client.Ping().Result(); err != nil {
		app.Logger.Errorf("redis ping failure %+v, %s", opts, err.Error())
		return err
	}
	app.Logger.Info("redis ping is ok")
	return nil
}
