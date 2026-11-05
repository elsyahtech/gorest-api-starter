package config

import (
	"fmt"

	"github.com/elsyahtech/gorest/log"
)

// LoggerConfig loads log file and rotation settings from environment variables.
func LoggerConfig() (log.Config, error) {
	cfg := log.ConfigDefault

	var err error

	cfg.LogDir = str("GOREST_LOG_DIR", log.DefaultLogDir)
	cfg.LogFile = str("GOREST_LOG_FILE", log.DefaultLogFile)

	cfg.Caller, err = strToBool("GOREST_LOG_CALLER", log.DefaultCaller)
	if err != nil {
		return log.Config{}, fmt.Errorf("failed to parse GOREST_LOG_CALLER: %w", err)
	}

	cfg.MaxFileSize = strToInt("GOREST_LOG_MAX_FILE_SIZE", log.DefaultMaxFileSize)
	cfg.MaxBackups = strToInt("GOREST_LOG_MAX_BACKUPS", log.DefaultMaxBackups)
	cfg.MaxAge = strToInt("GOREST_LOG_MAX_AGE", log.DefaultLogMaxAge)

	return cfg, nil
}
