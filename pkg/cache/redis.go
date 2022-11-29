// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cache

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v7"
)

// Redis struct for use redis as cache backend
type Redis struct {
	client *redis.Client
	prefix string
}

var _ Cache = (*Redis)(nil)

// NewRedis return redis struct
func NewRedis(client *redis.Client, prefix string) *Redis {
	return &Redis{
		client,
		prefix,
	}
}

// Get return the key of value
func (r *Redis) Get(key string, reply interface{}) (err error) {
	var val []byte
	if val, err = r.client.Get(r.Key(key)).Bytes(); err != nil {
		if err == redis.Nil {
			return Nil
		}
		return
	}
	return json.Unmarshal(val, reply)
}

// Remember use user define callback func to set cache item
// SetCallback should change reply value
func (r *Redis) Remember(key string, reply interface{}, fn SetCallback, expire time.Duration) (err error) {
	if err = r.Get(key, reply); err == nil {
		return nil
	}
	if err != Nil {
		return err
	}

	if err = fn(); err != nil {
		return
	}
	return r.Set(key, reply, expire)
}

// RememberForever use user define callback func to set cache item without expiration
// SetCallback should change reply value
func (r *Redis) RememberForever(key string, reply interface{}, fn SetCallback) (err error) {
	return r.Remember(key, reply, fn, 0)
}

// Set use for set cache item
func (r *Redis) Set(key string, value interface{}, expire time.Duration) (err error) {
	var val []byte
	if val, err = json.Marshal(value); err != nil {
		return
	}
	key = r.Key(key)
	return r.client.Set(key, val, expire).Err()
}

// Forever use for set cache item without expiration
func (r *Redis) Forever(key string, value interface{}) (err error) {
	return r.Set(key, value, 0)
}

// Delete use for delete cache item by key
func (r *Redis) Delete(key string) (err error) {
	return r.client.Del(r.Key(key)).Err()
}

// Key return the key of cache item
func (r *Redis) Key(key string) string {
	if r.prefix == "" {
		return key
	}
	return r.prefix + ":" + key
}
