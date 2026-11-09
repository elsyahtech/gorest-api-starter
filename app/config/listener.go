package config

import (
	"github.com/elsyahtech/gorest/server"
)

func listenerConfig() (*server.ListenerConfig, error) {
	listener := server.ListenerConfig{
		Network:           envStr("GOREST_SERVER_LISTENER_NETWORK", ""),
		TLSCertFile:       envStr("GOREST_SERVER_CERT_FILE", ""),
		TLSCertKeyFile:    envStr("GOREST_SERVER_CERT_KEY_FILE", ""),
		TLSCertClientFile: envStr("GOREST_SERVER_CERT_CLIENT_FILE", ""),
		UnixSocketPath:    envStr("GOREST_SERVER_UNIX_SOCKET_PATH", ""),
	}

	var err error

	if listener.PreforkRecoverInterval, err = envStrToTimeDuration("GOREST_SERVER_PREFORK_RECOVER_INTERVAL", 0); err != nil {
		return &server.ListenerConfig{}, err
	}

	if listener.ShutdownTimeout, err = envStrToTimeDuration("GOREST_SERVER_SHUTDOWN_TIMEOUT", 0); err != nil {
		return &server.ListenerConfig{}, err
	}

	if listener.PreforkShutdownGracePeriod, err = envStrToTimeDuration("GOREST_SERVER_PREFORK_SHUTDOWN_GRACE_PERIOD", 0); err != nil {
		return &server.ListenerConfig{}, err
	}

	if listener.PreforkRecoverThreshold, err = envNonNegativeInt("GOREST_SERVER_PREFORK_RECOVER_THRESHOLD"); err != nil {
		return &server.ListenerConfig{}, err
	}

	if listener.UnixSocketFileMode, err = envFileMode("GOREST_SERVER_UNIX_SOCKET_FILE_MODE"); err != nil {
		return &server.ListenerConfig{}, err
	}

	if listener.TLSMinVersion, err = envTLSMinVersion("GOREST_SERVER_TLS_MIN_VERSION"); err != nil {
		return &server.ListenerConfig{}, err
	}

	if listener.EnablePrefork, err = envStrToBool("GOREST_SERVER_ENABLE_PREFORK", false); err != nil {
		return &server.ListenerConfig{}, err
	}

	if listener.EnablePrintRoutes, err = envStrToBool("GOREST_SERVER_ENABLE_PRINT_ROUTES", false); err != nil {
		return &server.ListenerConfig{}, err
	}

	return &listener, nil
}
