// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package queue

import (
	"errors"
	"time"

	"github.com/cyub/hyper/app"
	"github.com/cyub/hyper/logger"
	"github.com/cyub/hyper/pkg/config"
	"github.com/cyub/hyper/pkg/queue"
	kafkaQueue "github.com/cyub/hyper/pkg/queue/kafka"
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

		var err error
		switch baseOpts.Driver {
		case "redis":
			queuer, err = createRedisQueue(app.Config, baseOpts)
		case "kafka":
			queuer, err = createKafkaQueue(app.Config, baseOpts)
		default:
			return queue.ErrInvalidProvider
		}
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

func createKafkaQueue(config *config.Config, baseOpts queue.Options) (queue.Queuer, error) {
	opts := kafkaQueue.Options{
		Options:      baseOpts,
		Brokers:      config.GetStringSlice("queue.kafka_brokers", []string{"localhost:9092"}),
		GroupID:      config.GetString("queue.kafka_group_id", "hyper-consume-group"),
		KafkaVersion: config.GetString("queue.kafka_version", "2.1.1"),
	}

	return kafkaQueue.New(opts), nil
}

// Instance return instance of Queue
func Instance() queue.Queuer {
	return queuer
}
