package config

import (
	"fmt"

	"github.com/elsyahtech/gorest/hardening"
	goresthelper "github.com/elsyahtech/gorest/helper"
)

// HardeningConfig loads the optional runtime hardening report settings.
func HardeningConfig() (*hardening.Config, error) {
	cfg := hardening.ConfigDefault

	var err error

	if cfg.Enable, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_HARDENING_ENABLE",
		cfg.Enable,
	); err != nil {
		return &hardening.Config{}, fmt.Errorf("failed to parse GOREST_HARDENING_ENABLE: %w", err)
	}

	cfg.HardeningAllowedIPs = goresthelper.ReadEnvironmentStringSliceWithDefault(
		"GOREST_HARDENING_ALLOWED_IPS",
		cfg.HardeningAllowedIPs,
	)

	return &cfg, nil
}
