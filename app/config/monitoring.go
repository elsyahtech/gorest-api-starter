package config

import (
	goresthelper "github.com/elsyahtech/gorest/helper"
	"github.com/elsyahtech/gorest/monitoring"
)

// MonitoringConfig loads monitoring endpoint and HTTP metrics settings.
func MonitoringConfig() (*monitoring.Config, error) {
	cfg := monitoring.ConfigDefault

	var err error

	if cfg.DisableMonitoringEndpoints, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_MONITORING_DISABLE",
		false,
	); err != nil {
		return &monitoring.Config{}, err
	}

	cfg.MonitoringAllowedIPs = goresthelper.ReadEnvironmentStringSliceWithDefault(
		"GOREST_MONITORING_ALLOWED_IPS",
		[]string{"127.0.0.1", "::1"},
	)

	if cfg.EnableHTTPMetricsMonitoring, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_MONITORING_ENABLE_HTTP_METRICS",
		false,
	); err != nil {
		return &monitoring.Config{}, err
	}

	return &cfg, nil
}
