package config

import (
	"os"

	"github.com/elsyahtech/gorest/monitoring"
)

func MonitoringConfig() (monitoring.Config, error) {
	cfg := monitoring.ConfigDefault

	cfg.DisableMonitoringEndpoints = strToBool(os.Getenv("GOREST_MONITORING_DISABLE"), false)
	cfg.MonitoringAllowedIPs = toStrArray("GOREST_MONITORING_ALLOWED_IPS", []string{"127.0.0.1", "::1"})
	cfg.EnableHTTPMetricsMonitoring = strToBool(os.Getenv("GOREST_MONITORING_ENABLE_HTTP_METRICS"), false)

	return cfg, nil
}
