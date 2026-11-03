package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elsyahtech/gorest"
)

func str(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	return value
}

func strPtr(value, defaultValue string) *string {
	if value == "" {
		return &defaultValue
	}

	return &value
}

func strToInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

func strToTimeDuration(value string, defaultValue time.Duration) time.Duration {
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

func strToBool(value string, defaultValue bool) bool {
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

func validateEnvironment(value string) error {
	switch value {
	case gorest.Development, gorest.Production, gorest.Testing:
		return nil
	default:
		return fmt.Errorf("GOREST_APP_ENV must be %s, %s, or %s", gorest.Development, gorest.Production, gorest.Testing)
	}
}
