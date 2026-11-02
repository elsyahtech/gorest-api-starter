package config

import (
	"github.com/elsyahtech/gorest/monitoring"
)

func Monitoring() monitoring.Config {
	cfg := monitoring.ConfigDefault

	cfg.HTTPMetricsAllowedIPs = []string{"127.0.0.1", "::1"}
	cfg.DisableMonitoringEndpoints = false
	cfg.EnablePerRouteHTTPMetrics = false

	return cfg
}

func LoadMonitoring() (monitoring.Config, error) {
	cfg := Monitoring()

	var err error

	cfg.HTTPMetricsAllowedIPs, err = envStringSlice("GOREST_HTTP_METRICS_ALLOWED_IPS", cfg.HTTPMetricsAllowedIPs)
	if err != nil {
		return monitoring.Config{}, err
	}

	if cfg.DisableMonitoringEndpoints, err = envBool("GOREST_DISABLE_MONITORING_ENDPOINTS", cfg.DisableMonitoringEndpoints); err != nil {
		return monitoring.Config{}, err
	}

	if cfg.EnablePerRouteHTTPMetrics, err = envBool("GOREST_ENABLE_PER_ROUTE_HTTP_METRICS", cfg.EnablePerRouteHTTPMetrics); err != nil {
		return monitoring.Config{}, err
	}

	return cfg, nil
}
