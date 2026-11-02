package config

import (
	"os"
	"testing"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest/filter"
)

var configurationEnvironmentKeys = []string{
	"APP_ENV",
	"GOREST_SERVER_HOST",
	"GOREST_SERVER_PORT",
	"GOREST_CORS_ORIGINS",
	"GOREST_DB_DRIVER",
	"GOREST_DB_HOST",
	"GOREST_DB_PORT",
	"GOREST_DB_NAME",
	"GOREST_DB_USERNAME",
	"GOREST_DB_PASSWORD",
	"GOREST_DB_TLS",
	"GOREST_DB_ALLOW_INSECURE_PRODUCTION",
	"GOREST_DB_LIB_DIR",
	"GOREST_DB_SSL_CA_PATH",
	"GOREST_DB_SSL_CERT_PATH",
	"GOREST_DB_SSL_KEY_PATH",
	"GOREST_DB_DIRECTORY",
	"GOREST_DB_MIGRATION",
	"GOREST_DB_SEEDER",
	"GOREST_REDIS_DISABLE",
	"GOREST_REDIS_TLS",
	"GOREST_REDIS_DATABASE",
	"GOREST_REDIS_PORT",
	"GOREST_REDIS_HOST",
	"GOREST_REDIS_USERNAME",
	"GOREST_REDIS_PASSWORD",
	"GOREST_REDIS_SSL_CA_PATH",
	"GOREST_REDIS_SSL_CERT_PATH",
	"GOREST_REDIS_SSL_KEY_PATH",
	"GOREST_FILTER_SIGNATURE_KEY",
	"GOREST_HTTP_METRICS_ALLOWED_IPS",
	"GOREST_DISABLE_MONITORING_ENDPOINTS",
	"GOREST_ENABLE_PER_ROUTE_HTTP_METRICS",
}

func TestLoadFiltersUsesDefaultConfiguration(t *testing.T) {
	isolateConfigurationEnvironment(t)

	cfg, err := LoadFilters()
	if err != nil {
		t.Fatalf("LoadFilters() error = %v", err)
	}

	if cfg.SignatureKey == nil || *cfg.SignatureKey != "your-signature-key" {
		got := "<nil>"
		if cfg.SignatureKey != nil {
			got = *cfg.SignatureKey
		}
		t.Fatalf("LoadFilters().SignatureKey = %q, want the default signing key", got)
	}
	if cfg.SigningMethod != "HS256" {
		t.Errorf("LoadFilters().SigningMethod = %q, want HS256", cfg.SigningMethod)
	}
	if cfg.AccessTokenDuration != 15*time.Minute {
		t.Errorf("LoadFilters().AccessTokenDuration = %s, want 15m", cfg.AccessTokenDuration)
	}
	if cfg.RefreshTokenDuration != 7*24*time.Hour {
		t.Errorf("LoadFilters().RefreshTokenDuration = %s, want 168h", cfg.RefreshTokenDuration)
	}
	if cfg.Issuer != "your-issuer" {
		t.Errorf("LoadFilters().Issuer = %q, want your-issuer", cfg.Issuer)
	}
	if cfg.CookieSameSite != "Lax" || cfg.CookiePath != "/" || cfg.CookieDomain != "" {
		t.Errorf("LoadFilters() cookie scope = (%q, %q, %q), want (Lax, /, empty)",
			cfg.CookieSameSite, cfg.CookiePath, cfg.CookieDomain)
	}
	if !cfg.CookieSecure || !cfg.CookieHTTPOnly || cfg.CookieSessionOnly {
		t.Errorf("LoadFilters() cookie flags = (secure:%t, httpOnly:%t, sessionOnly:%t), want (true, true, false)",
			cfg.CookieSecure, cfg.CookieHTTPOnly, cfg.CookieSessionOnly)
	}
	if cfg.LoggerConfig.DisableLogging ||
		len(cfg.LoggerConfig.EnabledEvents) != 1 ||
		cfg.LoggerConfig.EnabledEvents[0] != filter.AuthenticationEventMissingCredentials {
		t.Errorf("LoadFilters().LoggerConfig = %+v, want missing-credentials logging enabled", cfg.LoggerConfig)
	}
}

func TestLoadFiltersUsesSignatureKeyEnvironmentVariable(t *testing.T) {
	isolateConfigurationEnvironment(t)
	const signatureKey = "unit-test-signature-key"
	t.Setenv("GOREST_FILTER_SIGNATURE_KEY", signatureKey)

	cfg, err := LoadFilters()
	if err != nil {
		t.Fatalf("LoadFilters() error = %v", err)
	}
	if cfg.SignatureKey == nil || *cfg.SignatureKey != signatureKey {
		got := "<nil>"
		if cfg.SignatureKey != nil {
			got = *cfg.SignatureKey
		}
		t.Fatalf("LoadFilters().SignatureKey = %q, want configured key", got)
	}
}

