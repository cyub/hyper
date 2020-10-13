package queue

import (
	"encoding/json"
	"time"
)

// Jober define job interface
type Jober interface {
	GetName() string
	Serialize() ([]byte, error)
	Unserialize() (Job, error)
	GetMaxTries() int
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

// Serialize use for serialize job
func (j *Job) Serialize() (data []byte, err error) {
	data, err = json.Marshal(j)
	return
}

// Unserialize use for unserialize job
func (j *Job) Unserialize() (Job, error) {
	return *j, nil
}

// GetMaxTries return the job max try times
func (j *Job) GetMaxTries() int {
	return j.MaxTries
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
