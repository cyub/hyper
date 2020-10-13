package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cyub/hyper/app"
	"github.com/cyub/hyper/logger"
	_redis "github.com/cyub/hyper/redis"
	"github.com/go-redis/redis/v7"
)

const (
	defaultQueue   = "hyper:queue:default"
	defaultTimeout = 5
)

// Queue struct
type Queue struct {
	Name        string
	Backend     *redis.Client
	WaitTimeOut time.Duration
	Consumers   map[string]Consumer
}

var queue *Queue

// Consumer define job consumer
type Consumer func(job *Job) error

// Provider use for mount to app bootstrap
func Provider(consumers ...map[string]Consumer) app.ComponentMount {
	return func(app *app.App) error {
		redis := _redis.Instance()
		if redis == nil {
			return errors.New("please boot with redis as queue store backend")
		}
		queue = &Queue{
			Name:        app.Config.GetString("queue.name", defaultQueue),
			Backend:     redis,
			WaitTimeOut: time.Duration(app.Config.GetInt("queue.timeout", defaultTimeout)) * time.Second,
			Consumers:   make(map[string]Consumer),
		}

		for _, consumer := range consumers {
			for name, c := range consumer {
				RegisterConsumer(name, c)
			}
		}

		for name, w := range _consumers {
			queue.AddConsumer(name, w)
		}
		go queue.Run()
		return nil
	}
}

// Instance return instance of Queue
func Instance() *Queue {
	return queue
}

// AddConsumer use for queue consumer
func (q *Queue) AddConsumer(name string, f Consumer) {
	q.Consumers[name] = f
}

// In use job enqueue
func (q *Queue) In(job Jober) error {
	data, err := job.Serialize()
	if err != nil {
		fmt.Printf("in job error %s", err.Error())
		return err
	}
	err = q.Backend.LPush(q.Name, data).Err()
	if err != nil {
		return err
	}
	return nil
}

// Run use for queue run
func (q *Queue) Run() {
	for {
		data, err := q.Backend.BRPop(q.WaitTimeOut, q.Name).Result()
		if err != nil {
			continue
		}

		job, err := q.Parse(data[1])
		if err != nil {
			panic(err)
		}

		name := job.GetName()
		if w, exist := q.Consumers[name]; exist {
			tries := job.GetMaxTries()
			for {
				if w(&job) == nil {
					break
				}
				tries--
				if tries <= 0 {
					break
				}
			}
		} else {
			logger.Errorf("job[%s] consumer don't exist", name)
		}
	}
}

// Parse use for parse job info
func (q *Queue) Parse(payload string) (j Job, err error) {
	err = json.Unmarshal([]byte(payload), &j)
	return
}
