package queue

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/cyub/hyper/pkg/backoff"
	"github.com/cyub/hyper/pkg/semaphore"
	"github.com/sirupsen/logrus"
)

var (
	defaultQueue                 = "hyper:queue:default"
	defaultTimeout time.Duration = 5
	// ErrInvalidProvider when provider invalid return
	ErrInvalidProvider = errors.New("config: invalid config center provider")
)

// Queuer interface
type Queuer interface {
	In(job Jober) (err error)
	Parse(payload string) (j Jober, err error)
	Do(job Jober, fn Consumer) (err error)
	RegisterConsumer(name string, f Consumer)
	Run()
}

var _ Queuer = (*Base)(nil)

// Base struct
type Base struct {
	Opts      Options
	sem       semaphore.Semaphore
	backoff   *backoff.Backoff
	consumers map[string]Consumer
	Debug     bool
	Enqueue   func([]byte) error
}

// NewBase return base queue
func NewBase(opts Options) *Base {
	b := &Base{
		consumers: make(map[string]Consumer),
		sem:       semaphore.NewSemaphore(opts.Parallel),
		backoff:   backoff.Default(),
		Debug:     opts.Debug,
		Opts:      opts,
	}

	for name, c := range Consumers() {
		b.RegisterConsumer(name, c)
	}
	return b
}

// SemAcquire use for acquire semaphore
func (b *Base) SemAcquire(n int) {
	b.sem.P(n)
}

// SemRelease use for release semphore
func (b *Base) SemRelease(n int) {
	b.sem.V(n)
}

// Backoff use for backoff
func (b *Base) Backoff() {
	time.Sleep(b.backoff.Next())
}

// BackoffReset use for reset backoff
func (b *Base) BackoffReset() {
	b.backoff.Reset()
}

// GetConsumer option
func (b *Base) GetConsumer(name string) (Consumer, error) {
	fn, exist := b.consumers[name]
	if !exist {
		return nil, fmt.Errorf("consumer[%s] don't exist", name)
	}

	return fn, nil
}

// In use for job enqueue
func (b *Base) In(job Jober) (err error) {
	defer func() {
		if err == nil {
			InSuccessInc(job.GetName())
		} else {
			InFailInc(job.GetName())
		}
	}()

	data, err := job.Serialize()
	if err != nil {
		b.Opts.Logger.Errorf("job enqueue error %s", err.Error())
		return err
	}

	if b.Enqueue == nil {
		panic("job enqueue function don't exist")
	}
	return b.Enqueue(data)
}

// RegisterConsumer use for register queue consumer
func (b *Base) RegisterConsumer(name string, f Consumer) {
	b.consumers[name] = f
}

// Parse use for parse job info
func (b *Base) Parse(payload string) (j Jober, err error) {
	var job Job
	err = json.Unmarshal([]byte(payload), &job)
	return &job, err
}

// Run use for loop run job handle
func (b *Base) Run() {
	panic("job consume function don't exist")
}

// Do use for run queue consume action
func (b *Base) Do(job Jober, fn Consumer) (err error) {
	startTime := time.Now()
	defer func() {
		b.SemRelease(1)
		recoverErr := recover()
		if err == nil && recoverErr == nil {
			ConsumSuccessInc(job.GetName())
		} else {
			ConsumFailInc(job.GetName())
		}
		ProcessTimeHist(job.GetName(), startTime)
		if recoverErr != nil {
			fmt.Printf("[Queue Panic] %s\n[stacktrace]\n%s", err, stack(3))
		}
	}()
	bf := backoff.Default()
	tries := job.GetMaxTries()
	if tries <= 0 {
		tries = 1
	}
	for i := 1; i <= tries; i++ {
		if b.Debug {
			fmt.Printf("start handle job[%s] at %d time\n", job.GetName(), i)
		}
		if i != 1 {
			retryInterval := bf.Next()
			if b.Debug {
				fmt.Printf("job[%s] will retry handle after %s\n", job.GetName(), retryInterval)
			}
			time.Sleep(retryInterval)
		}
		err = fn(job)
		if err == nil {
			return nil
		}
		b.Opts.Logger.Errorf("job[%s] handle error: %s %+v", job.GetName(), job.GetPayload(), err)
	}

	if b.Debug {
		fmt.Printf("job[%s] handle failure, reach max try times[%d]\n", job.GetName(), tries)
	}
	return fmt.Errorf("job[%s] handle failure after try max times", job.GetName())
}

// Options struct
type Options struct {
	Name        string
	Driver      string
	Parallel    int
	Debug       bool
	Metric      bool
	Consume     bool
	WaitTimeOut time.Duration
	Logger      *logrus.Logger
}

// Init options
func (opts *Options) Init() {
	if opts.Name == "" {
		opts.Name = defaultQueue
	}

	if opts.Driver == "kafka" { // Name as kafka topic
		opts.Name = strings.Replace(opts.Name, ":", "-", -1)
	}

	if opts.Parallel <= 0 {
		opts.Parallel = 1
	}

	if opts.WaitTimeOut <= 0 {
		opts.WaitTimeOut = defaultTimeout * time.Second
	}

	if opts.Logger == nil {
		opts.Logger = logrus.New()
	}
}

func stack(skip int) []byte {
	buf := new(bytes.Buffer)
	index := 0
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Fprintf(buf, "#%d %s:%d (0x%x)\n", index, file, line, pc)
		index++
		fn := runtime.FuncForPC(pc)
		fnName := []byte("???")
		if fn != nil {
			fnName = []byte(fn.Name())
		}

		fmt.Fprintf(buf, "\t%s: %d\n", fnName, line)
	}
	return buf.Bytes()
}
