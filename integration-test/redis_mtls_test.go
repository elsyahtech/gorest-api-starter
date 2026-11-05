//go:build integration

package integration_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	gorestredis "github.com/elsyahtech/gorest/redis"
	"github.com/joho/godotenv"
)

func TestGORESTRedisMutualTLSRequiresClientCertificate(t *testing.T) {
	values, err := godotenv.Read("../.env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("read .env for Redis mTLS test: %v", err)
	}
	value := func(key string) string {
		if current := os.Getenv(key); current != "" {
			return current
		}
		return values[key]
	}
	setting := func(testKey, runtimeKey string) string {
		if configured := value(testKey); configured != "" {
			return configured
		}
		return value(runtimeKey)
	}
	if strings.EqualFold(strings.TrimSpace(value("GOREST_APP_ENV")), "production") {
		t.Fatal("Redis mTLS integration test refuses to run with GOREST_APP_ENV=production")
	}

	testKeys := []string{
		"GOREST_TEST_REDIS_MTLS_HOST",
		"GOREST_TEST_REDIS_MTLS_PORT",
		"GOREST_TEST_REDIS_MTLS_DATABASE",
		"GOREST_TEST_REDIS_MTLS_USER",
		"GOREST_TEST_REDIS_MTLS_PASSWORD",
		"GOREST_TEST_REDIS_MTLS_CA_FILE",
		"GOREST_TEST_REDIS_MTLS_CLIENT_CERT_FILE",
		"GOREST_TEST_REDIS_MTLS_CLIENT_KEY_FILE",
	}
	hasDedicatedSettings := false
	for _, key := range testKeys {
		if value(key) != "" {
			hasDedicatedSettings = true
		}
	}
	certPath := setting("GOREST_TEST_REDIS_MTLS_CLIENT_CERT_FILE", "GOREST_REDIS_CERT_PATH")
	keyPath := setting("GOREST_TEST_REDIS_MTLS_CLIENT_KEY_FILE", "GOREST_REDIS_KEY_PATH")
	if !hasDedicatedSettings && certPath == "" && keyPath == "" {
		t.Skip("set GOREST_TEST_REDIS_MTLS_* to test a Redis server that requires a client certificate")
	}
	if !hasDedicatedSettings && !strings.EqualFold(strings.TrimSpace(value("GOREST_REDIS_TLS")), "true") {
		t.Fatal("GOREST_REDIS_TLS=true is required when the Redis mTLS test uses runtime GOREST_REDIS_* settings")
	}
	if (certPath == "") != (keyPath == "") {
		t.Fatal("Redis mTLS requires both a client certificate and key")
	}
	host := strings.TrimSpace(setting("GOREST_TEST_REDIS_MTLS_HOST", "GOREST_REDIS_HOST"))
	if host == "" {
		t.Fatal("GOREST_TEST_REDIS_MTLS_HOST or GOREST_REDIS_HOST is required")
	}
	caPath := setting("GOREST_TEST_REDIS_MTLS_CA_FILE", "GOREST_REDIS_CA_PATH")
	if caPath == "" {
		t.Fatal("Redis mTLS requires a CA file to verify the server certificate")
	}
	caPEM, err := os.ReadFile(caPath)
	if err != nil {
		t.Fatalf("read Redis mTLS CA certificate: %v", err)
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(caPEM) {
		t.Fatal("Redis mTLS CA file contains no valid PEM certificates")
	}
	if certPath == "" {
		t.Fatal("Redis mTLS client certificate and key must be configured")
	}
	if _, err := tls.LoadX509KeyPair(certPath, keyPath); err != nil {
		t.Fatalf("load Redis mTLS client certificate and key: %v", err)
	}

	config := gorestredis.ConfigDefault
	config.Disable = false
	config.Host = host
	config.Port = strings.TrimSpace(setting("GOREST_TEST_REDIS_MTLS_PORT", "GOREST_REDIS_PORT"))
	if config.Port == "" {
		config.Port = "6379"
	}
	if database := strings.TrimSpace(setting("GOREST_TEST_REDIS_MTLS_DATABASE", "GOREST_REDIS_DATABASE")); database != "" {
		config.Database, err = strconv.Atoi(database)
		if err != nil || config.Database < 0 {
			t.Fatalf("Redis mTLS database index must be a non-negative integer: %v", err)
		}
	}
	config.Username = setting("GOREST_TEST_REDIS_MTLS_USER", "GOREST_REDIS_USER")
	config.Password = setting("GOREST_TEST_REDIS_MTLS_PASSWORD", "GOREST_REDIS_PASSWORD")
	config.TLS = true
	config.SSLCAPath = caPath

	newService := func(config gorestredis.Config) (*gorestredis.Redis, string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		serviceConfig := gorestredis.Run(config)
		return gorestredis.New(ctx, &serviceConfig)
	}
	service, message, err := newService(config)
	if err == nil {
		if message, closeErr := service.Close(); closeErr != nil {
			t.Errorf("close Redis connection without a client certificate (%s): %v", message, closeErr)
		}
		t.Fatal("Redis accepted a connection without the client certificate; configure the test endpoint to require mTLS")
	}
	config.SSLCertPath = certPath
	config.SSLKeyPath = keyPath
	service, message, err = newService(config)
	if err != nil {
		t.Fatalf("connect to Redis with the configured client certificate (%s): %v", message, err)
	}
	t.Cleanup(func() {
		if message, err := service.Close(); err != nil {
			t.Errorf("close Redis mTLS service (%s): %v", message, err)
		}
	})
	pingContext, cancelPing := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelPing()
	if err := service.Ping(pingContext); err != nil {
		t.Fatalf("Redis mTLS PING failed: %v", err)
	}
}
