package config

import (
	"os"
	"testing"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest/filter"
	gorestserver "github.com/elsyahtech/gorest/server"
	"github.com/gofiber/fiber/v3/middleware/compress"
)

var configurationEnvironmentKeys = []string{
	"GOREST_APP_NAME",
	"GOREST_APP_VERSION",
	"GOREST_APP_ENV",
	"GOREST_APP_TIMEZONE",
	"GOREST_APP_DEBUG",
	"GOREST_APP_CHECK_FOR_UPDATES",
	"GOREST_LOG_DIR",
	"GOREST_LOG_FILE",
	"GOREST_LOG_CALLER",
	"GOREST_LOG_MAX_FILE_SIZE",
	"GOREST_LOG_MAX_BACKUPS",
	"GOREST_LOG_MAX_AGE",
	"GOREST_SERVER_HOST",
	"GOREST_SERVER_PORT",
	"GOREST_SERVER_BODY_LIMIT",
	"GOREST_SERVER_READ_TIMEOUT",
	"GOREST_SERVER_WRITE_TIMEOUT",
	"GOREST_SERVER_IDLE_TIMEOUT",
	"GOREST_SERVER_HEADER",
	"GOREST_SERVER_STRICT_ROUTING",
	"GOREST_SERVER_CASE_SENSITIVE",
	"GOREST_SERVER_SKIP_UNMATCHED_ROUTES",
	"GOREST_SERVER_DISABLE_HEAD_AUTO_REGISTER",
	"GOREST_SERVER_DISABLE_FLASH_MESSAGES",
	"GOREST_SERVER_IMMUTABLE",
	"GOREST_SERVER_UNESCAPE_PATH",
	"GOREST_SERVER_MAX_RANGES",
	"GOREST_SERVER_CONCURRENCY",
	"GOREST_SERVER_VIEWS_LAYOUT",
	"GOREST_SERVER_PASS_LOCALS_TO_VIEWS",
	"GOREST_SERVER_PASS_LOCALS_TO_CONTEXT",
	"GOREST_SERVER_READ_BUFFER_SIZE",
	"GOREST_SERVER_WRITE_BUFFER_SIZE",
	"GOREST_SERVER_PROXY_HEADER",
	"GOREST_SERVER_GET_ONLY",
	"GOREST_SERVER_DISABLE_KEEPALIVE",
	"GOREST_SERVER_DISABLE_DEFAULT_DATE",
	"GOREST_SERVER_DISABLE_DEFAULT_CONTENT_TYPE",
	"GOREST_SERVER_DISABLE_HEADER_NORMALIZING",
	"GOREST_SERVER_SHARED_STATE_PREFIX",
	"GOREST_SERVER_STREAM_REQUEST_BODY",
	"GOREST_SERVER_DISABLE_PRE_PARSE_MULTIPART_FORM",
	"GOREST_SERVER_REDUCE_MEMORY_USAGE",
	"GOREST_SERVER_TRUST_PROXY",
	"GOREST_SERVER_TRUST_PROXY_PROXIES",
	"GOREST_SERVER_TRUST_PROXY_LINK_LOCAL",
	"GOREST_SERVER_TRUST_PROXY_LOOPBACK",
	"GOREST_SERVER_TRUST_PROXY_PRIVATE",
	"GOREST_SERVER_TRUST_PROXY_UNIX_SOCKET",
	"GOREST_SERVER_ENABLE_IP_VALIDATION",
	"GOREST_SERVER_REQUEST_METHODS",
	"GOREST_SERVER_ENABLE_SPLITTING_ON_PARSERS",
	"GOREST_SERVER_LISTENER_NETWORK",
	"GOREST_SERVER_CERT_FILE",
	"GOREST_SERVER_CERT_KEY_FILE",
	"GOREST_SERVER_CERT_CLIENT_FILE",
	"GOREST_SERVER_UNIX_SOCKET_PATH",
	"GOREST_SERVER_PREFORK_RECOVER_INTERVAL",
	"GOREST_SERVER_SHUTDOWN_TIMEOUT",
	"GOREST_SERVER_PREFORK_SHUTDOWN_GRACE_PERIOD",
	"GOREST_SERVER_PREFORK_RECOVER_THRESHOLD",
	"GOREST_SERVER_UNIX_SOCKET_FILE_MODE",
	"GOREST_SERVER_TLS_MIN_VERSION",
	"GOREST_SERVER_ENABLE_PREFORK",
	"GOREST_SERVER_ENABLE_PRINT_ROUTES",
	"GOREST_SERVER_RATE_LIMIT_MAX",
	"GOREST_SERVER_RATE_LIMIT_EXPIRATION",
	"GOREST_SERVER_RATE_LIMIT_SKIP_FAILED_REQUESTS",
	"GOREST_SERVER_RATE_LIMIT_SKIP_SUCCESSFUL_REQUESTS",
	"GOREST_SERVER_RATE_LIMIT_DISABLE_HEADERS",
	"GOREST_SERVER_RATE_LIMIT_DISABLE_VALUE_REDACTION",
	"GOREST_SERVER_RECOVER_ENABLE_STACK_TRACE",
	"GOREST_SERVER_REQUEST_ID_HEADER",
	"GOREST_SERVER_COMPRESS_LEVEL",
	"GOREST_SERVER_HELMET_XSS_PROTECTION",
	"GOREST_SERVER_HELMET_CONTENT_TYPE_NOSNIFF",
	"GOREST_SERVER_HELMET_X_FRAME_OPTIONS",
	"GOREST_SERVER_HELMET_CONTENT_SECURITY_POLICY",
	"GOREST_SERVER_HELMET_REFERRER_POLICY",
	"GOREST_SERVER_HELMET_PERMISSION_POLICY",
	"GOREST_SERVER_HELMET_CROSS_ORIGIN_EMBEDDER_POLICY",
	"GOREST_SERVER_HELMET_CROSS_ORIGIN_OPENER_POLICY",
	"GOREST_SERVER_HELMET_CROSS_ORIGIN_RESOURCE_POLICY",
	"GOREST_SERVER_HELMET_ORIGIN_AGENT_CLUSTER",
	"GOREST_SERVER_HELMET_X_DNS_PREFETCH_CONTROL",
	"GOREST_SERVER_HELMET_X_DOWNLOAD_OPTIONS",
	"GOREST_SERVER_HELMET_X_PERMITTED_CROSS_DOMAIN",
	"GOREST_SERVER_HELMET_HSTS_MAX_AGE",
	"GOREST_SERVER_HELMET_HSTS_EXCLUDE_SUBDOMAINS",
	"GOREST_SERVER_HELMET_CSP_REPORT_ONLY",
	"GOREST_SERVER_HELMET_HSTS_PRELOAD_ENABLED",
	"GOREST_SERVER_CORS_ORIGINS",
	"GOREST_SERVER_CORS_METHODS",
	"GOREST_SERVER_CORS_HEADERS",
	"GOREST_SERVER_CORS_EXPOSE_HEADERS",
	"GOREST_SERVER_CORS_MAX_AGE",
	"GOREST_SERVER_CORS_ALLOW_CREDENTIALS",
	"GOREST_SERVER_CORS_ALLOW_PRIVATE_NETWORK",
	"GOREST_SERVER_CORS_DISABLE_VALUE_REDACTION",
	"GOREST_DB_DISABLE",
	"GOREST_DB_DRIVER",
	"GOREST_DB_HOST",
	"GOREST_DB_PORT",
	"GOREST_DB_NAME",
	"GOREST_DB_USER",
	"GOREST_DB_PASSWORD",
	"GOREST_DB_TLS",
	"GOREST_DB_ALLOW_INSECURE_PRODUCTION",
	"GOREST_DB_CERT_PATH",
	"GOREST_DB_KEY_PATH",
	"GOREST_DB_CA_PATH",
	"GOREST_DB_CHARSET",
	"GOREST_DB_TIMEOUT",
	"GOREST_DB_CONNECTION_MAX_LIFETIME",
	"GOREST_DB_MAX_OPEN_CONNECTIONS",
	"GOREST_DB_MAX_IDLE_CONNECTIONS",
	"GOREST_DB_QUERY_LOGGING",
	"GOREST_DB_LIB_DIR",
	"GOREST_DB_DIRECTORY",
	"GOREST_DB_MIGRATION",
	"GOREST_DB_SEEDER",
	"GOREST_DB_MIGRATIONS_DIRECTORY",
	"GOREST_DB_SEEDERS_DIRECTORY",
	"GOREST_REDIS_DISABLE",
	"GOREST_REDIS_TLS",
	"GOREST_REDIS_DATABASE",
	"GOREST_REDIS_PORT",
	"GOREST_REDIS_HOST",
	"GOREST_REDIS_USERNAME",
	"GOREST_REDIS_PASSWORD",
	"GOREST_REDIS_USER",
	"GOREST_REDIS_CERT_PATH",
	"GOREST_REDIS_KEY_PATH",
	"GOREST_REDIS_CA_PATH",
	"GOREST_REDIS_SSL_CA_PATH",
	"GOREST_REDIS_SSL_CERT_PATH",
	"GOREST_REDIS_SSL_KEY_PATH",
	"GOREST_FILTER_SIGNATURE_KEY",
	"GOREST_FILTER_SIGNING_METHOD",
	"GOREST_FILTER_ACCESS_TOKEN_DURATION",
	"GOREST_FILTER_REFRESH_TOKEN_DURATION",
	"GOREST_FILTER_COOKIE_SAME_SITE",
	"GOREST_FILTER_COOKIE_DOMAIN",
	"GOREST_FILTER_COOKIE_PATH",
	"GOREST_FILTER_COOKIE_MAX_AGE",
	"GOREST_FILTER_COOKIE_SECURE",
	"GOREST_FILTER_COOKIE_HTTP_ONLY",
	"GOREST_FILTER_COOKIE_SESSION_ONLY",
	"GOREST_MONITORING_ALLOWED_IPS",
	"GOREST_MONITORING_DISABLE",
	"GOREST_MONITORING_ENABLE_HTTP_METRICS",
}

