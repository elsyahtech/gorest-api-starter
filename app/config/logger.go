package config

import (
	"os"

	"github.com/elsyahtech/gorest/log"
)

// Logger is the configuration for log file location and rotation (size, backup count, and age).
// GOREST runs the logger based on the values you set here,
// so adjust them to match your environment.
func LoggerConfig() (log.Config, error) {
	// Load the default logger configuration
	cfg := log.ConfigDefault

	// LogDir defines the root directory path where log files will be stored.
	// If the directory does not exist, it should be created automatically on startup.
	// Example: "./logs" or "/var/log/myapp"
	// Default: "./logs"
	cfg.LogDir = str(os.Getenv("GOREST_LOG_DIR"), log.DefaultLogDir)

	// InfoFile is the filename for general/info level logs.
	// Example: "app.log"
	// Default: "app.log"
	cfg.LogFile = str(os.Getenv("GOREST_LOG_FILE"), log.DefaultLogFile)

	cfg.Caller = strToBool(os.Getenv("GOREST_LOG_CALLER"), log.DefaultCaller)

	// MaxFileSize defines the maximum size in megabytes (MB) that a log file can reach
	// before it gets rotated (split into a new file).
	// Example: 100 (means 100 MB)
	// Default: 100
	cfg.MaxFileSize = strToInt(os.Getenv("GOREST_LOG_MAX_FILE_SIZE"), log.DefaultMaxFileSize)

	// MaxBackups defines the maximum number of old log files to retain after rotation.
	// 0 means retain all old log files (not recommended for disk space management).
	// Example: 5 (keep up to 5 old backup files)
	// Default: 5
	cfg.MaxBackups = strToInt(os.Getenv("GOREST_LOG_MAX_BACKUPS"), log.DefaultMaxBackups)

	// MaxAge defines the maximum number of days to retain old log files based on
	// the timestamp encoded in their filename.
	// 0 means forever (or controlled solely by MaxBackups).
	// Example: 7 (delete logs older than 7 days)
	// Default: 7
	cfg.MaxAge = strToInt(os.Getenv("GOREST_LOG_MAX_AGE"), log.DefaultLogMaxAge)

	return cfg, nil
}
