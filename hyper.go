package hyper

import (
	"encoding/json"
	"fmt"

	"github.com/cyub/hyper/app"
	_config "github.com/cyub/hyper/config"
	"github.com/cyub/hyper/logger"
	"github.com/cyub/hyper/mysql"
	"github.com/cyub/hyper/pkg/config"
	"github.com/cyub/hyper/queue"
	"github.com/cyub/hyper/redis"
	_redis "github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// NewApp return application
func NewApp(opts ...app.Option) *app.App {
	return app.NewApp(opts...)
}

// DB return gorm.DB
func DB() *gorm.DB {
	return mysql.Instance()
}

// Redis return redis.Client
func Redis() *_redis.Client {
	return redis.Instance()
}

// Config return config.Config
func Config() *config.Config {
	return _config.Instance()
}

// Logger return logrus.Logger
func Logger() *logrus.Logger {
	return logger.Instance()
}

// Queue return queue.Queue
func Queue() *queue.Queue {
	return queue.Instance()
}

// NewJob return job
func NewJob(name string, data interface{}) (*queue.Job, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("job payload can't marshal %s", err.Error())
	}
	job := &queue.Job{
		Name:    name,
		Payload: payload,
	}
	return job, nil
}

// InQueue use for job enqueue
func InQueue(name string, data interface{}) error {
	job, err := NewJob(name, data)
	if err != nil {
		return err
	}
	err = queue.Instance().In(job)
	if err != nil {
		return err
	}

	return nil
}

// WithName use for set app's name
func WithName(name string) app.Option {
	return func(o *app.Options) {
		o.Name = name
	}
}

// WithAddr use for set app's addr
func WithAddr(addr string) app.Option {
	return func(o *app.Options) {
		o.Addr = addr
	}
}

// WithRunMode use for set app's run mode
func WithRunMode(mode string) app.Option {
	return func(o *app.Options) {
		o.RunMode = mode
	}
}

// WithCfgAddr use for set app's config center addr
func WithCfgAddr(cfgAddr string) app.Option {
	return func(o *app.Options) {
		o.CfgCenterAddr = cfgAddr
	}
}

// WithCfgPath use for set app's config center path
func WithCfgPath(cfgPath string) app.Option {
	return func(o *app.Options) {
		o.CfgCenterPath = cfgPath
	}
}
