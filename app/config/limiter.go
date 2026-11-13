package config

import (
	"fmt"
	"time"

	goresthelper "github.com/elsyahtech/gorest/helper"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func limiterConfig() (*limiter.Config, error) {
	skipFailedRequests, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_RATE_LIMIT_SKIP_FAILED_REQUESTS",
		false,
	)
	if err != nil {
		return &limiter.Config{}, err
	}

	skipSuccessfulRequests, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_RATE_LIMIT_SKIP_SUCCESSFUL_REQUESTS",
		false,
	)
	if err != nil {
		return &limiter.Config{}, err
	}

	disableHeaders, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_RATE_LIMIT_DISABLE_HEADERS",
		false,
	)
	if err != nil {
		return &limiter.Config{}, err
	}

	disableValueRedaction, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_RATE_LIMIT_DISABLE_VALUE_REDACTION",
		false,
	)
	if err != nil {
		return &limiter.Config{}, err
	}

	expiration, err := goresthelper.ReadEnvironmentDurationWithDefault(
		"GOREST_SERVER_RATE_LIMIT_EXPIRATION",
		time.Minute,
	)
	if err != nil {
		return &limiter.Config{}, fmt.Errorf("%w", err)
	}

	max, err := goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_SERVER_RATE_LIMIT_MAX",
		50,
	)
	if err != nil {
		return &limiter.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_RATE_LIMIT_MAX: %w", err)
	}

	return &limiter.Config{
		Max:                    max,
		Expiration:             expiration,
		SkipFailedRequests:     skipFailedRequests,
		SkipSuccessfulRequests: skipSuccessfulRequests,
		DisableHeaders:         disableHeaders,
		DisableValueRedaction:  disableValueRedaction,
	}, nil
}
