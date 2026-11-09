//go:build integration

package integration_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest/redis"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestGORESTRedisServiceAgainstConfiguredInstance(t *testing.T) {
	values, err := godotenv.Read("../.env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("read .env for Redis integration test: %v", err)
	}
	loadRedisEnvironment(t, values)

	if strings.EqualFold(strings.TrimSpace(os.Getenv("GOREST_APP_ENV")), gorest.Production) {
		t.Fatal("Redis integration test refuses to run with GOREST_APP_ENV=production")
	}
	redisConfig, err := config.RedisConfig()
	if err != nil {
		t.Fatalf("load Redis integration configuration: %v", err)
	}
	if redisConfig.Disable {
		t.Skip("Redis is disabled; set GOREST_REDIS_DISABLE=false and configure a non-production test instance to run this test")
	}
	if strings.TrimSpace(redisConfig.Host) == "" || strings.TrimSpace(redisConfig.Port) == "" {
		t.Fatal("GOREST_REDIS_HOST and GOREST_REDIS_PORT are required when Redis integration testing is enabled")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	serviceConfig := redis.Run(*redisConfig)
	service, message, err := redis.New(ctx, &serviceConfig)
	if err != nil {
		t.Fatalf("initialize GOREST Redis service (%s): %v", message, err)
	}
	t.Cleanup(func() {
		if message, err := service.Close(); err != nil {
			t.Errorf("close GOREST Redis service (%s): %v", message, err)
		}
	})
	if err := service.Ping(ctx); err != nil {
		t.Fatalf("GOREST Redis service PING failed: %v", err)
	}
	client, message, status, err := service.Redis()
	if err != nil {
		t.Fatalf("obtain Redis client from GOREST service (%s, HTTP status %d): %v", message, status, err)
	}

	key := "gorest:integration-test:" + uuid.NewString()
	value := uuid.NewString()
	t.Cleanup(func() {
		cleanupContext, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		if err := client.Del(cleanupContext, key).Err(); err != nil {
			t.Errorf("remove Redis integration test key: %v", err)
		}
	})
	if err := client.Set(ctx, key, value, 2*time.Minute).Err(); err != nil {
		t.Fatalf("write unique Redis integration test key: %v", err)
	}
	got, err := client.Get(ctx, key).Result()
	if err != nil {
		t.Fatalf("read Redis integration test key: %v", err)
	}
	if got != value {
		t.Fatalf("Redis integration test value = %q, want the value written by this test", got)
	}
	if err := client.Del(ctx, key).Err(); err != nil {
		t.Fatalf("delete Redis integration test key: %v", err)
	}
}

func loadRedisEnvironment(t *testing.T, values map[string]string) {
	t.Helper()
	loaded := make([]string, 0)
	for key, value := range values {
		if key != "GOREST_APP_ENV" && !strings.HasPrefix(key, "GOREST_REDIS_") {
			continue
		}
		if _, exists := os.LookupEnv(key); exists || value == "" {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf("load %s from .env: %v", key, err)
		}
		loaded = append(loaded, key)
	}
	t.Cleanup(func() {
		for _, key := range loaded {
			if err := os.Unsetenv(key); err != nil {
				t.Errorf("restore environment after Redis integration test (%s): %v", key, err)
			}
		}
	})
}
