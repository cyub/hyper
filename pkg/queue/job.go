// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package queue

import (
	"encoding/json"
	"time"
)

// Jober define job interface
type Jober interface {
	GetName() string
	GetMaxTries() int
	Serialize() ([]byte, error)
	Unserialize() (Job, error)
	PackPayload() ([]byte, error)
	UnpackPayload(interface{}) error
	GetPayload() []byte
}

// Job struct
type Job struct {
	Name     string
	Payload  []byte
	MaxTries int
	Timeout  time.Duration
}

// GetName return the name of job
func (j *Job) GetName() string {
	return j.Name
}

// GetMaxTries return the max tries time of job
func (j *Job) GetMaxTries() int {
	return j.MaxTries
}

// GetPayload return the payload of job
func (j *Job) GetPayload() []byte {
	return j.Payload
}

// Serialize use for serialize job
func (j *Job) Serialize() (data []byte, err error) {
	data, err = json.Marshal(j)
	return
}

// Unserialize use for unserialize job
func (j *Job) Unserialize() (Job, error) {
	return *j, nil
}

// PackPayload use for pack the data of job
func (j *Job) PackPayload() (packed []byte, err error) {
	packed, err = json.Marshal(j.Payload)
	return
}

// UnpackPayload use for unpack the data of job
func (j *Job) UnpackPayload(v interface{}) error {
	return json.Unmarshal(j.Payload, v)
}
