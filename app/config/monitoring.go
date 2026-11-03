package config

import (
	"github.com/elsyahtech/gorest/monitoring"
)

func Monitoring() monitoring.Config {
	cfg := monitoring.ConfigDefault

	cfg.MonitoringAllowedIPs = []string{"127.0.0.1", "::1"}
	cfg.DisableMonitoringEndpoints = false
	cfg.EnableHTTPMetricsMonitoring = false

	return cfg
}

func LoadMonitoring() (monitoring.Config, error) {
	cfg := Monitoring()

	var err error

	cfg.MonitoringAllowedIPs, err = envStringSlice("GOREST_HTTP_METRICS_ALLOWED_IPS", cfg.MonitoringAllowedIPs)
	if err != nil {
		return monitoring.Config{}, err
	}

	if cfg.DisableMonitoringEndpoints, err = envBool("GOREST_DISABLE_MONITORING_ENDPOINTS", cfg.DisableMonitoringEndpoints); err != nil {
		return monitoring.Config{}, err
	}

	if cfg.EnableHTTPMetricsMonitoring, err = envBool("GOREST_ENABLE_PER_ROUTE_HTTP_METRICS", cfg.EnableHTTPMetricsMonitoring); err != nil {
		return monitoring.Config{}, err
	}

	return cfg, nil
}
