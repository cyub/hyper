// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package queue

import (
	"fmt"

	"github.com/cyub/hyper/pkg/queue"
	"github.com/go-redis/redis/v7"
)

// Queue struct
type Queue struct {
	*queue.Base
	Options
}

// Options struct
type Options struct {
	queue.Options
	Backend *redis.Client
}

var _ queue.Queuer = (*Queue)(nil)

// New return queue base redis
func New(opts Options) *Queue {
	opts.Init()
	if opts.Backend == nil {
		panic("please provider redis as store backend")
	}
	q := Queue{
		Options: opts,
	}
	q.Base = queue.NewBase(opts.Options)
	q.Base.Enqueue = q.doEnqueue
	return &q
}

func (q *Queue) doEnqueue(msg []byte) error {
	return q.Backend.LPush(q.Name, msg).Err()
}

// Run use for queue run
func (q *Queue) Run() {
	go func() {
		for {
			data, err := q.Backend.BRPop(q.WaitTimeOut, q.Name).Result()
			if err == redis.Nil {
				continue
			}
			// redis error
			if err != nil {
				q.Logger.Errorf("queue backend redis error %s", err.Error())
				q.Backoff()
				continue
			}
			q.BackoffReset()
			if q.Debug {
				fmt.Printf("job data %s\n", data[1])
			}

			job, err := q.Parse(data[1])
			if err != nil {
				q.Logger.Errorf("job[%s] parse failure %s", data[1], err.Error())
				continue
			}
			fn, err := q.GetConsumer(job.GetName())
			if err != nil {
				q.Logger.Error(err)
			}
			q.SemAcquire(1)
			go q.Do(job, fn)
		}
	}()
}