func TestStrLoadsEnvironmentValueAndUsesDefaultWhenEmpty(t *testing.T) {
	t.Setenv("CONFIG_TEST_STRING", "  configured value  ")
	if got := str("CONFIG_TEST_STRING", "fallback"); got != "configured value" {
		t.Fatalf("str() = %q, want trimmed environment value", got)
	}

	t.Setenv("CONFIG_TEST_STRING", "  ")
	if got := str("CONFIG_TEST_STRING", "fallback"); got != "fallback" {
		t.Fatalf("str() with blank environment value = %q, want fallback", got)
	}
}

func TestStrToIntLoadsEnvironmentValueAndUsesDefaultWhenInvalid(t *testing.T) {
	t.Setenv("CONFIG_TEST_INTEGER", " 42 ")
	if got := strToInt("CONFIG_TEST_INTEGER", 10); got != 42 {
		t.Fatalf("strToInt() = %d, want 42", got)
	}

	t.Setenv("CONFIG_TEST_INTEGER", "invalid")
	if got := strToInt("CONFIG_TEST_INTEGER", 10); got != 10 {
		t.Fatalf("strToInt() with invalid environment value = %d, want default 10", got)
	}

	t.Setenv("CONFIG_TEST_INTEGER", " ")
	if got := strToInt("CONFIG_TEST_INTEGER", 10); got != 10 {
		t.Fatalf("strToInt() with blank environment value = %d, want default 10", got)
	}
}

