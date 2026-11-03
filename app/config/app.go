package config

import (
	"fmt"
	"os"

	"github.com/elsyahtech/gorest"
)

func AppConfig() (gorest.Config, error) {
	cfg := gorest.ConfigDefault // Load the default application configuration

	cfg.AppName = str(os.Getenv("GOREST_APP_NAME"), gorest.DefaultAppName)
	cfg.Version = str(os.Getenv("GOREST_APP_VERSION"), "v1.0.0")
	cfg.Environment = str(os.Getenv("GOREST_APP_ENV"), gorest.Development)

	if err := validateEnvironment(cfg.Environment); err != nil {
		return gorest.Config{}, fmt.Errorf("%w", err)
	}

	cfg.Timezone = str(os.Getenv("GOREST_APP_TIMEZONE"), gorest.DefaultTimezone)
	cfg.Debug = strToBool(os.Getenv("GOREST_APP_DEBUG"), false)

	if value, ok := os.LookupEnv("GOREST_APP_CHECK_FOR_UPDATES"); ok && value != "" {
		checkForUpdates := strToBool(value, true)
		cfg.CheckForUpdates = &checkForUpdates
	}

	return cfg, nil
}
