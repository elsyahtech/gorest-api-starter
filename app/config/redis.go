package config

import (
	"errors"
	"fmt"

	"github.com/elsyahtech/gorest/redis"
)

// RedisConfig loads optional Redis settings; Redis is disabled by default.
func RedisConfig() (*redis.Config, error) {
	cfg := redis.ConfigDefault

	var err error

	cfg.Disable, err = envStrToBool("GOREST_REDIS_DISABLE", true)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_DISABLE: %w", err)
	}

	cfg.Database = envStrToInt("GOREST_REDIS_DATABASE", redis.DefaultRedisDatabase)
	cfg.Host = envStr("GOREST_REDIS_HOST", redis.DefaultRedisHost)
	cfg.Port = envStr("GOREST_REDIS_PORT", redis.DefaultRedisPort)
	cfg.Username = envStr("GOREST_REDIS_USER", "")
	cfg.Password = envStr("GOREST_REDIS_PASSWORD", "")

	cfg.TLS, err = envStrToBool("GOREST_REDIS_TLS", false)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_TLS: %w", err)
	}

	cfg.SSLCertPath = envStr("GOREST_REDIS_CERT_PATH", "")
	cfg.SSLKeyPath = envStr("GOREST_REDIS_KEY_PATH", "")
	if (cfg.SSLCertPath == "") != (cfg.SSLKeyPath == "") {
		return &redis.Config{}, errors.New("GOREST_REDIS_CERT_PATH and GOREST_REDIS_KEY_PATH must be configured together")
	}
	cfg.SSLCAPath = envStr("GOREST_REDIS_CA_PATH", "")

	cfg.Timeout, err = envStrToTimeDuration("GOREST_REDIS_TIMEOUT", redis.DefaultRedisTimeout)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_TIMEOUT: %w", err)
	}

	cfg.PoolTimeout, err = envStrToTimeDuration("GOREST_REDIS_POOL_TIMEOUT", redis.DefaultRedisPoolTimeout)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_POOL_TIMEOUT: %w", err)
	}

	cfg.ConnMaxLifetime, err = envStrToTimeDuration("GOREST_REDIS_CONN_MAX_LIFETIME", redis.DefaultRedisConnMaxLifetime)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_CONN_MAX_LIFETIME: %w", err)
	}

	cfg.ConnMaxIdleTime, err = envStrToTimeDuration("GOREST_REDIS_CONN_MAX_IDLE_TIME", redis.DefaultRedisConnMaxIdleTime)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_CONN_MAX_IDLE_TIME: %w", err)
	}

	cfg.AllowInsecureProduction, err = envStrToBool("GOREST_REDIS_ALLOW_INSECURE_PRODUCTION", false)
	if err != nil {
		return &redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_ALLOW_INSECURE_PRODUCTION: %w", err)
	}

	cfg.MaxRetries = envStrToInt("GOREST_REDIS_MAX_RETRIES", redis.DefaultRedisMaxRetries)
	cfg.PoolSize = envStrToInt("GOREST_REDIS_POOL_SIZE", redis.DefaultRedisPoolSize)
	cfg.MinIdleConns = envStrToInt("GOREST_REDIS_MIN_IDLE_CONNS", redis.DefaultRedisMinIdleConns)
	cfg.MaxActiveConns = envStrToInt("GOREST_REDIS_MAX_ACTIVE_CONNS", redis.DefaultRedisMaxActiveConns)
	cfg.ReadBufferSize = envStrToInt("GOREST_REDIS_READ_BUFFER_SIZE", redis.DefaultRedisReadBufferSize)
	cfg.WriteBufferSize = envStrToInt("GOREST_REDIS_WRITE_BUFFER_SIZE", redis.DefaultRedisWriteBufferSize)

	return &cfg, nil
}