func TestStrToTimeDurationLoadsEnvironmentValueAndUsesDefaultWhenInvalid(t *testing.T) {
	t.Setenv("CONFIG_TEST_DURATION", " 2m ")
	if got := strToTimeDuration("CONFIG_TEST_DURATION", time.Second); got != 2*time.Minute {
		t.Fatalf("strToTimeDuration() = %s, want 2m", got)
	}

	t.Setenv("CONFIG_TEST_DURATION", "invalid")
	if got := strToTimeDuration("CONFIG_TEST_DURATION", time.Second); got != time.Second {
		t.Fatalf("strToTimeDuration() with invalid value = %s, want default 1s", got)
	}

	t.Setenv("CONFIG_TEST_DURATION", " ")
	if got := strToTimeDuration("CONFIG_TEST_DURATION", time.Second); got != time.Second {
		t.Fatalf("strToTimeDuration() with blank value = %s, want default 1s", got)
	}
}

func TestToStrArrayUsesEnvironmentValueAndDefaults(t *testing.T) {
	const key = "CONFIG_TEST_STRING_ARRAY"
	previous, wasSet := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("unset %s: %v", key, err)
	}
	t.Cleanup(func() {
		var err error
		if wasSet {
			err = os.Setenv(key, previous)
		} else {
			err = os.Unsetenv(key)
		}
		if err != nil {
			t.Errorf("restore %s: %v", key, err)
		}
	})

	defaultValue := []string{"default"}
	if got := toStrArray(key, defaultValue); len(got) != 1 || got[0] != "default" {
		t.Fatalf("toStrArray() with unset variable = %v, want %v", got, defaultValue)
	}

	t.Setenv(key, " first, second ,, third ")
	got := toStrArray(key, defaultValue)
	want := []string{"first", "second", "third"}
	if len(got) != len(want) {
		t.Fatalf("toStrArray() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("toStrArray() = %v, want %v", got, want)
		}
	}

	t.Setenv(key, " ")
	if got := toStrArray(key, defaultValue); got == nil || len(got) != 0 {
		t.Fatalf("toStrArray() with explicitly empty variable = %v, want empty non-nil slice", got)
	}
}

func TestLoadFiltersUsesDefaultConfiguration(t *testing.T) {
	isolateConfigurationEnvironment(t)

	cfg, err := FilterConfig()
	if err != nil {
		t.Fatalf("LoadFilters() error = %v", err)
	}

	if cfg.SignatureKey == nil || *cfg.SignatureKey != "" {
		got := "<nil>"
		if cfg.SignatureKey != nil {
			got = *cfg.SignatureKey
		}
		t.Fatalf("FilterConfig().SignatureKey = %q, want an empty development placeholder", got)
	}
	if cfg.SigningMethod != "HS256" {
		t.Errorf("LoadFilters().SigningMethod = %q, want HS256", cfg.SigningMethod)
	}
	if cfg.AccessTokenDuration != 15*time.Minute {
		t.Errorf("LoadFilters().AccessTokenDuration = %s, want 15m", cfg.AccessTokenDuration)
	}
	if cfg.RefreshTokenDuration != 30*24*time.Hour {
		t.Errorf("FilterConfig().RefreshTokenDuration = %s, want 720h", cfg.RefreshTokenDuration)
	}
	if cfg.Issuer != "gorest" {
		t.Errorf("FilterConfig().Issuer = %q, want gorest", cfg.Issuer)
	}
	if cfg.CookieSameSite != "Lax" || cfg.CookiePath != "/" || cfg.CookieDomain != "" {
		t.Errorf("LoadFilters() cookie scope = (%q, %q, %q), want (Lax, /, empty)",
			cfg.CookieSameSite, cfg.CookiePath, cfg.CookieDomain)
	}
	if cfg.CookieSecure || !cfg.CookieHTTPOnly || cfg.CookieSessionOnly {
		t.Errorf("LoadFilters() cookie flags = (secure:%t, httpOnly:%t, sessionOnly:%t), want (true, true, false)",
			cfg.CookieSecure, cfg.CookieHTTPOnly, cfg.CookieSessionOnly)
	}
	if cfg.LoggerConfig.DisableLogging ||
		len(cfg.LoggerConfig.EnabledEvents) != 1 ||
		cfg.LoggerConfig.EnabledEvents[0] != filter.AuthenticationEventMissingCredentials {
		t.Errorf("LoadFilters().LoggerConfig = %+v, want missing-credentials logging enabled", cfg.LoggerConfig)
	}
}

