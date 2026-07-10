package config

import "fmt"

// RedisAddr returns the host:port address for the Redis instance.
func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}
