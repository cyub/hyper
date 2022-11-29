// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"time"

	"github.com/cyub/hyper/pkg/config"
)

var c *config.Config

// Init for init config
func Init(addr, path string) (err error) {
	c, err = config.New(addr, path)
	return err
}

// Instance return the instance of Config
func Instance() *config.Config {
	return c
}

// IsSet use for detected key isset
func IsSet(key string) bool {
	return c.IsSet(key)
}

// GetString use for get string-type value
func GetString(key string, fallback string) string {
	return c.GetString(key, fallback)
}

// GetInt use for get int-type value
func GetInt(key string, fallback int) int {
	return c.GetInt(key, fallback)
}

// GetBool use for get bool-type value
func GetBool(key string, fallback bool) bool {
	return c.GetBool(key, fallback)
}

// GetIntSlice exporter
func GetIntSlice(key string, fallback []int) []int {
	return c.GetIntSlice(key, fallback)
}

// GetStringSlice exporter
func GetStringSlice(key string, fallback []string) []string {
	return c.GetStringSlice(key, fallback)
}

// GetStringMap exporter
func GetStringMap(key string, fallback map[string]interface{}) map[string]interface{} {
	return c.GetStringMap(key, fallback)
}

// GetStringMapString exporter
func GetStringMapString(key string, fallback map[string]string) map[string]string {
	return c.GetStringMapString(key, fallback)
}

// GetStringMapStringSlice exporter
func GetStringMapStringSlice(key string, fallback map[string][]string) map[string][]string {
	return c.GetStringMapStringSlice(key, fallback)
}

// GetTime exporter
func GetTime(key string, fallback time.Time) time.Time {
	return c.GetTime(key, fallback)
}

// GetDuration exporter
func GetDuration(key string, fallback time.Duration) time.Duration {
	return c.GetDuration(key, fallback)
}

// GetFloat64 exporter
func GetFloat64(key string, fallback float64) float64 {
	return c.GetFloat64(key, fallback)
}
