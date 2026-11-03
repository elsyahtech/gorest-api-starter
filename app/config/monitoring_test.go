package config

import "testing"

func TestLoadMonitoringUsesLoopbackAllowlistByDefault(t *testing.T) {
	isolateConfigurationEnvironment(t)

	cfg, err := LoadMonitoring()
	if err != nil {
		t.Fatalf("LoadMonitoring() error = %v", err)
	}
	if len(cfg.MonitoringAllowedIPs) != 2 ||
		cfg.MonitoringAllowedIPs[0] != "127.0.0.1" ||
		cfg.MonitoringAllowedIPs[1] != "::1" {
		t.Errorf("LoadMonitoring().MonitoringAllowedIPs = %v, want IPv4 and IPv6 loopback addresses", cfg.MonitoringAllowedIPs)
	}
	if cfg.DisableMonitoringEndpoints {
		t.Error("LoadMonitoring().DisableMonitoringEndpoints = true, want monitoring routes registered by default")
	}
	if cfg.EnableHTTPMetricsMonitoring {
		t.Error("LoadMonitoring().EnableHTTPMetricsMonitoring = true, want per-route metrics disabled by default")
	}
}

func TestLoadMonitoringCanDenyAllWithEmptyAllowlistOverride(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_HTTP_METRICS_ALLOWED_IPS", "")

	cfg, err := LoadMonitoring()
	if err != nil {
		t.Fatalf("LoadMonitoring() error = %v", err)
	}
	if len(cfg.MonitoringAllowedIPs) != 0 {
		t.Errorf("LoadMonitoring().MonitoringAllowedIPs = %v, want an empty allowlist", cfg.MonitoringAllowedIPs)
	}
}

func TestLoadMonitoringAppliesEnvironmentOverrides(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_HTTP_METRICS_ALLOWED_IPS", "127.0.0.1/32,::1/128")
	t.Setenv("GOREST_DISABLE_MONITORING_ENDPOINTS", "true")
	t.Setenv("GOREST_ENABLE_PER_ROUTE_HTTP_METRICS", "true")

	cfg, err := LoadMonitoring()
	if err != nil {
		t.Fatalf("LoadMonitoring() error = %v", err)
	}
	if len(cfg.MonitoringAllowedIPs) != 2 ||
		cfg.MonitoringAllowedIPs[0] != "127.0.0.1/32" ||
		cfg.MonitoringAllowedIPs[1] != "::1/128" {
		t.Errorf("LoadMonitoring().MonitoringAllowedIPs = %v, want loopback IPv4 and IPv6 CIDRs", cfg.MonitoringAllowedIPs)
	}
	if !cfg.DisableMonitoringEndpoints {
		t.Error("LoadMonitoring().DisableMonitoringEndpoints = false, want true")
	}
	if !cfg.EnableHTTPMetricsMonitoring {
		t.Error("LoadMonitoring().EnableHTTPMetricsMonitoring = false, want true")
	}
}

func TestLoadMonitoringRejectsInvalidBoolean(t *testing.T) {
	isolateConfigurationEnvironment(t)
	t.Setenv("GOREST_DISABLE_MONITORING_ENDPOINTS", "sometimes")

	if _, err := LoadMonitoring(); err == nil {
		t.Fatal("LoadMonitoring() succeeded with an invalid boolean")
	}
}
