package redis

import (
	"time"

	"github.com/cyub/hyper/app"
	"github.com/go-redis/redis/v7"
)

var (
	client        *redis.Client
	clusterClient *redis.ClusterClient
)

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

// ClusterProvider use for mount to app bootstrap
func ClusterProvider() app.ComponentMount {
	return func(app *app.App) error {
		return InitCluster(app)
	}
}

// ClusterInstance return the instance of redis.Client
func ClusterInstance() *redis.ClusterClient {
	return clusterClient
}

// Init use for init redis client
func Init(app *app.App) error {
	opts := &redis.Options{
		Addr:         app.Config.GetString("redis.addr", "localhost:6379"),
		Password:     app.Config.GetString("redis.password", ""),
		DB:           app.Config.GetInt("redis.db", 0),
		ReadTimeout:  app.Config.GetDuration("redis.read_timeout", 3) * time.Second,
		WriteTimeout: app.Config.GetDuration("redis.write_timeout", 3) * time.Second,
		PoolSize:     app.Config.GetInt("redis.pool_size", 50),
		PoolTimeout:  app.Config.GetDuration("redis.pool_timeout", 5) * time.Second,
	}

	client = redis.NewClient(opts)
	if _, err := client.Ping().Result(); err != nil {
		app.Logger.Errorf("redis ping failure %+v, %s", opts, err.Error())
		return err
	}
	app.Logger.Info("redis ping is ok")
	return nil
}

// InitCluster init cluster client
func InitCluster(app *app.App) error {
	opts := &redis.ClusterOptions{
		Addrs:           app.Config.GetStringSlice("redis_cluster.addrs", []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"}),
		MaxRedirects:    app.Config.GetInt("redis_cluster.max_redirects", 8),
		PoolSize:        app.Config.GetInt("redis_cluster.pool_size", 0),
		Password:        app.Config.GetString("redis_cluster.password", ""),
		ReadOnly:        app.Config.GetBool("redis_cluster.read_only", false),
		ReadTimeout:     app.Config.GetDuration("redis_cluster.read_timeout", 3) * time.Second,
		WriteTimeout:    app.Config.GetDuration("redis_cluster.write_timeout", 3) * time.Second,
		MaxRetryBackoff: app.Config.GetDuration("redis_cluster.max_retry_backoff", 512) * time.Millisecond,
		MinRetryBackoff: app.Config.GetDuration("redis_cluster.min_retry_backoff", 8) * time.Millisecond,
	}
	clusterClient = redis.NewClusterClient(opts)
	err := clusterClient.ForEachNode(func(rdb *redis.Client) error {
		return rdb.Ping().Err()
	})
	if err != nil {
		app.Logger.Errorf("redis-cluster ping failure %+v, %s", opts, err.Error())
		return err
	}
	app.Logger.Info("redis-cluster ping is ok")
	return nil
}
