package config

var c *Config

// Instance return the instance of Config
func Instance() *Config {
	return c
}

// GetString use for get string-type value
func GetString(key string, fallback string) string {
	if !IsSet(key) {
		return fallback
	}
	return c.viper.GetString(key)
}

// GetInt use for get int-type value
func GetInt(key string, fallback int) int {
	if !IsSet(key) {
		return fallback
	}
	return c.viper.GetInt(key)
}

// GetBool use for get bool-type value
func GetBool(key string, fallback bool) bool {
	if !IsSet(key) {
		return fallback
	}
	return c.viper.GetBool(key)
}

// IsSet use for detected key isset
func IsSet(key string) bool {
	return c.viper.IsSet(key)
}
