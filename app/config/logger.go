package config

import (
	"fmt"

	"github.com/elsyahtech/gorest/log"
)

// LoggerConfig loads log file and rotation settings from environment variables.
func LoggerConfig() (*log.Config, error) {
	cfg := log.ConfigDefault

	var err error

	cfg.LogDir = envStr("GOREST_LOG_DIR", log.DefaultLogDir)
	cfg.LogFile = envStr("GOREST_LOG_FILE", log.DefaultLogFile)

	cfg.Caller, err = envStrToBool("GOREST_LOG_CALLER", log.DefaultCaller)
	if err != nil {
		return &log.Config{}, fmt.Errorf("failed to parse GOREST_LOG_CALLER: %w", err)
	}

	cfg.MaxFileSize = envStrToInt("GOREST_LOG_MAX_FILE_SIZE", log.DefaultMaxFileSize)
	cfg.MaxBackups = envStrToInt("GOREST_LOG_MAX_BACKUPS", log.DefaultMaxBackups)
	cfg.MaxAge = envStrToInt("GOREST_LOG_MAX_AGE", log.DefaultLogMaxAge)

	return &cfg, nil
}
