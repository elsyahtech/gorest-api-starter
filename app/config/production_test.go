package config

import (
	"os"
	"testing"

	"github.com/elsyahtech/gorest"
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
		t.Setenv("GOREST_DB_DIRECTORY", "/var/lib/gorest")
		if _, err := LoadDatabase(); err != nil {
			t.Fatalf("LoadDatabase() with explicit SQLite directory error = %v", err)
		}
	})
}