func TestLoadRedisDefaultsDisabled(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("APP_ENV", gorest.Development)

	cfg, err := LoadRedis()
	if err != nil {
		t.Fatalf("LoadRedis() error = %v", err)
	}
	if !cfg.Disable {
		t.Fatal("LoadRedis().Disable = false, want Redis disabled by default")
	}
}

func TestLoadRedisCanBeExplicitlyEnabled(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("APP_ENV", gorest.Development)
	t.Setenv("GOREST_REDIS_DISABLE", "false")

	cfg, err := LoadRedis()
	if err != nil {
		t.Fatalf("LoadRedis() error = %v", err)
	}
	if cfg.Disable {
		t.Fatal("LoadRedis().Disable = true, want Redis enabled when explicitly configured")
	}
}

func TestValidateEnvironmentAcceptsOnlyGORESTModes(t *testing.T) {
	for _, mode := range []string{gorest.Development, gorest.Production, gorest.Testing} {
		t.Run("accepts "+mode, func(t *testing.T) {
			isolateConfigurationEnvironment(t)
			t.Setenv("APP_ENV", mode)
			if err := validateEnvironment(); err != nil {
				t.Fatalf("validateEnvironment() error = %v", err)
			}
		})
	}

	for _, mode := range []string{"test", "staging"} {
		t.Run("rejects "+mode, func(t *testing.T) {
			isolateConfigurationEnvironment(t)
			t.Setenv("APP_ENV", mode)
			if err := validateEnvironment(); err == nil {
				t.Fatalf("validateEnvironment() accepted APP_ENV=%q", mode)
			}
		})
	}

	t.Run("defaults an empty value to development", func(t *testing.T) {
		isolateConfigurationEnvironment(t)
		t.Setenv("APP_ENV", "")
		if got := Environment(); got != gorest.Development {
			t.Fatalf("Environment() = %q, want %q", got, gorest.Development)
		}
		if err := validateEnvironment(); err != nil {
			t.Fatalf("validateEnvironment() error = %v", err)
		}
	})
}

func isolateConfigurationEnvironment(t *testing.T) {
	t.Helper()

	type environmentValue struct {
		value   string
		present bool
	}

	previous := make(map[string]environmentValue, len(configurationEnvironmentKeys))
	for _, key := range configurationEnvironmentKeys {
		value, present := os.LookupEnv(key)
		previous[key] = environmentValue{value: value, present: present}
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("unset %s: %v", key, err)
		}
	}
	t.Cleanup(func() {
		for key, old := range previous {
			var err error
			if old.present {
				err = os.Setenv(key, old.value)
			} else {
				err = os.Unsetenv(key)
			}
			if err != nil {
				t.Errorf("restore %s: %v", key, err)
			}
		}
	})
}

func setValidProductionEnvironment(t *testing.T) {
	t.Helper()
	isolateConfigurationEnvironment(t)
	for key, value := range map[string]string{
		"APP_ENV":             "production",
		"GOREST_SERVER_HOST":  "0.0.0.0",
		"GOREST_SERVER_PORT":  "8080",
		"GOREST_CORS_ORIGINS": "https://app.example.com",
		"GOREST_DB_DRIVER":    "postgres",
		"GOREST_DB_HOST":      "db.example.internal",
		"GOREST_DB_PORT":      "5432",
		"GOREST_DB_NAME":      "app",
		"GOREST_DB_USERNAME":  "app_user",
		"GOREST_DB_PASSWORD":  "test-only-password",
		"GOREST_DB_TLS":       "require",
		"GOREST_DB_MIGRATION": "false",
		"GOREST_DB_SEEDER":    "false",
	} {
		t.Setenv(key, value)
	}
}

func TestLoadServerProductionConfiguration(t *testing.T) {
	t.Run("accepts explicit bind and CORS settings", func(t *testing.T) {
		setValidProductionEnvironment(t)
		cfg, err := LoadServer()
		if err != nil {
			t.Fatalf("LoadServer() error = %v", err)
		}
		if cfg.Host != "0.0.0.0" || cfg.Port != 8080 {
			t.Fatalf("server bind = %s:%d, want 0.0.0.0:8080", cfg.Host, cfg.Port)
		}
	})

	for _, missing := range []string{"GOREST_SERVER_HOST", "GOREST_SERVER_PORT", "GOREST_CORS_ORIGINS"} {
		t.Run("requires "+missing, func(t *testing.T) {
			setValidProductionEnvironment(t)
			if err := os.Unsetenv(missing); err != nil {
				t.Fatalf("unset %s: %v", missing, err)
			}
			if _, err := LoadServer(); err == nil {
				t.Fatalf("LoadServer() succeeded without %s", missing)
			}
		})
	}

	t.Run("rejects an empty CORS origin list", func(t *testing.T) {
		setValidProductionEnvironment(t)
		t.Setenv("GOREST_CORS_ORIGINS", " , ")
		if _, err := LoadServer(); err == nil {
			t.Fatal("LoadServer() succeeded with an empty CORS origin list")
		}
	})
}

