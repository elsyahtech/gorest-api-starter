package config

import (
	"time"

	"github.com/elsyahtech/gorest/redis"
)

// Redis configures the optional Redis connection. The starter does not require
// Redis, so it remains disabled until the application uses it.
func Redis() redis.Config {
	return redis.Config{
		// Enable Redis only when the application needs it. If enabled, use a
		// private network and configure authentication; enable TLS when traffic
		// crosses a network that is not otherwise trusted.
		Enabled: false,

		// For local development. In containers, use the Redis service name or
		// reachable private host instead of 127.0.0.1.
		Host: "127.0.0.1",
		Port: "6379",

		// Supply credentials when Redis authentication is enabled. Avoid
		// committing production credentials to source control.
		Username: "",
		Password: "",

		// Set certificate paths when TLS or mutual TLS is configured.
		SSLCAPath:   "",
		SSLCertPath: "",
		SSLKeyPath:  "",

		// Connection and pool timeouts.
		Timeout:         5 * time.Second,
		PoolTimeout:     6 * time.Second,
		ConnMaxLifetime: 20 * time.Minute,
		ConnMaxIdleTime: 30 * time.Minute,

		TLS:        false,
		Database:   0,
		MaxRetries: 3,

		// Keep the GOREST defaults unless measured traffic requires tuning.
		PoolSize:        0,
		MinIdleConns:    5,
		MaxActiveConns:  0,
		ReadBufferSize:  32 * 1024,
		WriteBufferSize: 32 * 1024,
	}
}
