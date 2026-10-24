package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/elsyahtech/gorest"
)

// Environment returns the deployment environment. Development keeps the
// starter's local defaults when APP_ENV is not set.
func Environment() string {
	env := strings.TrimSpace(os.Getenv("APP_ENV"))

	if env == "" {
		return gorest.Development
	}

	return strings.ToLower(env)
}

func envString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(value)
	}

	return fallback
}

func envSecret(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func envBool(key string, fallback bool) (bool, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback, nil
	}

	parsed, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean: %w", key, err)
	}

	return parsed, nil
}

func requireProductionEnv(key string) error {
	if strings.TrimSpace(os.Getenv(key)) == "" {
		return fmt.Errorf("%s must be set when APP_ENV=production", key)
	}

	return nil
}

func validateEnvironment() error {
	switch Environment() {
	case gorest.Development, gorest.Production, gorest.Testing:
		return nil
	default:
		return fmt.Errorf("APP_ENV must be %s, %s, or %s", gorest.Development, gorest.Production, gorest.Testing)
	}
}
