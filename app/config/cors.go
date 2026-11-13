package config

import (
	"fmt"

	goresthelper "github.com/elsyahtech/gorest/helper"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func corsConfig() (*cors.Config, error) {
	origins := goresthelper.ReadEnvironmentStringSliceWithDefault(
		"GOREST_SERVER_CORS_ORIGINS",
		[]string{"http://localhost:3000", "http://127.0.0.1:3000"},
	)
	if len(origins) == 0 {
		origins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}

	methods := goresthelper.ReadEnvironmentStringSliceWithDefault(
		"GOREST_SERVER_CORS_METHODS",
		[]string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	)
	if len(methods) == 0 {
		methods = []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}

	headers := goresthelper.ReadEnvironmentStringSliceWithDefault(
		"GOREST_SERVER_CORS_HEADERS",
		[]string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
	)
	if len(headers) == 0 {
		headers = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"}
	}

	allowCredentials, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_CORS_ALLOW_CREDENTIALS",
		false,
	)
	if err != nil {
		return &cors.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_CORS_ALLOW_CREDENTIALS: %w", err)
	}

	allowPrivateNetwork, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_CORS_ALLOW_PRIVATE_NETWORK",
		false,
	)
	if err != nil {
		return &cors.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_CORS_ALLOW_PRIVATE_NETWORK: %w", err)
	}

	disableCORSValueRedaction, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_CORS_DISABLE_VALUE_REDACTION",
		false,
	)
	if err != nil {
		return &cors.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_CORS_DISABLE_VALUE_REDACTION: %w", err)
	}

	maxAge, err := goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_SERVER_CORS_MAX_AGE",
		300,
	)
	if err != nil {
		return &cors.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_CORS_MAX_AGE: %w", err)
	}

	return &cors.Config{
		AllowOrigins:          origins,
		AllowMethods:          methods,
		AllowHeaders:          headers,
		ExposeHeaders:         goresthelper.ReadEnvironmentStringSliceWithDefault("GOREST_SERVER_CORS_EXPOSE_HEADERS", make([]string, 0)),
		MaxAge:                maxAge,
		AllowCredentials:      allowCredentials,
		AllowPrivateNetwork:   allowPrivateNetwork,
		DisableValueRedaction: disableCORSValueRedaction,
	}, nil
}
