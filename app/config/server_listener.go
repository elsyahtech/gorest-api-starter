package config

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elsyahtech/gorest/server"
)

func loadServerListenerConfig() (server.ListenerConfig, error) {
	listener := server.ListenerConfig{
		Network:           listenerString("GOREST_SERVER_LISTENER_NETWORK"),
		TLSCertFile:       listenerString("GOREST_SERVER_CERT_FILE"),
		TLSCertKeyFile:    listenerString("GOREST_SERVER_CERT_KEY_FILE"),
		TLSCertClientFile: listenerString("GOREST_SERVER_CERT_CLIENT_FILE"),
		UnixSocketPath:    listenerString("GOREST_SERVER_UNIX_SOCKET_PATH"),
	}

	var err error
	if listener.PreforkRecoverInterval, err = listenerDuration("GOREST_SERVER_PREFORK_RECOVER_INTERVAL"); err != nil {
		return server.ListenerConfig{}, err
	}
	if listener.ShutdownTimeout, err = listenerDuration("GOREST_SERVER_SHUTDOWN_TIMEOUT"); err != nil {
		return server.ListenerConfig{}, err
	}
	if listener.PreforkShutdownGracePeriod, err = listenerDuration("GOREST_SERVER_PREFORK_SHUTDOWN_GRACE_PERIOD"); err != nil {
		return server.ListenerConfig{}, err
	}
	if listener.PreforkRecoverThreshold, err = listenerNonNegativeInt("GOREST_SERVER_PREFORK_RECOVER_THRESHOLD"); err != nil {
		return server.ListenerConfig{}, err
	}
	if listener.UnixSocketFileMode, err = listenerFileMode("GOREST_SERVER_UNIX_SOCKET_FILE_MODE"); err != nil {
		return server.ListenerConfig{}, err
	}
	if listener.TLSMinVersion, err = listenerTLSMinVersion("GOREST_SERVER_TLS_MIN_VERSION"); err != nil {
		return server.ListenerConfig{}, err
	}
	if listener.EnablePrefork, err = listenerBool("GOREST_SERVER_ENABLE_PREFORK"); err != nil {
		return server.ListenerConfig{}, err
	}
	if listener.EnablePrintRoutes, err = listenerBool("GOREST_SERVER_ENABLE_PRINT_ROUTES"); err != nil {
		return server.ListenerConfig{}, err
	}

	return listener, nil
}

func listenerString(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func listenerDuration(key string) (time.Duration, error) {
	value := listenerString(key)
	if value == "" {
		return 0, nil
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid Go duration: %w", key, err)
	}
	if duration < 0 {
		return 0, fmt.Errorf("%s cannot be negative", key)
	}

	return duration, nil
}

func listenerNonNegativeInt(key string) (int, error) {
	value := listenerString(key)
	if value == "" {
		return 0, nil
	}

	number, err := strconv.Atoi(value)
	if err != nil || number < 0 {
		return 0, fmt.Errorf("%s must be a non-negative integer", key)
	}

	return number, nil
}

func listenerFileMode(key string) (os.FileMode, error) {
	value := listenerString(key)
	if value == "" {
		return 0, nil
	}

	mode, err := strconv.ParseUint(value, 8, 32)
	if err != nil || mode > 0o777 {
		return 0, fmt.Errorf("%s must be an octal permission mode between 0000 and 0777", key)
	}

	return os.FileMode(mode), nil
}

func listenerTLSMinVersion(key string) (uint16, error) {
	switch strings.ToLower(listenerString(key)) {
	case "":
		return 0, nil
	case "1.2", "tls1.2", "tlsv1.2":
		return tls.VersionTLS12, nil
	case "1.3", "tls1.3", "tlsv1.3":
		return tls.VersionTLS13, nil
	default:
		return 0, fmt.Errorf("%s must be 1.2 or 1.3", key)
	}
}

func listenerBool(key string) (bool, error) {
	value := listenerString(key)
	if value == "" {
		return false, nil
	}

	enabled, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean: %w", key, err)
	}

	return enabled, nil
}
