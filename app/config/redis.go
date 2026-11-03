package config

import (
	"errors"
	"os"

	"github.com/elsyahtech/gorest/redis"
)

// Redis configures the Redis connection used by the Redis query examples.
func RedisConfig() (redis.Config, error) {
	cfg := redis.ConfigDefault

	cfg.Disable = strToBool(os.Getenv("GOREST_REDIS_DISABLE"), redis.RedisDisabled)
	cfg.Database = strToInt(os.Getenv("GOREST_REDIS_DATABASE"), redis.DefaultRedisDatabase)
	cfg.Host = str(os.Getenv("GOREST_REDIS_HOST"), redis.DefaultRedisHost)
	cfg.Port = str(os.Getenv("GOREST_REDIS_PORT"), redis.DefaultRedisPort)
	cfg.Username = str(os.Getenv("GOREST_REDIS_USER"), "")
	cfg.Password = str(os.Getenv("GOREST_REDIS_PASSWORD"), "")
	cfg.TLS = strToBool(os.Getenv("GOREST_REDIS_TLS"), false)

	cfg.SSLCertPath = str(os.Getenv("GOREST_REDIS_CERT_PATH"), "")
	if cfg.SSLCertPath == "" {
		if cfg.TLS {
			return redis.Config{}, errors.New("GOREST_REDIS_CERT_PATH cannot be empty when GOREST_REDIS_TLS is true")
		}
	}

	cfg.SSLKeyPath = str(os.Getenv("GOREST_REDIS_KEY_PATH"), "")
	cfg.SSLCAPath = str(os.Getenv("GOREST_REDIS_CA_PATH"), "")
	cfg.Timeout = strToTimeDuration(os.Getenv("GOREST_REDIS_TIMEOUT"), redis.DefaultRedisTimeout)
	cfg.PoolTimeout = strToTimeDuration(os.Getenv("GOREST_REDIS_POOL_TIMEOUT"), redis.DefaultRedisPoolTimeout)
	cfg.ConnMaxLifetime = strToTimeDuration(os.Getenv("GOREST_REDIS_CONN_MAX_LIFETIME"), redis.DefaultRedisConnMaxLifetime)
	cfg.ConnMaxIdleTime = strToTimeDuration(os.Getenv("GOREST_REDIS_CONN_MAX_IDLE_TIME"), redis.DefaultRedisConnMaxIdleTime)
	cfg.AllowInsecureProduction = strToBool(os.Getenv("GOREST_REDIS_ALLOW_INSECURE_PRODUCTION"), false)
	cfg.MaxRetries = strToInt(os.Getenv("GOREST_REDIS_MAX_RETRIES"), redis.DefaultRedisMaxRetries)
	cfg.PoolSize = strToInt(os.Getenv("GOREST_REDIS_POOL_SIZE"), redis.DefaultRedisPoolSize)
	cfg.MinIdleConns = strToInt(os.Getenv("GOREST_REDIS_MIN_IDLE_CONNS"), redis.DefaultRedisMinIdleConns)
	cfg.MaxActiveConns = strToInt(os.Getenv("GOREST_REDIS_MAX_ACTIVE_CONNS"), redis.DefaultRedisMaxActiveConns)
	cfg.ReadBufferSize = strToInt(os.Getenv("GOREST_REDIS_READ_BUFFER_SIZE"), redis.DefaultRedisReadBufferSize)
	cfg.WriteBufferSize = strToInt(os.Getenv("GOREST_REDIS_WRITE_BUFFER_SIZE"), redis.DefaultRedisWriteBufferSize)

	return cfg, nil
}
