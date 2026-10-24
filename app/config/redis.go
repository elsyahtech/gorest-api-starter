package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elsyahtech/gorest/redis"
)

// Redis configures the Redis connection used by the Redis query examples.
func Redis() redis.Config {
	cfg, err := LoadRedis()
	if err != nil {
		panic(err)
	}

	return cfg
}

// LoadRedis reads optional Redis connection settings from the environment.
func LoadRedis() (redis.Config, error) {
	if err := validateEnvironment(); err != nil {
		return redis.Config{}, err
	}

	disable, err := envBool("GOREST_REDIS_DISABLE", true)
	if err != nil {
		return redis.Config{}, err
	}

	tls, err := envBool("GOREST_REDIS_TLS", false)
	if err != nil {
		return redis.Config{}, err
	}

	databaseNumber := 0

	if value, ok := os.LookupEnv("GOREST_REDIS_DATABASE"); ok {
		databaseNumber, err = strconv.Atoi(strings.TrimSpace(value))

		if err != nil || databaseNumber < 0 {
			return redis.Config{}, fmt.Errorf("GOREST_REDIS_DATABASE must be a non-negative integer")
		}
	}

	redisPort := envString("GOREST_REDIS_PORT", "6379")

	if value, ok := os.LookupEnv("GOREST_REDIS_PORT"); ok {
		port, parseErr := strconv.Atoi(strings.TrimSpace(value))

		if parseErr != nil || port < 1 || port > 65535 {
			return redis.Config{}, fmt.Errorf("GOREST_REDIS_PORT must be an integer between 1 and 65535")
		}
	}

	return redis.Config{
		// Enable Redis only when the application needs it. If enabled, use a
		// private network and configure authentication; enable TLS when traffic
		// crosses a network that is not otherwise trusted.
		Disable: disable,

		// For local development. In containers, use the Redis service name or
		// reachable private host instead of 127.0.0.1.
		Host: envString("GOREST_REDIS_HOST", "127.0.0.1"),
		Port: redisPort,

		// Supply credentials when Redis authentication is enabled. Avoid
		// committing production credentials to source control.
		Username: envString("GOREST_REDIS_USERNAME", ""),
		Password: envSecret("GOREST_REDIS_PASSWORD", ""),

		// Set certificate paths when TLS or mutual TLS is configured.
		SSLCAPath:   envString("GOREST_REDIS_SSL_CA_PATH", ""),
		SSLCertPath: envString("GOREST_REDIS_SSL_CERT_PATH", ""),
		SSLKeyPath:  envString("GOREST_REDIS_SSL_KEY_PATH", ""),

		// Connection and pool timeouts.
		Timeout:         5 * time.Second,
		PoolTimeout:     6 * time.Second,
		ConnMaxLifetime: 20 * time.Minute,
		ConnMaxIdleTime: 30 * time.Minute,

		TLS:        tls,
		Database:   databaseNumber,
		MaxRetries: 3,

		// Keep the GOREST defaults unless measured traffic requires tuning.
		PoolSize:        0,
		MinIdleConns:    5,
		MaxActiveConns:  0,
		ReadBufferSize:  32 * 1024,
		WriteBufferSize: 32 * 1024,
	}, nil
}
