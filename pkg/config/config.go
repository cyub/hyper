package config

import (
	// viper remote
	"errors"
	"strings"
	"time"

	"github.com/cyub/hyper/logger"
	"github.com/spf13/viper"

	// viper remote config
	_ "github.com/spf13/viper/remote"
)

// Configer interface
type Configer interface {
	GetString(key string, fallback string) string
	GetInt(key string, fallback int) int
	GetBool(key string, fallback bool) bool
	GetIntSlice(key string, fallback []int) []int
	GetStringSlice(key string, fallback []string) []string
	GetStringMap(key string, fallback map[string]interface{}) map[string]interface{}
	GetStringMapString(key string, fallback map[string]string) map[string]string
	GetStringMapStringSlice(key string, fallback map[string][]string) map[string][]string
	GetTime(key string, fallback time.Time) time.Time
	GetDuration(key string, fallback time.Duration) time.Duration
	GetFloat64(key string, fallback float64) float64
	IsSet(key string) bool
}

// Config struct
type Config struct {
	Provider string
	Addr     string
	Path     string
	viper    *viper.Viper
}

var (
	defaultProvider = "consul"
	// ErrInvalidProvider when provider invalid return
	ErrInvalidProvider          = errors.New("config: invalid config center provider")
	_                  Configer = (*Config)(nil)
)

// New return config instance
func New(addr, path string) (*Config, error) {
	var provider string
	addrs := strings.SplitN(addr, "://", 2)
	if len(addrs) < 2 { // localhost:8500
		provider = defaultProvider
	} else if addrs[0] == defaultProvider { // consul://localhost:8500
		provider = defaultProvider
		addr = addrs[1]
	} else if addrs[0] == "etcd" { // etcd://localhost:2379
		provider = "etcd"
		addr = addrs[1]
	} else if addrs[0] == "file" { // file:///etc/hyper/
		provider = "file"
		addr = addrs[1]
	} else {
		return nil, ErrInvalidProvider
	}

	c := &Config{
		Provider: provider,
		Addr:     addr,
		Path:     path,
		viper:    viper.New(),
	}
	err := c.Load()
	return c, err
}

// Load use for load config
func (c *Config) Load() error {
	switch c.Provider {
	case "consul":
		return c.loadFromConsul()
	case "etcd":
		return c.loadFromEtcd()
	case "file":
		return c.loadFromFile()
	default:
		return ErrInvalidProvider
	}
}

/**
 * ReadInConfig 非线程安全
 * https://github.com/spf13/viper/issues/174
 */
func (c *Config) loadFromFile() error {
	logger.Infof("load local file config path[%s] file[config.yml]", c.Addr)
	c.viper.SetConfigName("config")
	c.viper.SetConfigType("yaml")
	c.viper.AddConfigPath(c.Addr)
	if err := c.viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func (c *Config) loadFromEtcd() error {
	logger.Infof("load etcd config addr[%s] path[%s]", c.Addr, c.Path)
	if err := c.viper.AddRemoteProvider("etcd", c.Addr, c.Path); err != nil {
		return err
	}
	c.viper.SetConfigType("yaml")
	if err := c.viper.ReadRemoteConfig(); err != nil {
		return err
	}
	return nil
}

func (c *Config) loadFromConsul() error {
	logger.Infof("load consul config addr[%s] path[%s]", c.Addr, c.Path)
	if err := c.viper.AddRemoteProvider("consul", c.Addr, c.Path); err != nil {
		return err
	}
	c.viper.SetConfigType("yaml")
	if err := c.viper.ReadRemoteConfig(); err != nil {
		return err
	}
	return nil
}

// IsSet use for detected key isset
func (c *Config) IsSet(key string) bool {
	return c.viper.IsSet(key)
}

// GetString use for get string-type value
func (c *Config) GetString(key string, fallback string) string {
	if !c.IsSet(key) {
		return fallback
	}
	return c.viper.GetString(key)
}

// GetInt use for get int-type value
func (c *Config) GetInt(key string, fallback int) int {
	if !c.IsSet(key) {
		return fallback
	}
	return c.viper.GetInt(key)
}

// GetBool use for get bool-type value
func (c *Config) GetBool(key string, fallback bool) bool {
	if !c.IsSet(key) {
		return fallback
	}
	return c.viper.GetBool(key)
}

// GetIntSlice returns the value associated with the key as a slice of int values
func (c *Config) GetIntSlice(key string, fallback []int) []int {
	if !c.IsSet(key) {
		return fallback
	}
	return c.viper.GetIntSlice(key)
}

// GetStringSlice returns the value associated with the key as a slice of strings
func (c *Config) GetStringSlice(key string, fallback []string) []string {
	if !c.IsSet(key) {
		return fallback
	}
	return c.viper.GetStringSlice(key)
}

// GetStringMap returns the value associated with the key as a map of interfaces
func (c *Config) GetStringMap(key string, fallback map[string]interface{}) map[string]interface{} {
	if !c.IsSet(key) {
		return fallback
	}
	return c.viper.GetStringMap(key)
}

// GetStringMapString returns the value associated with the key as a map of strings
func (c *Config) GetStringMapString(key string, fallback map[string]string) map[string]string {
	if !c.IsSet(key) {
		return fallback
	}

	return c.viper.GetStringMapString(key)
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings
func (c *Config) GetStringMapStringSlice(key string, fallback map[string][]string) map[string][]string {
	if !c.IsSet(key) {
		return fallback
	}
	return c.viper.GetStringMapStringSlice(key)
}

// GetTime returns the value associated with the key as time
func (c *Config) GetTime(key string, fallback time.Time) time.Time {
	if !c.IsSet(key) {
		return fallback
	}

	return c.viper.GetTime(key)
}

// GetDuration returns the value associated with the key as a duration
func (c *Config) GetDuration(key string, fallback time.Duration) time.Duration {
	if !c.IsSet(key) {
		return fallback
	}
	return c.viper.GetDuration(key)
}

// GetFloat64 returns the value associated with the key as a float64
func (c *Config) GetFloat64(key string, fallback float64) float64 {
	if !c.IsSet(key) {
		return fallback
	}
	return c.viper.GetFloat64(key)
}
