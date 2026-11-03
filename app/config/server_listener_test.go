package config

import (
	"crypto/tls"
	"os"
	"testing"
	"time"

	"github.com/elsyahtech/gorest/server"
)

func TestServerListenerUsesGORESTDefaultsWhenEnvironmentIsEmpty(t *testing.T) {
	isolateConfigurationEnvironment(t)

	cfg, err := ServerConfig()
	if err != nil {
		t.Fatalf("ServerConfig() error = %v", err)
	}

	listener := server.Run(cfg).Listener
	if listener.Network != "tcp4" ||
		listener.TLSMinVersion != tls.VersionTLS12 ||
		listener.UnixSocketFileMode != 0o770 ||
		listener.ShutdownTimeout != server.DefaultServerShutdownTimeout ||
		listener.PreforkShutdownGracePeriod != 5*time.Second ||
		listener.PreforkRecoverThreshold < 1 {
		t.Fatalf("listener defaults = %+v, want GOREST defaults", listener)
	}
	if listener.TLSCertFile != "" || listener.TLSCertKeyFile != "" ||
		listener.TLSCertClientFile != "" || listener.UnixSocketPath != "" ||
		listener.PreforkRecoverInterval != 0 || listener.EnablePrefork || listener.EnablePrintRoutes {
		t.Errorf("optional listener values = %+v, want empty or disabled", listener)
	}
	if _, _, err := server.New(&cfg); err != nil {
		t.Fatalf("server.New() error with empty optional listener values = %v", err)
	}
}

func TestServerListenerMapsEnvironmentOverrides(t *testing.T) {
	isolateConfigurationEnvironment(t)
	for key, value := range map[string]string{
		"GOREST_SERVER_LISTENER_NETWORK":              "tcp6",
		"GOREST_SERVER_CERT_FILE":                     "/etc/example/server.crt",
		"GOREST_SERVER_CERT_KEY_FILE":                 "/etc/example/server.key",
		"GOREST_SERVER_CERT_CLIENT_FILE":              "/etc/example/client-ca.crt",
		"GOREST_SERVER_PREFORK_RECOVER_INTERVAL":      "2s",
		"GOREST_SERVER_SHUTDOWN_TIMEOUT":              "45s",
		"GOREST_SERVER_PREFORK_SHUTDOWN_GRACE_PERIOD": "8s",
		"GOREST_SERVER_PREFORK_RECOVER_THRESHOLD":     "4",
		"GOREST_SERVER_UNIX_SOCKET_FILE_MODE":         "0700",
		"GOREST_SERVER_TLS_MIN_VERSION":               "1.3",
		"GOREST_SERVER_ENABLE_PREFORK":                "true",
		"GOREST_SERVER_ENABLE_PRINT_ROUTES":           "true",
	} {
		t.Setenv(key, value)
	}

	cfg, err := ServerConfig()
	if err != nil {
		t.Fatalf("ServerConfig() error = %v", err)
	}

	listener := cfg.Listener
	if listener.Network != "tcp6" ||
		listener.TLSCertFile != "/etc/example/server.crt" ||
		listener.TLSCertKeyFile != "/etc/example/server.key" ||
		listener.TLSCertClientFile != "/etc/example/client-ca.crt" ||
		listener.PreforkRecoverInterval != 2*time.Second ||
		listener.ShutdownTimeout != 45*time.Second ||
		listener.PreforkShutdownGracePeriod != 8*time.Second ||
		listener.PreforkRecoverThreshold != 4 ||
		listener.UnixSocketFileMode != 0o700 ||
		listener.TLSMinVersion != tls.VersionTLS13 ||
		!listener.EnablePrefork || !listener.EnablePrintRoutes {
		t.Fatalf("listener config = %+v, want configured environment values", listener)
	}
}

func TestServerListenerMapsUnixSocket(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_SERVER_LISTENER_NETWORK", "unix")
	t.Setenv("GOREST_SERVER_UNIX_SOCKET_PATH", "/run/example/server.sock")

	cfg, err := ServerConfig()
	if err != nil {
		t.Fatalf("ServerConfig() error = %v", err)
	}
	if cfg.Listener.Network != "unix" || cfg.Listener.UnixSocketPath != "/run/example/server.sock" {
		t.Fatalf("listener config = %+v, want Unix socket settings", cfg.Listener)
	}
}

func TestServerListenerRejectsInvalidEnvironment(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{name: "invalid shutdown duration", key: "GOREST_SERVER_SHUTDOWN_TIMEOUT", value: "soon"},
		{name: "negative shutdown duration", key: "GOREST_SERVER_SHUTDOWN_TIMEOUT", value: "-1s"},
		{name: "negative prefork recovery threshold", key: "GOREST_SERVER_PREFORK_RECOVER_THRESHOLD", value: "-1"},
		{name: "invalid Unix socket mode", key: "GOREST_SERVER_UNIX_SOCKET_FILE_MODE", value: "0888"},
		{name: "unsupported TLS version", key: "GOREST_SERVER_TLS_MIN_VERSION", value: "1.1"},
		{name: "invalid prefork boolean", key: "GOREST_SERVER_ENABLE_PREFORK", value: "yes"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isolateConfigurationEnvironment(t)
			t.Setenv(test.key, test.value)
			if _, err := ServerConfig(); err == nil {
				t.Fatalf("ServerConfig() succeeded with %s=%q", test.key, test.value)
			}
		})
	}
}

func TestServerListenerInvalidTLSCombinationIsRejectedByGOREST(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_SERVER_CERT_FILE", "/etc/example/server.crt")

	cfg, err := ServerConfig()
	if err != nil {
		t.Fatalf("ServerConfig() error = %v", err)
	}
	if _, _, err := server.New(&cfg); err == nil {
		t.Fatal("server.New() succeeded with an incomplete TLS certificate pair")
	}
}

func TestServerListenerFileModeUsesOctal(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_SERVER_UNIX_SOCKET_FILE_MODE", "0700")

	cfg, err := ServerConfig()
	if err != nil {
		t.Fatalf("ServerConfig() error = %v", err)
	}
	if cfg.Listener.UnixSocketFileMode != os.FileMode(0o700) {
		t.Fatalf("UnixSocketFileMode = %04o, want 0700", cfg.Listener.UnixSocketFileMode)
	}
}
