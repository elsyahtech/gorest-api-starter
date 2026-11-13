package config

import (
	"fmt"

	"github.com/elsyahtech/gorest"
	goresthelper "github.com/elsyahtech/gorest/helper"
)

// AppConfig loads application settings from environment variables and validates the environment.
func AppConfig() (*gorest.Config, error) {
	cfg := gorest.ConfigDefault

	var err error

	cfg.AppName = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_APP_NAME",
		gorest.DefaultAppName,
	)

	cfg.Version = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_APP_VERSION",
		"v1.0.0",
	)

	cfg.Environment = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_APP_ENV",
		gorest.Development,
	)

	if err := goresthelper.ValidateApplicationEnvironmentMode(cfg.Environment); err != nil {
		return &gorest.Config{}, fmt.Errorf("%w", err)
	}

	switch cfg.Environment {
	case gorest.Development, gorest.Production, gorest.Testing:
	default:
		return &gorest.Config{}, fmt.Errorf("GOREST_APP_ENV must be %s, %s, or %s", gorest.Development, gorest.Production, gorest.Testing)
	}

	cfg.Timezone = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_APP_TIMEZONE",
		gorest.DefaultTimezone,
	)

	cfg.Debug, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_APP_DEBUG",
		false,
	)
	if err != nil {
		return &gorest.Config{}, fmt.Errorf("failed to parse GOREST_APP_DEBUG: %w", err)
	}

	if checkForUpdates, err := goresthelper.ReadEnvironmentBooleanPointer("GOREST_APP_CHECK_FOR_UPDATES"); err != nil {
		return &gorest.Config{}, fmt.Errorf("failed to parse GOREST_APP_CHECK_FOR_UPDATES: %w", err)
	} else if checkForUpdates != nil {
		cfg.CheckForUpdates = checkForUpdates
	}

	return &cfg, nil
}
