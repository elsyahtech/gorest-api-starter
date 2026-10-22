package config

import (
	"github.com/elsyahtech/gorest/log"
)

// Logger is the configuration for log file location and rotation (size, backup count, and age).
// GOREST runs the logger based on the values you set here,
// so adjust them to match your environment.
func Logger() log.Config {
	return log.Config{
		// LogDir defines the root directory path where log files will be stored.
		// If the directory does not exist, it should be created automatically on startup.
		// Example: "./logs" or "/var/log/myapp"
		// Default: "./logs"
		LogDir: "./logs",

		// InfoFile is the filename for general/info level logs.
		// Example: "app.log"
		// Default: "app.log"
		LogFile: "app.log",

		// Caller enables or disables the inclusion of caller information (such as
		// the file name and line number) in the log entries.
		// Example: true
		// Default: true
		Caller: true,

		// MaxFileSize defines the maximum size in megabytes (MB) that a log file can reach
		// before it gets rotated (split into a new file).
		// Example: 100 (means 100 MB)
		// Default: 100
		MaxFileSize: 100,

		// MaxBackups defines the maximum number of old log files to retain after rotation.
		// 0 means retain all old log files (not recommended for disk space management).
		// Example: 5 (keep up to 5 old backup files)
		// Default: 5
		MaxBackups: 5,

		// MaxAge defines the maximum number of days to retain old log files based on
		// the timestamp encoded in their filename.
		// 0 means forever (or controlled solely by MaxBackups).
		// Example: 7 (delete logs older than 7 days)
		// Default: 7
		MaxAge: 7,
	}
}
