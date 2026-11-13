package config

import (
	goresthelper "github.com/elsyahtech/gorest/helper"
	"github.com/elsyahtech/gorest/server"
)

func listenerConfig() (*server.ListenerConfig, error) {
	listener := server.ListenerConfig{
		Network:           goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_LISTENER_NETWORK", ""),
		TLSCertFile:       goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_CERT_FILE", ""),
		TLSCertKeyFile:    goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_CERT_KEY_FILE", ""),
		TLSCertClientFile: goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_CERT_CLIENT_FILE", ""),
		UnixSocketPath:    goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_UNIX_SOCKET_PATH", ""),
	}

	var err error

	if listener.PreforkRecoverInterval, err = goresthelper.ReadEnvironmentDurationWithDefault(
		"GOREST_SERVER_PREFORK_RECOVER_INTERVAL",
		0,
	); err != nil {
		return &server.ListenerConfig{}, err
	}

	if listener.ShutdownTimeout, err = goresthelper.ReadEnvironmentDurationWithDefault(
		"GOREST_SERVER_SHUTDOWN_TIMEOUT",
		0,
	); err != nil {
		return &server.ListenerConfig{}, err
	}

	if listener.PreforkShutdownGracePeriod, err = goresthelper.ReadEnvironmentDurationWithDefault(
		"GOREST_SERVER_PREFORK_SHUTDOWN_GRACE_PERIOD",
		0,
	); err != nil {
		return &server.ListenerConfig{}, err
	}

	if listener.PreforkRecoverThreshold, err = goresthelper.ReadEnvironmentNonNegativeInteger("GOREST_SERVER_PREFORK_RECOVER_THRESHOLD"); err != nil {
		return &server.ListenerConfig{}, err
	}

	if listener.UnixSocketFileMode, err = goresthelper.ReadEnvironmentFileMode("GOREST_SERVER_UNIX_SOCKET_FILE_MODE"); err != nil {
		return &server.ListenerConfig{}, err
	}

	if listener.TLSMinVersion, err = goresthelper.ReadEnvironmentTLSMinimumVersion("GOREST_SERVER_TLS_MIN_VERSION"); err != nil {
		return &server.ListenerConfig{}, err
	}

	if listener.EnablePrefork, err = goresthelper.ReadEnvironmentBooleanWithDefault("GOREST_SERVER_ENABLE_PREFORK", false); err != nil {
		return &server.ListenerConfig{}, err
	}

	if listener.EnablePrintRoutes, err = goresthelper.ReadEnvironmentBooleanWithDefault("GOREST_SERVER_ENABLE_PRINT_ROUTES", false); err != nil {
		return &server.ListenerConfig{}, err
	}

	return &listener, nil
}