func TestAppConfigUsesEnvironmentFields(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_APP_NAME", "mysql-starter")
	t.Setenv("GOREST_APP_VERSION", "v2.3.4")
	t.Setenv("GOREST_APP_ENV", gorest.Production)
	t.Setenv("GOREST_APP_TIMEZONE", "Asia/Jakarta")
	t.Setenv("GOREST_APP_DEBUG", "true")
	t.Setenv("GOREST_APP_CHECK_FOR_UPDATES", "false")

	cfg, err := AppConfig()
	if err != nil {
		t.Fatalf("AppConfig() error = %v", err)
	}
	if cfg.AppName != "mysql-starter" || cfg.Version != "v2.3.4" ||
		cfg.Environment != gorest.Production || cfg.Timezone != "Asia/Jakarta" ||
		!cfg.Debug || cfg.CheckForUpdates == nil || *cfg.CheckForUpdates {
		t.Errorf("AppConfig() = %+v, want configured app fields and updates disabled", cfg)
	}
}

func TestLoggerConfigUsesEnvironmentFields(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_LOG_DIR", "./custom-logs")
	t.Setenv("GOREST_LOG_FILE", "service.log")
	t.Setenv("GOREST_LOG_CALLER", "true")
	t.Setenv("GOREST_LOG_MAX_FILE_SIZE", "128")
	t.Setenv("GOREST_LOG_MAX_BACKUPS", "9")
	t.Setenv("GOREST_LOG_MAX_AGE", "14")

	cfg, err := LoggerConfig()
	if err != nil {
		t.Fatalf("LoggerConfig() error = %v", err)
	}
	if cfg.LogDir != "./custom-logs" || cfg.LogFile != "service.log" ||
		!cfg.Caller || cfg.MaxFileSize != 128 || cfg.MaxBackups != 9 || cfg.MaxAge != 14 {
		t.Errorf("LoggerConfig() = %+v, want configured logger fields", cfg)
	}
}

func TestLoadFiltersUsesSignatureKeyEnvironmentVariable(t *testing.T) {
	isolateConfigurationEnvironment(t)
	const signatureKey = "unit-test-signature-key"
	t.Setenv("GOREST_FILTER_SIGNATURE_KEY", signatureKey)

	cfg, err := FilterConfig()
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
	t.Setenv("GOREST_APP_ENV", gorest.Development)

	cfg, err := RedisConfig()
	if err != nil {
		t.Fatalf("LoadRedis() error = %v", err)
	}
	if !cfg.Disable {
		t.Fatal("LoadRedis().Disable = false, want Redis disabled by default")
	}
}

