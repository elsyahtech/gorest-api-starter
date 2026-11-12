package config

import "testing"

func TestLoadHardeningDefaultsDisabled(t *testing.T) {
	isolateConfigurationEnvironment(t)

	cfg, err := LoadHardening()
	if err != nil {
		t.Fatalf("LoadHardening() error = %v", err)
	}
	if cfg.Enable {
		t.Error("LoadHardening().Enable = true, want hardening endpoint disabled by default")
	}
	if len(cfg.HardeningAllowedIPs) != 0 {
		t.Errorf("LoadHardening().HardeningAllowedIPs = %v, want an empty allowlist", cfg.HardeningAllowedIPs)
	}
}

func TestLoadHardeningAppliesEnvironmentOverrides(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_HARDENING_ENABLE", "true")
	t.Setenv("GOREST_HARDENING_ALLOWED_IPS", "127.0.0.1/32,::1/128")

	cfg, err := LoadHardening()
	if err != nil {
		t.Fatalf("LoadHardening() error = %v", err)
	}
	if !cfg.Enable {
		t.Error("LoadHardening().Enable = false, want true")
	}
	if len(cfg.HardeningAllowedIPs) != 2 ||
		cfg.HardeningAllowedIPs[0] != "127.0.0.1/32" ||
		cfg.HardeningAllowedIPs[1] != "::1/128" {
		t.Errorf("LoadHardening().HardeningAllowedIPs = %v, want loopback IPv4 and IPv6 CIDRs", cfg.HardeningAllowedIPs)
	}
}

func TestLoadHardeningRejectsInvalidBoolean(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_HARDENING_ENABLE", "sometimes")

	if _, err := LoadHardening(); err == nil {
		t.Fatal("LoadHardening() succeeded with an invalid enable flag")
	}
}

func TestLoadHardeningRejectsInvalidAllowlist(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_HARDENING_ALLOWED_IPS", "not-an-ip")

	if _, err := LoadHardening(); err == nil {
		t.Fatal("LoadHardening() succeeded with an invalid IP allowlist entry")
	}
}
