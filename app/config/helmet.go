package config

import (
	"fmt"

	goresthelper "github.com/elsyahtech/gorest/helper"
	"github.com/gofiber/fiber/v3/middleware/helmet"
)

func helmetConfig() (*helmet.Config, error) {
	hSTSExcludeSubdomains, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_HELMET_HSTS_EXCLUDE_SUBDOMAINS",
		false,
	)
	if err != nil {
		return &helmet.Config{}, err
	}

	cspReportOnly, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_HELMET_CSP_REPORT_ONLY",
		false,
	)
	if err != nil {
		return &helmet.Config{}, err
	}

	hSTSPreloadEnabled, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_HELMET_HSTS_PRELOAD_ENABLED",
		false,
	)
	if err != nil {
		return &helmet.Config{}, err
	}

	hSTSMaxAge, err := goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_SERVER_HELMET_HSTS_MAX_AGE",
		0,
	)
	if err != nil {
		return &helmet.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_HELMET_HSTS_MAX_AGE: %w", err)
	}

	return &helmet.Config{
		XSSProtection:             goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_HELMET_XSS_PROTECTION", ""),
		ContentTypeNosniff:        goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_HELMET_CONTENT_TYPE_NOSNIFF", ""),
		XFrameOptions:             goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_HELMET_X_FRAME_OPTIONS", ""),
		ContentSecurityPolicy:     goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_HELMET_CONTENT_SECURITY_POLICY", ""),
		ReferrerPolicy:            goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_HELMET_REFERRER_POLICY", ""),
		PermissionPolicy:          goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_HELMET_PERMISSION_POLICY", ""),
		CrossOriginEmbedderPolicy: goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_HELMET_CROSS_ORIGIN_EMBEDDER_POLICY", ""),
		CrossOriginOpenerPolicy:   goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_HELMET_CROSS_ORIGIN_OPENER_POLICY", ""),
		CrossOriginResourcePolicy: goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_HELMET_CROSS_ORIGIN_RESOURCE_POLICY", ""),
		OriginAgentCluster:        goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_HELMET_ORIGIN_AGENT_CLUSTER", ""),
		XDNSPrefetchControl:       goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_HELMET_X_DNS_PREFETCH_CONTROL", ""),
		XDownloadOptions:          goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_HELMET_X_DOWNLOAD_OPTIONS", ""),
		XPermittedCrossDomain:     goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_HELMET_X_PERMITTED_CROSS_DOMAIN", ""),
		HSTSMaxAge:                hSTSMaxAge,
		HSTSExcludeSubdomains:     hSTSExcludeSubdomains,
		CSPReportOnly:             cspReportOnly,
		HSTSPreloadEnabled:        hSTSPreloadEnabled,
	}, nil
}
