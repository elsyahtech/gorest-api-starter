package config

import (
	"testing"
)

func TestLoadRedisUsesGORESTBooleanDefaultForEmptyValue(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_REDIS_DISABLE", "")

	cfg, err := LoadRedis()
	if err != nil {
		t.Fatalf("LoadRedis() error = %v", err)
	}
	if !cfg.Disable {
		t.Fatal("LoadRedis().Disable = false, want GOREST helper's configured default")
	}
}

func TestLoadServerRejectsExplicitlyEmptyPort(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_SERVER_PORT", "")

	if _, err := LoadServer(); err == nil {
		t.Fatal("LoadServer() accepted an explicitly empty port")
	}
}

func TestLoadServerRejectsPortOutsideValidRange(t *testing.T) {
	for _, value := range []string{"0", "65536", "not-a-port"} {
		t.Run(value, func(t *testing.T) {
			isolateConfigurationEnvironment(t)
			t.Setenv("GOREST_SERVER_PORT", value)

			if _, err := LoadServer(); err == nil {
				t.Fatalf("LoadServer() accepted GOREST_SERVER_PORT=%q", value)
			}
		})
	}
}

func TestLoadRedisRejectsExplicitlyEmptyNumericValues(t *testing.T) {
	for _, key := range []string{"GOREST_REDIS_DATABASE", "GOREST_REDIS_PORT"} {
		t.Run(key, func(t *testing.T) {
			isolateConfigurationEnvironment(t)
			t.Setenv(key, "")

			if _, err := LoadRedis(); err == nil {
				t.Fatalf("LoadRedis() accepted an explicitly empty %s", key)
			}
		})
	}
}

func TestLoadRedisRejectsOutOfRangeNumericValues(t *testing.T) {
	tests := []struct {
		key   string
		value string
	}{
		{key: "GOREST_REDIS_DATABASE", value: "-1"},
		{key: "GOREST_REDIS_PORT", value: "0"},
		{key: "GOREST_REDIS_PORT", value: "65536"},
		{key: "GOREST_REDIS_PORT", value: "not-a-port"},
	}
	for _, test := range tests {
		t.Run(test.key+"="+test.value, func(t *testing.T) {
			isolateConfigurationEnvironment(t)
			t.Setenv(test.key, test.value)

			if _, err := LoadRedis(); err == nil {
				t.Fatalf("LoadRedis() accepted %s=%q", test.key, test.value)
			}
		})
	}
}
