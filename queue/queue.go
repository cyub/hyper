package queue

import (
	"errors"
	"time"

	"github.com/cyub/hyper/app"
	"github.com/cyub/hyper/logger"
	"github.com/cyub/hyper/pkg/config"
	"github.com/cyub/hyper/pkg/queue"
	redisQueue "github.com/cyub/hyper/pkg/queue/redis"
	"github.com/cyub/hyper/redis"
)

var (
	queuer queue.Queuer
)

// Provider provide queue
func Provider(consumers ...map[string]queue.Consumer) app.ComponentMount {
	return func(app *app.App) error {
		baseOpts := queue.Options{
			Name:        app.Config.GetString("queue.name", ""),
			Driver:      app.Config.GetString("queue.driver", "redis"),
			WaitTimeOut: app.Config.GetDuration("queue.timeout", 0) * time.Second,
			Debug:       app.Config.GetBool("queue.debug", false),
			Metric:      app.Config.GetBool("queue.metric", true),
			Consume:     app.Config.GetBool("queue.consume", true),
			Parallel:    app.Config.GetInt("queue.parallel_number", 1),
			Logger:      logger.Instance(),
		}
		if baseOpts.Driver != "redis" {
			return queue.ErrInvalidProvider
		}

		var err error
		queuer, err = createRedisQueue(app.Config, baseOpts)
		if err != nil {
			return err
		}

		for _, consumer := range consumers {
			for name, c := range consumer {
				queue.RegisterConsumer(name, c)
			}
		}

		if baseOpts.Metric {
			queue.MetricTurnon()
		}

		if baseOpts.Consume {
			queuer.Run()
		}
		app.Logger.Infof("queue-%s ping is ok", baseOpts.Driver)
		return nil
	}
}

func createRedisQueue(config *config.Config, baseOpts queue.Options) (queue.Queuer, error) {
	redis := redis.Instance()
	if redis == nil {
		return nil, errors.New("please boot with redis as queue store backend")
	}

	opts := redisQueue.Options{
		Options: baseOpts,
		Backend: redis,
	}

	return redisQueue.New(opts), nil
}

// Instance return instance of Queue
func Instance() queue.Queuer {
	return queuer
}
