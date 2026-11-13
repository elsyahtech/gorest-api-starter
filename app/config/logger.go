package config

import (
	"fmt"

	goresthelper "github.com/elsyahtech/gorest/helper"
	"github.com/elsyahtech/gorest/log"
)

// LoggerConfig loads log file and rotation settings from environment variables.
func LoggerConfig() (*log.Config, error) {
	cfg := log.ConfigDefault

	var err error

	cfg.LogDir = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_LOG_DIR",
		log.DefaultLogDir,
	)

	cfg.LogFile = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_LOG_FILE",
		log.DefaultLogFile,
	)

	cfg.Caller, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_LOG_CALLER",
		log.DefaultCaller,
	)
	if err != nil {
		return &log.Config{}, fmt.Errorf("failed to parse GOREST_LOG_CALLER: %w", err)
	}

	cfg.MaxFileSize, err = goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_LOG_MAX_FILE_SIZE",
		log.DefaultMaxFileSize,
	)
	if err != nil {
		return &log.Config{}, fmt.Errorf("failed to parse GOREST_LOG_MAX_FILE_SIZE: %w", err)
	}

	cfg.MaxBackups, err = goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_LOG_MAX_BACKUPS",
		log.DefaultMaxBackups,
	)
	if err != nil {
		return &log.Config{}, fmt.Errorf("failed to parse GOREST_LOG_MAX_BACKUPS: %w", err)
	}

	cfg.MaxAge, err = goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_LOG_MAX_AGE",
		log.DefaultLogMaxAge,
	)
	if err != nil {
		return &log.Config{}, fmt.Errorf("failed to parse GOREST_LOG_MAX_AGE: %w", err)
	}

	return &cfg, nil
}
