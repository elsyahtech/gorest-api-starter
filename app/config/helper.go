package config

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elsyahtech/gorest"
)

func envString(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func envDuration(key string) (time.Duration, error) {
	value := envString(key)
	if value == "" {
		return 0, nil
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid Go duration: %w", key, err)
	}

	if duration < 0 {
		return 0, fmt.Errorf("%s cannot be negative", key)
	}

	return duration, nil
}

func envNonNegativeInt(key string) (int, error) {
	value := envString(key)
	if value == "" {
		return 0, nil
	}

	number, err := strconv.Atoi(value)
	if err != nil || number < 0 {
		return 0, fmt.Errorf("%s must be a non-negative integer", key)
	}

	return number, nil
}

func envFileMode(key string) (os.FileMode, error) {
	value := envString(key)
	if value == "" {
		return 0, nil
	}

	mode, err := strconv.ParseUint(value, 8, 32)
	if err != nil || mode > 0o777 {
		return 0, fmt.Errorf("%s must be an octal permission mode between 0000 and 0777", key)
	}

	return os.FileMode(mode), nil
}

func envTLSMinVersion(key string) (uint16, error) {
	switch strings.ToLower(envString(key)) {
	case "":
		return 0, nil
	case "1.2", "tls1.2", "tlsv1.2":
		return tls.VersionTLS12, nil
	case "1.3", "tls1.3", "tlsv1.3":
		return tls.VersionTLS13, nil
	default:
		return 0, fmt.Errorf("%s must be 1.2 or 1.3", key)
	}
}

func str(key, defaultValue string) string {
	value := envString(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func strPtr(key, defaultValue string) *string {
	value := str(key, defaultValue)

	return &value
}

func strToInt(key string, defaultValue int) int {
	value := envString(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

func strToTimeDuration(key string, defaultValue time.Duration) time.Duration {
	value := envString(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

func toStrArray(key string, defaultValue []string) []string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	value = strings.TrimSpace(value)

	if value == "" {
		return make([]string, 0)
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

func strToBool(key string, defaultValue bool) (bool, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean: %w", key, err)
	}

	return parsed, nil
}

func validateEnvironment(value string) error {
	switch value {
	case gorest.Development, gorest.Production, gorest.Testing:
		return nil
	default:
		return fmt.Errorf("GOREST_APP_ENV must be %s, %s, or %s", gorest.Development, gorest.Production, gorest.Testing)
	}
}
