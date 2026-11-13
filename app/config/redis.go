package config

import (
	"errors"
	"fmt"

	goresthelper "github.com/elsyahtech/gorest/helper"
	"github.com/elsyahtech/gorest/redis"
)

// RedisConfig loads optional Redis settings; Redis is disabled by default.
func RedisConfig() (*redis.Config, error) {
	cfg := redis.ConfigDefault

	var err error

	cfg.Disable, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_REDIS_DISABLE",
		true,
	)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_DISABLE: %w", err)
	}

	cfg.Database, err = goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_REDIS_DATABASE",
		redis.DefaultRedisDatabase,
	)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_DATABASE: %w", err)
	}

	cfg.Host = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_REDIS_HOST",
		redis.DefaultRedisHost,
	)

	cfg.Port = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_REDIS_PORT",
		redis.DefaultRedisPort,
	)

	cfg.Username = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_REDIS_USER",
		"",
	)

	cfg.Password = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_REDIS_PASSWORD",
		"",
	)

	cfg.TLS, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_REDIS_TLS",
		false,
	)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_TLS: %w", err)
	}

	cfg.SSLCertPath = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_REDIS_CERT_PATH",
		"",
	)

	cfg.SSLKeyPath = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_REDIS_KEY_PATH",
		"",
	)

	if (cfg.SSLCertPath == "") != (cfg.SSLKeyPath == "") {
		return &redis.Config{}, errors.New("GOREST_REDIS_CERT_PATH and GOREST_REDIS_KEY_PATH must be configured together")
	}

	cfg.SSLCAPath = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_REDIS_CA_PATH",
		"",
	)

	cfg.Timeout, err = goresthelper.ReadEnvironmentDurationWithDefault(
		"GOREST_REDIS_TIMEOUT",
		redis.DefaultRedisTimeout,
	)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_TIMEOUT: %w", err)
	}

	cfg.PoolTimeout, err = goresthelper.ReadEnvironmentDurationWithDefault(
		"GOREST_REDIS_POOL_TIMEOUT",
		redis.DefaultRedisPoolTimeout,
	)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_POOL_TIMEOUT: %w", err)
	}

	cfg.ConnMaxLifetime, err = goresthelper.ReadEnvironmentDurationWithDefault(
		"GOREST_REDIS_CONN_MAX_LIFETIME",
		redis.DefaultRedisConnMaxLifetime,
	)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_CONN_MAX_LIFETIME: %w", err)
	}

	cfg.ConnMaxIdleTime, err = goresthelper.ReadEnvironmentDurationWithDefault(
		"GOREST_REDIS_CONN_MAX_IDLE_TIME",
		redis.DefaultRedisConnMaxIdleTime,
	)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_CONN_MAX_IDLE_TIME: %w", err)
	}

	cfg.AllowInsecureProduction, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_REDIS_ALLOW_INSECURE_PRODUCTION",
		false,
	)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_ALLOW_INSECURE_PRODUCTION: %w", err)
	}

	cfg.MaxRetries, err = goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_REDIS_MAX_RETRIES",
		redis.DefaultRedisMaxRetries,
	)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_MAX_RETRIES: %w", err)
	}

	cfg.PoolSize, err = goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_REDIS_POOL_SIZE",
		redis.DefaultRedisPoolSize,
	)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_POOL_SIZE: %w", err)
	}

	cfg.MinIdleConns, err = goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_REDIS_MIN_IDLE_CONNS",
		redis.DefaultRedisMinIdleConns,
	)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_MIN_IDLE_CONNS: %w", err)
	}

	cfg.MaxActiveConns, err = goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_REDIS_MAX_ACTIVE_CONNS",
		redis.DefaultRedisMaxActiveConns,
	)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_MAX_ACTIVE_CONNS: %w", err)
	}

	cfg.ReadBufferSize, err = goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_REDIS_READ_BUFFER_SIZE",
		redis.DefaultRedisReadBufferSize,
	)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_READ_BUFFER_SIZE: %w", err)
	}

	cfg.WriteBufferSize, err = goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_REDIS_WRITE_BUFFER_SIZE",
		redis.DefaultRedisWriteBufferSize,
	)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_WRITE_BUFFER_SIZE: %w", err)
	}

	return &cfg, nil
}
