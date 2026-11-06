package config

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func limiterConfig() (*limiter.Config, error) {
	skipFailedRequests, err := envStrToBool("GOREST_SERVER_RATE_LIMIT_SKIP_FAILED_REQUESTS", false)
	if err != nil {
		return &limiter.Config{}, err
	}

	skipSuccessfulRequests, err := envStrToBool("GOREST_SERVER_RATE_LIMIT_SKIP_SUCCESSFUL_REQUESTS", false)
	if err != nil {
		return &limiter.Config{}, err
	}

	disableHeaders, err := envStrToBool("GOREST_SERVER_RATE_LIMIT_DISABLE_HEADERS", false)
	if err != nil {
		return &limiter.Config{}, err
	}

	disableValueRedaction, err := envStrToBool("GOREST_SERVER_RATE_LIMIT_DISABLE_VALUE_REDACTION", false)
	if err != nil {
		return &limiter.Config{}, err
	}

	expiration, err := envStrToTimeDuration("GOREST_SERVER_RATE_LIMIT_EXPIRATION", time.Minute)
	if err != nil {
		return &limiter.Config{}, fmt.Errorf("%w", err)
	}

	return &limiter.Config{
		Max:                    envStrToInt("GOREST_SERVER_RATE_LIMIT_MAX", 50),
		Expiration:             expiration,
		SkipFailedRequests:     skipFailedRequests,
		SkipSuccessfulRequests: skipSuccessfulRequests,
		DisableHeaders:         disableHeaders,
		DisableValueRedaction:  disableValueRedaction,
	}, nil
}
