package config

import (
	"errors"

	"github.com/cyub/hyper/logger"
	"github.com/spf13/viper"

	// viper remote config
	_ "github.com/spf13/viper/remote"
)

// Config struct define consul config
type Config struct {
	Provider string
	Addr     string
	Path     string
	viper    *viper.Viper
}

var (
	defaultProvider = "consul"
	// ErrInvalidProvider when provider invalid return
	ErrInvalidProvider = errors.New("config: invalid config center provider")
)

// Init for init config
func Init(provider, addr, key string) error {
	c = &Config{
		Provider: provider,
		Addr:     addr,
		Path:     key,
		viper:    viper.New(),
	}

	switch provider {
	case "consul":
		return c.loadFromConsul()
	default:
		return ErrInvalidProvider
	}
}

func (c *Config) loadFromConsul() error {
	logger.Infof("consul connect info addr[%s] path[%s]", c.Addr, c.Path)
	if err := c.viper.AddRemoteProvider("consul", c.Addr, c.Path); err != nil {
		return err
	}
	c.viper.SetConfigType("yaml")
	if err := c.viper.ReadRemoteConfig(); err != nil {
		return err
	}

	return nil
}

// GetString use for get string-type value
func (c *Config) GetString(key string, fallback string) string {
	if !IsSet(key) {
		return fallback
	}
	return c.viper.GetString(key)
}

// GetInt use for get int-type value
func (c *Config) GetInt(key string, fallback int) int {
	if !IsSet(key) {
		return fallback
	}
	return c.viper.GetInt(key)
}

// GetBool use for get bool-type value
func (c *Config) GetBool(key string, fallback bool) bool {
	if !IsSet(key) {
		return fallback
	}
	return c.viper.GetBool(key)
}

// IsSet use for detected key isset
func (c *Config) IsSet(key string) bool {
	return c.viper.IsSet(key)
}
