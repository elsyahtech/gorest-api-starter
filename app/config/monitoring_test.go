package config

import "testing"

func TestLoadMonitoringUsesLoopbackAllowlistByDefault(t *testing.T) {
	isolateConfigurationEnvironment(t)

	cfg, err := MonitoringConfig()
	if err != nil {
		t.Fatalf("MonitoringConfig() error = %v", err)
	}
	if len(cfg.MonitoringAllowedIPs) != 2 ||
		cfg.MonitoringAllowedIPs[0] != "127.0.0.1" ||
		cfg.MonitoringAllowedIPs[1] != "::1" {
		t.Errorf("MonitoringConfig().MonitoringAllowedIPs = %v, want IPv4 and IPv6 loopback addresses", cfg.MonitoringAllowedIPs)
	}
	if cfg.DisableMonitoringEndpoints {
		t.Error("MonitoringConfig().DisableMonitoringEndpoints = true, want monitoring routes registered by default")
	}
	if cfg.EnableHTTPMetricsMonitoring {
		t.Error("MonitoringConfig().EnableHTTPMetricsMonitoring = true, want per-route metrics disabled by default")
	}
}

func TestLoadMonitoringCanDenyAllWithEmptyAllowlistOverride(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_MONITORING_ALLOWED_IPS", "")

	cfg, err := MonitoringConfig()
	if err != nil {
		t.Fatalf("LoadMonitoring() error = %v", err)
	}
	if len(cfg.MonitoringAllowedIPs) != 0 {
		t.Errorf("LoadMonitoring().HTTPMetricsAllowedIPs = %v, want an empty allowlist", cfg.MonitoringAllowedIPs)
	}
}

func TestLoadMonitoringAppliesEnvironmentOverrides(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_MONITORING_ALLOWED_IPS", "127.0.0.1/32,::1/128")
	t.Setenv("GOREST_MONITORING_DISABLE", "true")
	t.Setenv("GOREST_MONITORING_ENABLE_HTTP_METRICS", "true")

	cfg, err := MonitoringConfig()
	if err != nil {
		t.Fatalf("LoadMonitoring() error = %v", err)
	}
	if len(cfg.MonitoringAllowedIPs) != 2 ||
		cfg.MonitoringAllowedIPs[0] != "127.0.0.1/32" ||
		cfg.MonitoringAllowedIPs[1] != "::1/128" {
		t.Errorf("LoadMonitoring().HTTPMetricsAllowedIPs = %v, want loopback IPv4 and IPv6 CIDRs", cfg.MonitoringAllowedIPs)
	}
	if !cfg.DisableMonitoringEndpoints {
		t.Error("LoadMonitoring().DisableMonitoringEndpoints = false, want true")
	}
	if !cfg.EnableHTTPMetricsMonitoring {
		t.Error("LoadMonitoring().EnablePerRouteHTTPMetrics = false, want true")
	}
}

func TestLoadMonitoringRejectsInvalidBoolean(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_MONITORING_DISABLE", "sometimes")

	if _, err := MonitoringConfig(); err == nil {
		t.Fatal("MonitoringConfig() succeeded with an invalid boolean")
	}
}
