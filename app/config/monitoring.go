package config

import (
	"github.com/elsyahtech/gorest/monitoring"
)

// MonitoringConfig loads monitoring endpoint and HTTP metrics settings.
func MonitoringConfig() (*monitoring.Config, error) {
	cfg := monitoring.ConfigDefault

	var err error
	if cfg.DisableMonitoringEndpoints, err = envStrToBool("GOREST_MONITORING_DISABLE", false); err != nil {
		return &monitoring.Config{}, err
	}

	cfg.MonitoringAllowedIPs = envToStrArray("GOREST_MONITORING_ALLOWED_IPS", []string{"127.0.0.1", "::1"})

	if cfg.EnableHTTPMetricsMonitoring, err = envStrToBool("GOREST_MONITORING_ENABLE_HTTP_METRICS", false); err != nil {
		return &monitoring.Config{}, err
	}

	return &cfg, nil
}