func TestLoadRedisCanBeExplicitlyEnabled(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_APP_ENV", gorest.Development)
	t.Setenv("GOREST_REDIS_DISABLE", "false")

	cfg, err := RedisConfig()
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
			if err := validateEnvironment(mode); err != nil {
				t.Fatalf("validateEnvironment() error = %v", err)
			}
		})
	}

	for _, mode := range []string{"test", "staging"} {
		t.Run("rejects "+mode, func(t *testing.T) {
			if err := validateEnvironment(mode); err == nil {
				t.Fatalf("validateEnvironment() accepted GOREST_APP_ENV=%q", mode)
			}
		})
	}

	t.Run("defaults an empty value to development", func(t *testing.T) {
		isolateConfigurationEnvironment(t)
		cfg, err := AppConfig()
		if err != nil {
			t.Fatalf("AppConfig() error = %v", err)
		}
		if cfg.Environment != gorest.Development {
			t.Fatalf("AppConfig().Environment = %q, want %q", cfg.Environment, gorest.Development)
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

func TestServerConfigMapsConfiguredHostAndPort(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_SERVER_HOST", "127.0.0.2")
	t.Setenv("GOREST_SERVER_PORT", "8080")

	cfg, err := ServerConfig()
	if err != nil {
		t.Fatalf("ServerConfig() error = %v", err)
	}
	if cfg.Host != "127.0.0.2" || cfg.Port != 8080 {
		t.Fatalf("server bind = %s:%d, want 127.0.0.2:8080", cfg.Host, cfg.Port)
	}
}

func TestLoadServerEnvironmentOverrides(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_APP_ENV", gorest.Development)
	t.Setenv("GOREST_SERVER_HOST", "0.0.0.0")
	t.Setenv("GOREST_SERVER_PORT", "8081")
	t.Setenv("GOREST_SERVER_BODY_LIMIT", "1048576")
	t.Setenv("GOREST_SERVER_READ_TIMEOUT", "5s")
	t.Setenv("GOREST_SERVER_WRITE_TIMEOUT", "15s")
	t.Setenv("GOREST_SERVER_IDLE_TIMEOUT", "45s")
	t.Setenv("GOREST_SERVER_RATE_LIMIT_MAX", "25")
	t.Setenv("GOREST_SERVER_RATE_LIMIT_EXPIRATION", "2m")
	t.Setenv("GOREST_SERVER_RATE_LIMIT_SKIP_FAILED_REQUESTS", "true")
	t.Setenv("GOREST_SERVER_RATE_LIMIT_SKIP_SUCCESSFUL_REQUESTS", "true")
	t.Setenv("GOREST_SERVER_RATE_LIMIT_DISABLE_HEADERS", "true")
	t.Setenv("GOREST_SERVER_RATE_LIMIT_DISABLE_VALUE_REDACTION", "true")
	t.Setenv("GOREST_SERVER_RECOVER_ENABLE_STACK_TRACE", "true")
	t.Setenv("GOREST_SERVER_REQUEST_ID_HEADER", "X-Correlation-ID")
	t.Setenv("GOREST_SERVER_COMPRESS_LEVEL", "2")
	t.Setenv("GOREST_SERVER_HELMET_XSS_PROTECTION", "1; mode=block")
	t.Setenv("GOREST_SERVER_HELMET_CONTENT_TYPE_NOSNIFF", "nosniff-test")
	t.Setenv("GOREST_SERVER_HELMET_X_FRAME_OPTIONS", "DENY")
	t.Setenv("GOREST_SERVER_HELMET_CONTENT_SECURITY_POLICY", "default-src 'self'")
	t.Setenv("GOREST_SERVER_HELMET_REFERRER_POLICY", "same-origin")
	t.Setenv("GOREST_SERVER_HELMET_PERMISSION_POLICY", "camera=()")
	t.Setenv("GOREST_SERVER_HELMET_CROSS_ORIGIN_EMBEDDER_POLICY", "unsafe-none")
	t.Setenv("GOREST_SERVER_HELMET_CROSS_ORIGIN_OPENER_POLICY", "same-origin-allow-popups")
	t.Setenv("GOREST_SERVER_HELMET_CROSS_ORIGIN_RESOURCE_POLICY", "same-site")
	t.Setenv("GOREST_SERVER_HELMET_ORIGIN_AGENT_CLUSTER", "?0")
	t.Setenv("GOREST_SERVER_HELMET_X_DNS_PREFETCH_CONTROL", "on")
	t.Setenv("GOREST_SERVER_HELMET_X_DOWNLOAD_OPTIONS", "open")
	t.Setenv("GOREST_SERVER_HELMET_X_PERMITTED_CROSS_DOMAIN", "none")
	t.Setenv("GOREST_SERVER_HELMET_HSTS_MAX_AGE", "31536000")
	t.Setenv("GOREST_SERVER_HELMET_HSTS_EXCLUDE_SUBDOMAINS", "false")
	t.Setenv("GOREST_SERVER_HELMET_CSP_REPORT_ONLY", "true")
	t.Setenv("GOREST_SERVER_HELMET_HSTS_PRELOAD_ENABLED", "true")
	t.Setenv("GOREST_SERVER_CORS_ORIGINS", "https://client.example")
	t.Setenv("GOREST_SERVER_CORS_METHODS", "GET, POST")
	t.Setenv("GOREST_SERVER_CORS_HEADERS", "Content-Type, X-Request-ID")
	t.Setenv("GOREST_SERVER_CORS_EXPOSE_HEADERS", "X-Request-ID")
	t.Setenv("GOREST_SERVER_CORS_MAX_AGE", "600")
	t.Setenv("GOREST_SERVER_CORS_ALLOW_CREDENTIALS", "true")
	t.Setenv("GOREST_SERVER_CORS_ALLOW_PRIVATE_NETWORK", "true")
	t.Setenv("GOREST_SERVER_CORS_DISABLE_VALUE_REDACTION", "true")

	cfg, err := ServerConfig()
	if err != nil {
		t.Fatalf("LoadServer() error = %v", err)
	}
	if cfg.Host != "0.0.0.0" || cfg.Port != 8081 {
		t.Errorf("server bind = %s:%d, want 0.0.0.0:8081", cfg.Host, cfg.Port)
	}
	if cfg.Core.BodyLimit != 1048576 ||
		cfg.Core.ReadTimeout != 5*time.Second ||
		cfg.Core.WriteTimeout != 15*time.Second ||
		cfg.Core.IdleTimeout != 45*time.Second {
		t.Errorf("server core config = %+v, want configured body limit and timeouts", cfg.Core)
	}
	if cfg.Cors.MaxAge != 600 ||
		len(cfg.Cors.AllowOrigins) != 1 || cfg.Cors.AllowOrigins[0] != "https://client.example" ||
		len(cfg.Cors.AllowMethods) != 2 || cfg.Cors.AllowMethods[1] != "POST" ||
		len(cfg.Cors.AllowHeaders) != 2 || cfg.Cors.AllowHeaders[1] != "X-Request-ID" {
		t.Errorf("CORS config = %+v, want configured origins, methods, headers, and max age", cfg.Cors)
	}
	if cfg.Limiter.Max != 25 || cfg.Limiter.Expiration != 2*time.Minute {
		t.Errorf("rate limiter config = %+v, want max 25 and expiration 2m", cfg.Limiter)
	}
	if !cfg.Limiter.SkipFailedRequests || !cfg.Limiter.SkipSuccessfulRequests ||
		!cfg.Limiter.DisableHeaders || !cfg.Limiter.DisableValueRedaction {
		t.Errorf("rate limiter flags = %+v, want configured flags enabled", cfg.Limiter)
	}
	if !cfg.Recover.EnableStackTrace {
		t.Error("recover stack traces are disabled, want enabled")
	}
	if cfg.Requestid.Header != "X-Correlation-ID" {
		t.Errorf("request ID header = %q, want X-Correlation-ID", cfg.Requestid.Header)
	}
	if cfg.Compress.Level != compress.LevelBestCompression {
		t.Errorf("compression level = %d, want %d", cfg.Compress.Level, compress.LevelBestCompression)
	}
	if cfg.Cors.ExposeHeaders[0] != "X-Request-ID" ||
		!cfg.Cors.AllowCredentials || !cfg.Cors.AllowPrivateNetwork || !cfg.Cors.DisableValueRedaction {
		t.Errorf("CORS additional config = %+v, want configured expose header and flags", cfg.Cors)
	}
	if cfg.Helmet.XSSProtection != "1; mode=block" ||
		cfg.Helmet.ContentTypeNosniff != "nosniff-test" ||
		cfg.Helmet.XFrameOptions != "DENY" ||
		cfg.Helmet.ContentSecurityPolicy != "default-src 'self'" ||
		cfg.Helmet.ReferrerPolicy != "same-origin" ||
		cfg.Helmet.PermissionPolicy != "camera=()" ||
		cfg.Helmet.CrossOriginEmbedderPolicy != "unsafe-none" ||
		cfg.Helmet.CrossOriginOpenerPolicy != "same-origin-allow-popups" ||
		cfg.Helmet.CrossOriginResourcePolicy != "same-site" ||
		cfg.Helmet.OriginAgentCluster != "?0" ||
		cfg.Helmet.XDNSPrefetchControl != "on" ||
		cfg.Helmet.XDownloadOptions != "open" ||
		cfg.Helmet.XPermittedCrossDomain != "none" ||
		cfg.Helmet.HSTSMaxAge != 31536000 ||
		cfg.Helmet.HSTSExcludeSubdomains ||
		!cfg.Helmet.CSPReportOnly ||
		!cfg.Helmet.HSTSPreloadEnabled {
		t.Errorf("Helmet config = %+v, want configured header values and flags", cfg.Helmet)
	}
}

func TestLoadServerFiberCoreEnvironmentOverrides(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_APP_NAME", "configured-app")
	t.Setenv("GOREST_SERVER_HEADER", "api-server")
	t.Setenv("GOREST_SERVER_STRICT_ROUTING", "true")
	t.Setenv("GOREST_SERVER_CASE_SENSITIVE", "true")
	t.Setenv("GOREST_SERVER_SKIP_UNMATCHED_ROUTES", "true")
	t.Setenv("GOREST_SERVER_DISABLE_HEAD_AUTO_REGISTER", "true")
	t.Setenv("GOREST_SERVER_DISABLE_FLASH_MESSAGES", "true")
	t.Setenv("GOREST_SERVER_IMMUTABLE", "true")
	t.Setenv("GOREST_SERVER_UNESCAPE_PATH", "true")
	t.Setenv("GOREST_SERVER_MAX_RANGES", "24")
	t.Setenv("GOREST_SERVER_CONCURRENCY", "4096")
	t.Setenv("GOREST_SERVER_VIEWS_LAYOUT", "layouts/main")
	t.Setenv("GOREST_SERVER_PASS_LOCALS_TO_VIEWS", "true")
	t.Setenv("GOREST_SERVER_PASS_LOCALS_TO_CONTEXT", "true")
	t.Setenv("GOREST_SERVER_READ_BUFFER_SIZE", "16384")
	t.Setenv("GOREST_SERVER_WRITE_BUFFER_SIZE", "16384")
	t.Setenv("GOREST_SERVER_PROXY_HEADER", "X-Forwarded-For")
	t.Setenv("GOREST_SERVER_GET_ONLY", "true")
	t.Setenv("GOREST_SERVER_DISABLE_KEEPALIVE", "true")
	t.Setenv("GOREST_SERVER_DISABLE_DEFAULT_DATE", "true")
	t.Setenv("GOREST_SERVER_DISABLE_DEFAULT_CONTENT_TYPE", "true")
	t.Setenv("GOREST_SERVER_DISABLE_HEADER_NORMALIZING", "true")
	t.Setenv("GOREST_SERVER_SHARED_STATE_PREFIX", "api-state")
	t.Setenv("GOREST_SERVER_STREAM_REQUEST_BODY", "true")
	t.Setenv("GOREST_SERVER_DISABLE_PRE_PARSE_MULTIPART_FORM", "true")
	t.Setenv("GOREST_SERVER_REDUCE_MEMORY_USAGE", "true")
	t.Setenv("GOREST_SERVER_TRUST_PROXY", "true")
	t.Setenv("GOREST_SERVER_TRUST_PROXY_PROXIES", "10.0.0.1, 10.0.0.0/8")
	t.Setenv("GOREST_SERVER_TRUST_PROXY_LINK_LOCAL", "true")
	t.Setenv("GOREST_SERVER_TRUST_PROXY_LOOPBACK", "true")
	t.Setenv("GOREST_SERVER_TRUST_PROXY_PRIVATE", "true")
	t.Setenv("GOREST_SERVER_TRUST_PROXY_UNIX_SOCKET", "true")
	t.Setenv("GOREST_SERVER_ENABLE_IP_VALIDATION", "true")
	t.Setenv("GOREST_SERVER_REQUEST_METHODS", "GET,POST,PURGE")
	t.Setenv("GOREST_SERVER_ENABLE_SPLITTING_ON_PARSERS", "true")

	cfg, err := ServerConfig()
	if err != nil {
		t.Fatalf("ServerConfig() error = %v", err)
	}

	core := cfg.Core
	if core.ServerHeader != "api-server" || !core.StrictRouting || !core.CaseSensitive ||
		!core.SkipUnmatchedRoutes || !core.DisableHeadAutoRegister || !core.DisableFlashMessages ||
		!core.Immutable || !core.UnescapePath {
		t.Errorf("Fiber routing/config flags = %+v, want configured flags", core)
	}
	if core.MaxRanges != 24 || core.Concurrency != 4096 ||
		core.ViewsLayout != "layouts/main" || !core.PassLocalsToViews || !core.PassLocalsToContext {
		t.Errorf("Fiber limits/views config = %+v, want configured values", core)
	}
	if core.ReadBufferSize != 16384 || core.WriteBufferSize != 16384 ||
		core.ProxyHeader != "X-Forwarded-For" || !core.GETOnly || !core.DisableKeepalive ||
		!core.DisableDefaultDate || !core.DisableDefaultContentType || !core.DisableHeaderNormalizing {
		t.Errorf("Fiber HTTP config = %+v, want configured values", core)
	}
	if core.AppName != "configured-app" || core.SharedStatePrefix != "api-state" ||
		!core.StreamRequestBody || !core.DisablePreParseMultipartForm || !core.ReduceMemoryUsage {
		t.Errorf("Fiber app/body config = %+v, want configured values", core)
	}
	if !core.TrustProxy || len(core.TrustProxyConfig.Proxies) != 2 ||
		core.TrustProxyConfig.Proxies[0] != "10.0.0.1" ||
		core.TrustProxyConfig.Proxies[1] != "10.0.0.0/8" ||
		!core.TrustProxyConfig.LinkLocal || !core.TrustProxyConfig.Loopback ||
		!core.TrustProxyConfig.Private || !core.TrustProxyConfig.UnixSocket {
		t.Errorf("Fiber trusted proxy config = %+v, want configured proxy values", core.TrustProxyConfig)
	}
	if !core.EnableIPValidation || len(core.RequestMethods) != 3 ||
		core.RequestMethods[0] != "GET" || core.RequestMethods[1] != "POST" ||
		core.RequestMethods[2] != "PURGE" || !core.EnableSplittingOnParsers {
		t.Errorf("Fiber request parsing config = %+v, want configured values", core)
	}
}

func TestLoadServerBlankEnvironmentUsesDefaults(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_APP_ENV", gorest.Development)
	for _, key := range []string{
		"GOREST_SERVER_HOST",
		"GOREST_SERVER_PORT",
		"GOREST_SERVER_BODY_LIMIT",
		"GOREST_SERVER_READ_TIMEOUT",
		"GOREST_SERVER_WRITE_TIMEOUT",
		"GOREST_SERVER_IDLE_TIMEOUT",
		"GOREST_SERVER_RATE_LIMIT_MAX",
		"GOREST_SERVER_RATE_LIMIT_EXPIRATION",
		"GOREST_SERVER_RATE_LIMIT_SKIP_FAILED_REQUESTS",
		"GOREST_SERVER_RATE_LIMIT_SKIP_SUCCESSFUL_REQUESTS",
		"GOREST_SERVER_RATE_LIMIT_DISABLE_HEADERS",
		"GOREST_SERVER_RATE_LIMIT_DISABLE_VALUE_REDACTION",
		"GOREST_SERVER_RECOVER_ENABLE_STACK_TRACE",
		"GOREST_SERVER_REQUEST_ID_HEADER",
		"GOREST_SERVER_COMPRESS_LEVEL",
		"GOREST_SERVER_HELMET_XSS_PROTECTION",
		"GOREST_SERVER_HELMET_CONTENT_TYPE_NOSNIFF",
		"GOREST_SERVER_HELMET_X_FRAME_OPTIONS",
		"GOREST_SERVER_HELMET_CONTENT_SECURITY_POLICY",
		"GOREST_SERVER_HELMET_REFERRER_POLICY",
		"GOREST_SERVER_HELMET_PERMISSION_POLICY",
		"GOREST_SERVER_HELMET_CROSS_ORIGIN_EMBEDDER_POLICY",
		"GOREST_SERVER_HELMET_CROSS_ORIGIN_OPENER_POLICY",
		"GOREST_SERVER_HELMET_CROSS_ORIGIN_RESOURCE_POLICY",
		"GOREST_SERVER_HELMET_ORIGIN_AGENT_CLUSTER",
		"GOREST_SERVER_HELMET_X_DNS_PREFETCH_CONTROL",
		"GOREST_SERVER_HELMET_X_DOWNLOAD_OPTIONS",
		"GOREST_SERVER_HELMET_X_PERMITTED_CROSS_DOMAIN",
		"GOREST_SERVER_HELMET_HSTS_MAX_AGE",
		"GOREST_SERVER_HELMET_HSTS_EXCLUDE_SUBDOMAINS",
		"GOREST_SERVER_HELMET_CSP_REPORT_ONLY",
		"GOREST_SERVER_HELMET_HSTS_PRELOAD_ENABLED",
		"GOREST_SERVER_CORS_ORIGINS",
		"GOREST_SERVER_CORS_METHODS",
		"GOREST_SERVER_CORS_HEADERS",
		"GOREST_SERVER_CORS_EXPOSE_HEADERS",
		"GOREST_SERVER_CORS_MAX_AGE",
		"GOREST_SERVER_CORS_ALLOW_CREDENTIALS",
		"GOREST_SERVER_CORS_ALLOW_PRIVATE_NETWORK",
		"GOREST_SERVER_CORS_DISABLE_VALUE_REDACTION",
	} {
		t.Setenv(key, "")
	}

	cfg, err := ServerConfig()
	if err != nil {
		t.Fatalf("LoadServer() error with blank template values = %v", err)
	}
	if cfg.Host != "127.0.0.1" || cfg.Port != 3000 {
		t.Errorf("server bind = %s:%d, want default 127.0.0.1:3000", cfg.Host, cfg.Port)
	}
	if cfg.Core.BodyLimit != 4*1024*1024 ||
		cfg.Core.ReadTimeout != 10*time.Second ||
		cfg.Core.WriteTimeout != 30*time.Second ||
		cfg.Core.IdleTimeout != 120*time.Second {
		t.Errorf("server core config = %+v, want defaults", cfg.Core)
	}
	if cfg.Limiter.Max != 50 || cfg.Limiter.Expiration != time.Minute {
		t.Errorf("rate limiter config = %+v, want defaults", cfg.Limiter)
	}
	if _, _, err := gorestserver.New(&cfg); err != nil {
		t.Fatalf("server.New() error with blank optional config values = %v", err)
	}
}

func TestDatabaseConfigRequiresCredentialsAndDefaultsToMySQL(t *testing.T) {
	isolateConfigurationEnvironment(t)
	if _, err := DatabaseConfig(); err == nil {
		t.Fatal("DatabaseConfig() succeeded without a MySQL username and password")
	}

	t.Setenv("GOREST_DB_USER", "integration-user")
	t.Setenv("GOREST_DB_PASSWORD", "integration-password")
	cfg, err := DatabaseConfig()
	if err != nil {
		t.Fatalf("DatabaseConfig() error = %v", err)
	}
	if cfg.Driver != "mysql" || cfg.Username != "integration-user" || cfg.Password != "integration-password" {
		t.Errorf("DatabaseConfig() = %+v, want MySQL and configured credentials", cfg)
	}
}

func TestDatabaseConfigRejectsNonMySQLDrivers(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_DB_USER", "integration-user")
	t.Setenv("GOREST_DB_PASSWORD", "integration-password")

	for _, driver := range []string{"postgres", "sqlserver", "oracle", "mongo", "scylla", "sqlite"} {
		t.Run(driver, func(t *testing.T) {
			t.Setenv("GOREST_DB_DRIVER", driver)
			if _, err := DatabaseConfig(); err == nil {
				t.Fatalf("DatabaseConfig() accepted non-MySQL driver %q", driver)
			}
		})
	}
}

func TestDatabaseConfigMapsConnectionPoolFields(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_DB_DRIVER", "mysql")
	t.Setenv("GOREST_DB_USER", "test-user")
	t.Setenv("GOREST_DB_PASSWORD", "test-password")
	t.Setenv("GOREST_DB_TLS", "disable")
	t.Setenv("GOREST_DB_CONNECTION_MAX_LIFETIME", "1h")
	t.Setenv("GOREST_DB_MAX_OPEN_CONNECTIONS", "100")
	t.Setenv("GOREST_DB_MAX_IDLE_CONNECTIONS", "10")

	cfg, err := DatabaseConfig()
	if err != nil {
		t.Fatalf("DatabaseConfig() error = %v", err)
	}
	if cfg.ConnectionMaxLifetime != time.Hour ||
		cfg.MaxOpenConnections != 100 ||
		cfg.MaxIdleConnections != 10 {
		t.Errorf("database pool config = (lifetime:%s, maxOpen:%d, maxIdle:%d), want (1h, 100, 10)",
			cfg.ConnectionMaxLifetime, cfg.MaxOpenConnections, cfg.MaxIdleConnections)
	}
}

func TestDatabaseConfigAllowsMySQLTLSWithoutCAPath(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_DB_DRIVER", "mysql")
	t.Setenv("GOREST_DB_USER", "test-user")
	t.Setenv("GOREST_DB_PASSWORD", "test-password")
	t.Setenv("GOREST_DB_TLS", "true")
	t.Setenv("GOREST_DB_CERT_PATH", "")
	t.Setenv("GOREST_DB_CA_PATH", "")

	cfg, err := DatabaseConfig()
	if err != nil {
		t.Fatalf("DatabaseConfig() with MySQL TLS and no CA path error = %v", err)
	}
	if cfg.TLS != "true" || cfg.SSLCAPath != "" || cfg.SSLCertPath != "" {
		t.Errorf("MySQL TLS fields = (TLS:%q, CA:%q, cert:%q), want TLS enabled without CA or client certificate",
			cfg.TLS, cfg.SSLCAPath, cfg.SSLCertPath)
	}
}
