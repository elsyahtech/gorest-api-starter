package config

import (
	"fmt"
	"strings"

	"github.com/elsyahtech/gorest"
	goresthelper "github.com/elsyahtech/gorest/helper"
)

func runtimeEnvironment() string {
	value := strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault("APP_ENV", gorest.Development))
	if value == "" {
		return gorest.Development
	}

	return strings.ToLower(value)
}

func validateEnvironment() error {
	environment := runtimeEnvironment()
	if err := goresthelper.ValidateApplicationEnvironmentMode(environment); err != nil {
		return err
	}

	switch environment {
	case gorest.Development, gorest.Production, gorest.Testing:
		return nil
	default:
		return fmt.Errorf("APP_ENV must be %s, %s, or %s", gorest.Development, gorest.Production, gorest.Testing)
	}
}

func requireProductionEnv(key string) error {
	if goresthelper.ReadTrimmedEnvironmentVariable(key) == "" {
		return fmt.Errorf("%s must be set when APP_ENV=production", key)
	}

	return nil
}
