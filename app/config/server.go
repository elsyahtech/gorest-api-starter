package config

import (
	"fmt"
	"os"

	"github.com/elsyahtech/gorest/server"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

// ServerConfig loads HTTP server, middleware, and listener settings.
// Authentication and authorization must be configured separately for protected routes.
func ServerConfig() (*server.Config, error) {
	cfg := server.ConfigDefault

	cfg.Host = envStr("GOREST_SERVER_HOST", server.DefaultServerHost)
	cfg.Port = envStrToInt("GOREST_SERVER_PORT", server.DefaultServerPort)

	stackTrace, err := envStrToBool("GOREST_SERVER_RECOVER_ENABLE_STACK_TRACE", false)
	if err != nil {
		return &server.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_RECOVER_ENABLE_STACK_TRACE: %w", err)
	}

	cfg.Recover = &recover.Config{
		EnableStackTrace: stackTrace,
	}

	cfg.Requestid = &requestid.Config{
		Header: os.Getenv("GOREST_SERVER_REQUEST_ID_HEADER"),
	}

	cfg.Compress = &compress.Config{
		Level: compress.Level(envStrToInt("GOREST_SERVER_COMPRESS_LEVEL", int(compress.LevelDefault))),
	}

	cfg.Core, err = coreConfig()
	if err != nil {
		return &server.Config{}, fmt.Errorf("failed to configure fiber core: %w", err)
	}

	origins := envToStrArray("GOREST_SERVER_CORS_ORIGINS", []string{"http://localhost:3000", "http://127.0.0.1:3000"})
	if len(origins) == 0 {
		origins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}

	methods := envToStrArray("GOREST_SERVER_CORS_METHODS", []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	if len(methods) == 0 {
		methods = []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}

	headers := envToStrArray("GOREST_SERVER_CORS_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"})
	if len(headers) == 0 {
		headers = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"}
	}

	cfg.Cors, err = corsConfig()
	if err != nil {
		return &server.Config{}, fmt.Errorf("failed to configure CORS: %w", err)
	}

	cfg.Helmet, err = helmetConfig()
	if err != nil {
		return &server.Config{}, fmt.Errorf("failed to configure Helmet: %w", err)
	}

	cfg.Limiter, err = limiterConfig()
	if err != nil {
		return &server.Config{}, fmt.Errorf("failed to configure rate limiter: %w", err)
	}

	listener, err := listenerConfig()
	if err != nil {
		return &server.Config{}, err
	}

	cfg.Listener = *listener

	return &cfg, nil
}