func TestLoadDatabaseProductionConfiguration(t *testing.T) {
	t.Run("accepts explicit server database settings with startup tasks disabled", func(t *testing.T) {
		setValidProductionEnvironment(t)
		cfg, err := LoadDatabase()
		if err != nil {
			t.Fatalf("LoadDatabase() error = %v", err)
		}
		if cfg.Migration || cfg.Seeder {
			t.Fatalf("production startup tasks = migration:%t seeder:%t, want both disabled", cfg.Migration, cfg.Seeder)
		}
	})

	t.Run("allows insecure connections when explicitly enabled outside production", func(t *testing.T) {
		isolateConfigurationEnvironment(t)
		t.Setenv("APP_ENV", gorest.Testing)
		t.Setenv("GOREST_DB_ALLOW_INSECURE_PRODUCTION", "true")

		cfg, err := LoadDatabase()
		if err != nil {
			t.Fatalf("LoadDatabase() error = %v", err)
		}
		if !cfg.AllowInsecureProduction {
			t.Fatal("LoadDatabase().AllowInsecureProduction = false, want true")
		}
	})

	t.Run("rejects invalid insecure-connection setting", func(t *testing.T) {
		isolateConfigurationEnvironment(t)
		t.Setenv("GOREST_DB_ALLOW_INSECURE_PRODUCTION", "sometimes")

		if _, err := LoadDatabase(); err == nil {
			t.Fatal("LoadDatabase() succeeded with an invalid GOREST_DB_ALLOW_INSECURE_PRODUCTION value")
		}
	})

	t.Run("rejects insecure connections in production", func(t *testing.T) {
		setValidProductionEnvironment(t)
		t.Setenv("GOREST_DB_ALLOW_INSECURE_PRODUCTION", "true")

		if _, err := LoadDatabase(); err == nil {
			t.Fatal("LoadDatabase() succeeded with GOREST_DB_ALLOW_INSECURE_PRODUCTION=true in production")
		}
	})

	for _, missing := range []string{"GOREST_DB_DRIVER", "GOREST_DB_HOST", "GOREST_DB_NAME", "GOREST_DB_USERNAME", "GOREST_DB_PASSWORD", "GOREST_DB_TLS"} {
		t.Run("requires "+missing, func(t *testing.T) {
			setValidProductionEnvironment(t)
			if err := os.Unsetenv(missing); err != nil {
				t.Fatalf("unset %s: %v", missing, err)
			}
			if _, err := LoadDatabase(); err == nil {
				t.Fatalf("LoadDatabase() succeeded without %s", missing)
			}
		})
	}

	for _, startupTask := range []string{"GOREST_DB_MIGRATION", "GOREST_DB_SEEDER"} {
		t.Run("rejects enabled "+startupTask, func(t *testing.T) {
			setValidProductionEnvironment(t)
			t.Setenv(startupTask, "true")
			if _, err := LoadDatabase(); err == nil {
				t.Fatalf("LoadDatabase() succeeded with %s=true", startupTask)
			}
		})
	}

	t.Run("requires an explicit directory for SQLite", func(t *testing.T) {
		setValidProductionEnvironment(t)
		t.Setenv("GOREST_DB_DRIVER", "sqlite")
		if _, err := LoadDatabase(); err == nil {
			t.Fatal("LoadDatabase() succeeded for production SQLite without GOREST_DB_DIRECTORY")
		}
		directory := "/var/lib/gorest"
		t.Setenv("GOREST_DB_DIRECTORY", directory)
		cfg, err := LoadDatabase()
		if err != nil {
			t.Fatalf("LoadDatabase() with explicit SQLite directory error = %v", err)
		}
		if cfg.Driver != "sqlite" || cfg.DatabaseDirectory != directory {
			t.Errorf("production SQLite location = (%q, %q), want (sqlite, %q)",
				cfg.Driver, cfg.DatabaseDirectory, directory)
		}
		if cfg.Migration || cfg.Seeder {
			t.Errorf("production SQLite startup tasks = migration:%t seeder:%t, want both disabled",
				cfg.Migration, cfg.Seeder)
		}
		if cfg.MaxOpenConnections != 1 || cfg.MaxIdleConnections != 1 || cfg.ConnectionMaxLifetime != 0 {
			t.Errorf("production SQLite connection pool = (%d, %d, %s), want (1, 1, no lifetime limit)",
				cfg.MaxOpenConnections, cfg.MaxIdleConnections, cfg.ConnectionMaxLifetime)
		}
	})

	for _, host := range []string{"", "localhost", "127.0.0.1", "::1"} {
		t.Run("preserves production SQLite host "+host, func(t *testing.T) {
			setValidProductionEnvironment(t)
			t.Setenv("GOREST_DB_DRIVER", "sqlite")
			t.Setenv("GOREST_DB_DIRECTORY", "/var/lib/gorest")
			t.Setenv("GOREST_DB_HOST", host)

			cfg, err := LoadDatabase()
			if err != nil {
				t.Fatalf("LoadDatabase() for local SQLite host %q: %v", host, err)
			}
			if cfg.Host != host {
				t.Errorf("LoadDatabase().Host = %q, want %q", cfg.Host, host)
			}
		})
	}
}
