package queue

import "github.com/cyub/hyper/logger"

var (
	_consumers = make(map[string]Consumer)
	_debug     = func(job *Job) (err error) {
		logger.Debugf("the debug job info: %+v\n", job)
		return
	}
)

// RegisterConsumer use Register Queue job consumer
func RegisterConsumer(name string, f Consumer) {
	_consumers[name] = f
}

func init() {
	RegisterConsumer("debug", _debug)
}
