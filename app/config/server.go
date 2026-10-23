package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elsyahtech/gorest/server"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

// Server configures the HTTP server and its middleware.
// These middleware provide request handling controls; applications must still
// implement authentication, authorization, and any required CSRF protections.
// Configure TLS at a trusted ingress/reverse proxy or in the server deployment.
func Server() server.Config {
	cfg, err := LoadServer()
	if err != nil {
		panic(err)
	}
	return cfg
}

// LoadServer applies deployment-specific bind and CORS settings. In production,
// bind settings must be explicit so the local loopback example is not reused.
func LoadServer() (server.Config, error) {
	if err := validateEnvironment(); err != nil {
		return server.Config{}, err
	}
	port := 3000
	if value, ok := os.LookupEnv("GOREST_SERVER_PORT"); ok {
		parsed, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil || parsed < 1 || parsed > 65535 {
			return server.Config{}, fmt.Errorf("GOREST_SERVER_PORT must be an integer between 1 and 65535")
		}
		port = parsed
	}
	host := envString("GOREST_SERVER_HOST", "127.0.0.1")
	origins := csvEnv("GOREST_CORS_ORIGINS")
	if len(origins) == 0 {
		origins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}
	if Environment() == "production" {
		if err := requireProductionEnv("GOREST_SERVER_HOST"); err != nil {
			return server.Config{}, err
		}
		if err := requireProductionEnv("GOREST_SERVER_PORT"); err != nil {
			return server.Config{}, err
		}
		if err := requireProductionEnv("GOREST_CORS_ORIGINS"); err != nil {
			return server.Config{}, err
		}
		if len(csvEnv("GOREST_CORS_ORIGINS")) == 0 {
			return server.Config{}, fmt.Errorf("GOREST_CORS_ORIGINS must contain at least one origin when APP_ENV=production")
		}
	}
	return server.Config{
		// Keep request timeouts finite and cap the body size. These timeout values
		// match GOREST defaults; adjust them for the application's workload.
		Core: &fiber.Config{
			BodyLimit:    4 * 1024 * 1024, // 4 MiB
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		},

		Cors: &cors.Config{
			AllowOrigins: origins,
			AllowMethods: []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
			MaxAge:       300,
		},

		// Recover converts handler panics into HTTP errors. It does not replace
		// input validation or application error handling.
		Recover: &recover.ConfigDefault,

		// Attach a request ID to make logs and request troubleshooting easier.
		Requestid: &requestid.ConfigDefault,

		// Compress eligible responses to reduce response size.
		Compress: &compress.ConfigDefault,

		// Add common HTTP security headers. These headers do not replace
		// authentication, authorization, or secure deployment configuration.
		Helmet: &helmet.ConfigDefault,

		// The default limiter uses in-memory storage, so limits apply per process.
		// Use shared storage or an upstream gateway if limits must span instances.
		// Review proxy trust settings before relying on client IP behind a proxy.
		Limiter: &limiter.Config{
			Max:        50,
			Expiration: time.Minute,
		},

		// Local-only by default. In a container, bind to 0.0.0.0 only when the
		// deployment network/firewall controls which clients can reach the service.
		Host: host,
		Port: port,
	}, nil
}

func csvEnv(key string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return []string{}
	}
	items := strings.Split(value, ",")
	result := make([]string, 0, len(items))
	for _, item := range items {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
