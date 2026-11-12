package config

import (
	"fmt"

	"github.com/elsyahtech/gorest/hardening"
	goresthelper "github.com/elsyahtech/gorest/helper"
)

func Hardening() hardening.Config {
	cfg, err := LoadHardening()
	if err != nil {
		panic(err)
	}

	return cfg
}

func LoadHardening() (hardening.Config, error) {
	cfg := hardening.ConfigDefault

	var err error
	if cfg.Enable, err = goresthelper.ReadEnvironmentBooleanWithDefault("GOREST_HARDENING_ENABLE", cfg.Enable); err != nil {
		return hardening.Config{}, err
	}
	cfg.HardeningAllowedIPs = goresthelper.ReadEnvironmentStringSliceWithDefault(
		"GOREST_HARDENING_ALLOWED_IPS",
		cfg.HardeningAllowedIPs,
	)
	if _, err := hardening.New(cfg); err != nil {
		return hardening.Config{}, fmt.Errorf("invalid hardening configuration: %w", err)
	}

	return cfg, nil
}
