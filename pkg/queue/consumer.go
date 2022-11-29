// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package queue

import "fmt"

var (
	_consumers = make(map[string]Consumer)
	_debug     = func(job Jober) (err error) {
		var payload interface{}
		job.UnpackPayload(&payload)
		fmt.Printf("the debug job info: name[%s] payload[%#v]\n", job.GetName(), payload)
		return
	}
)

// Consumer define job consumer
type Consumer func(job Jober) error

// RegisterConsumer use Register Queue job consumer
func RegisterConsumer(name string, f Consumer) {
	_consumers[name] = f
}

// Consumers return consumers
func Consumers() map[string]Consumer {
	return _consumers
}

func init() {
	RegisterConsumer("debug", _debug)
}
