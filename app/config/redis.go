package config

import (
	"errors"
	"fmt"

	"github.com/elsyahtech/gorest/redis"
)

// RedisConfig loads optional Redis settings; Redis is disabled by default.
func RedisConfig() (redis.Config, error) {
	cfg := redis.ConfigDefault

	var err error

	cfg.Disable, err = strToBool("GOREST_REDIS_DISABLE", true)
	if err != nil {
		return redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_DISABLE: %w", err)
	}

	cfg.Database = strToInt("GOREST_REDIS_DATABASE", redis.DefaultRedisDatabase)
	cfg.Host = str("GOREST_REDIS_HOST", redis.DefaultRedisHost)
	cfg.Port = str("GOREST_REDIS_PORT", redis.DefaultRedisPort)
	cfg.Username = str("GOREST_REDIS_USER", "")
	cfg.Password = str("GOREST_REDIS_PASSWORD", "")

	cfg.TLS, err = strToBool("GOREST_REDIS_TLS", false)
	if err != nil {
		return redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_TLS: %w", err)
	}

	cfg.SSLCertPath = str("GOREST_REDIS_CERT_PATH", "")
	cfg.SSLKeyPath = str("GOREST_REDIS_KEY_PATH", "")
	if (cfg.SSLCertPath == "") != (cfg.SSLKeyPath == "") {
		return redis.Config{}, errors.New("GOREST_REDIS_CERT_PATH and GOREST_REDIS_KEY_PATH must be configured together")
	}
	cfg.SSLCAPath = str("GOREST_REDIS_CA_PATH", "")
	cfg.Timeout = strToTimeDuration("GOREST_REDIS_TIMEOUT", redis.DefaultRedisTimeout)
	cfg.PoolTimeout = strToTimeDuration("GOREST_REDIS_POOL_TIMEOUT", redis.DefaultRedisPoolTimeout)
	cfg.ConnMaxLifetime = strToTimeDuration("GOREST_REDIS_CONN_MAX_LIFETIME", redis.DefaultRedisConnMaxLifetime)
	cfg.ConnMaxIdleTime = strToTimeDuration("GOREST_REDIS_CONN_MAX_IDLE_TIME", redis.DefaultRedisConnMaxIdleTime)

	cfg.AllowInsecureProduction, err = strToBool("GOREST_REDIS_ALLOW_INSECURE_PRODUCTION", false)
	if err != nil {
		return redis.Config{}, fmt.Errorf("failed to parse GOREST_REDIS_ALLOW_INSECURE_PRODUCTION: %w", err)
	}

	cfg.MaxRetries = strToInt("GOREST_REDIS_MAX_RETRIES", redis.DefaultRedisMaxRetries)
	cfg.PoolSize = strToInt("GOREST_REDIS_POOL_SIZE", redis.DefaultRedisPoolSize)
	cfg.MinIdleConns = strToInt("GOREST_REDIS_MIN_IDLE_CONNS", redis.DefaultRedisMinIdleConns)
	cfg.MaxActiveConns = strToInt("GOREST_REDIS_MAX_ACTIVE_CONNS", redis.DefaultRedisMaxActiveConns)
	cfg.ReadBufferSize = strToInt("GOREST_REDIS_READ_BUFFER_SIZE", redis.DefaultRedisReadBufferSize)
	cfg.WriteBufferSize = strToInt("GOREST_REDIS_WRITE_BUFFER_SIZE", redis.DefaultRedisWriteBufferSize)

	return cfg, nil
}
