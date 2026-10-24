package main

import (
	"testing"

	"github.com/elsyahtech/gorest/redis"
)

func TestRequireProductionRedis(t *testing.T) {
	t.Run("allows disabled Redis", func(t *testing.T) {
		setRedisProductionEnvironment(t, "", "")
		if err := requireProductionRedis(redis.Config{Disable: true}); err != nil {
			t.Fatalf("requireProductionRedis() error = %v", err)
		}
	})

	t.Run("requires credentials and TLS when enabled", func(t *testing.T) {
		for name, cfg := range map[string]redis.Config{
			"host":     {Disable: false, Password: "secret", TLS: true},
			"password": {Disable: false, Host: "redis.example.internal", TLS: true},
			"TLS":      {Disable: false, Host: "redis.example.internal", Password: "secret"},
		} {
			t.Run(name, func(t *testing.T) {
				host, password := "redis.example.internal", "secret"
				if name == "host" {
					host = ""
				}
				if name == "password" {
					password = ""
				}
				setRedisProductionEnvironment(t, host, password)
				if err := requireProductionRedis(cfg); err == nil {
					t.Fatalf("requireProductionRedis() succeeded without %s", name)
				}
			})
		}
	})

	t.Run("accepts enabled Redis with host password and TLS", func(t *testing.T) {
		setRedisProductionEnvironment(t, "redis.example.internal", "test-only-password")
		cfg := redis.Config{
			Disable:  false,
			Host:     "redis.example.internal",
			Password: "test-only-password",
			TLS:      true,
		}
		if err := requireProductionRedis(cfg); err != nil {
			t.Fatalf("requireProductionRedis() error = %v", err)
		}
	})
}

func setRedisProductionEnvironment(t *testing.T, host, password string) {
	t.Helper()
	t.Setenv("GOREST_REDIS_HOST", host)
	t.Setenv("GOREST_REDIS_PASSWORD", password)
}
