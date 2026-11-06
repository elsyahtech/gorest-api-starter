package config

import (
	"fmt"
	"os"

	"github.com/elsyahtech/gorest"
)

// AppConfig loads application settings from environment variables and validates the environment.
func AppConfig() (*gorest.Config, error) {
	cfg := gorest.ConfigDefault

	var err error

	cfg.AppName = envStr("GOREST_APP_NAME", gorest.DefaultAppName)
	cfg.Version = envStr("GOREST_APP_VERSION", "v1.0.0")
	cfg.Environment = envStr("GOREST_APP_ENV", gorest.Development)

	if err := validateEnvironment(cfg.Environment); err != nil {
		return &gorest.Config{}, fmt.Errorf("%w", err)
	}

	cfg.Timezone = envStr("GOREST_APP_TIMEZONE", gorest.DefaultTimezone)

	cfg.Debug, err = envStrToBool("GOREST_APP_DEBUG", false)
	if err != nil {
		return &gorest.Config{}, fmt.Errorf("failed to parse GOREST_APP_DEBUG: %w", err)
	}

	if value, ok := os.LookupEnv("GOREST_APP_CHECK_FOR_UPDATES"); ok && value != "" {
		checkForUpdates, err := envStrToBool("GOREST_APP_CHECK_FOR_UPDATES", true)
		if err != nil {
			return &gorest.Config{}, fmt.Errorf("failed to parse GOREST_APP_CHECK_FOR_UPDATES: %w", err)
		}

		cfg.CheckForUpdates = &checkForUpdates
	}

	return &cfg, nil
}
