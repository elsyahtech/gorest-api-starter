package config

import (
	"github.com/elsyahtech/gorest/monitoring"
)

// MonitoringConfig loads monitoring endpoint and HTTP metrics settings.
func MonitoringConfig() (monitoring.Config, error) {
	cfg := monitoring.ConfigDefault

	var err error
	if cfg.DisableMonitoringEndpoints, err = strToBool("GOREST_MONITORING_DISABLE", false); err != nil {
		return monitoring.Config{}, err
	}

	cfg.MonitoringAllowedIPs = toStrArray("GOREST_MONITORING_ALLOWED_IPS", []string{"127.0.0.1", "::1"})

	if cfg.EnableHTTPMetricsMonitoring, err = strToBool("GOREST_MONITORING_ENABLE_HTTP_METRICS", false); err != nil {
		return monitoring.Config{}, err
	}

	return cfg, nil
}
