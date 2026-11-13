package config

import (
	"fmt"

	goresthelper "github.com/elsyahtech/gorest/helper"
	"github.com/elsyahtech/gorest/server"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

// ServerConfig loads HTTP server, middleware, and listener settings.
// Authentication and authorization must be configured separately for protected routes.
func ServerConfig() (*server.Config, error) {
	cfg := server.ConfigDefault

	var err error

	cfg.Host = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_SERVER_HOST",
		server.DefaultServerHost,
	)

	cfg.Port, err = goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_SERVER_PORT",
		server.DefaultServerPort,
	)
	if err != nil {
		return &server.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_PORT: %w", err)
	}

	stackTrace, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_RECOVER_ENABLE_STACK_TRACE",
		false,
	)
	if err != nil {
		return &server.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_RECOVER_ENABLE_STACK_TRACE: %w", err)
	}

	cfg.Recover = &recover.Config{
		EnableStackTrace: stackTrace,
	}

	cfg.Requestid = &requestid.Config{
		Header: goresthelper.ReadEnvironmentStringWithDefault(
			"GOREST_SERVER_REQUEST_ID_HEADER",
			"",
		),
	}

	compressLevel, err := goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_SERVER_COMPRESS_LEVEL",
		int(compress.LevelDefault),
	)
	if err != nil {
		return &server.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_COMPRESS_LEVEL: %w", err)
	}
	cfg.Compress = &compress.Config{
		Level: compress.Level(compressLevel),
	}

	cfg.Core, err = coreConfig()
	if err != nil {
		return &server.Config{}, fmt.Errorf("failed to configure fiber core: %w", err)
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
