package config

import "testing"

func TestHardeningConfigDefaultsDisabled(t *testing.T) {
	isolateConfigurationEnvironment(t)

	cfg, err := HardeningConfig()
	if err != nil {
		t.Fatalf("HardeningConfig() error = %v", err)
	}
	if cfg.Enable {
		t.Fatal("HardeningConfig().Enable = true, want endpoint disabled by default")
	}
	if len(cfg.HardeningAllowedIPs) != 0 {
		t.Fatalf("HardeningConfig().HardeningAllowedIPs = %v, want an empty allowlist", cfg.HardeningAllowedIPs)
	}
}

func TestHardeningConfigLoadsEnvironmentValues(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_HARDENING_ENABLE", "true")
	t.Setenv("GOREST_HARDENING_ALLOWED_IPS", "127.0.0.1/32, ::1/128")

	cfg, err := HardeningConfig()
	if err != nil {
		t.Fatalf("HardeningConfig() error = %v", err)
	}
	if !cfg.Enable {
		t.Fatal("HardeningConfig().Enable = false, want true")
	}
	if len(cfg.HardeningAllowedIPs) != 2 ||
		cfg.HardeningAllowedIPs[0] != "127.0.0.1/32" ||
		cfg.HardeningAllowedIPs[1] != "::1/128" {
		t.Fatalf("HardeningConfig().HardeningAllowedIPs = %v, want IPv4 and IPv6 loopback CIDRs", cfg.HardeningAllowedIPs)
	}
}

func TestHardeningConfigRejectsInvalidEnableFlag(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_HARDENING_ENABLE", "sometimes")

	if _, err := HardeningConfig(); err == nil {
		t.Fatal("HardeningConfig() succeeded with an invalid enable flag")
	}
}
