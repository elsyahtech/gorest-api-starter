package config

import (
	"github.com/gofiber/fiber/v3/middleware/helmet"
)

func helmetConfig() (*helmet.Config, error) {
	hSTSExcludeSubdomains, err := envStrToBool("GOREST_SERVER_HELMET_HSTS_EXCLUDE_SUBDOMAINS", false)
	if err != nil {
		return &helmet.Config{}, err
	}

	cspReportOnly, err := envStrToBool("GOREST_SERVER_HELMET_CSP_REPORT_ONLY", false)
	if err != nil {
		return &helmet.Config{}, err
	}

	hSTSPreloadEnabled, err := envStrToBool("GOREST_SERVER_HELMET_HSTS_PRELOAD_ENABLED", false)
	if err != nil {
		return &helmet.Config{}, err
	}

	return &helmet.Config{
		XSSProtection:             envStr("GOREST_SERVER_HELMET_XSS_PROTECTION", ""),
		ContentTypeNosniff:        envStr("GOREST_SERVER_HELMET_CONTENT_TYPE_NOSNIFF", ""),
		XFrameOptions:             envStr("GOREST_SERVER_HELMET_X_FRAME_OPTIONS", ""),
		ContentSecurityPolicy:     envStr("GOREST_SERVER_HELMET_CONTENT_SECURITY_POLICY", ""),
		ReferrerPolicy:            envStr("GOREST_SERVER_HELMET_REFERRER_POLICY", ""),
		PermissionPolicy:          envStr("GOREST_SERVER_HELMET_PERMISSION_POLICY", ""),
		CrossOriginEmbedderPolicy: envStr("GOREST_SERVER_HELMET_CROSS_ORIGIN_EMBEDDER_POLICY", ""),
		CrossOriginOpenerPolicy:   envStr("GOREST_SERVER_HELMET_CROSS_ORIGIN_OPENER_POLICY", ""),
		CrossOriginResourcePolicy: envStr("GOREST_SERVER_HELMET_CROSS_ORIGIN_RESOURCE_POLICY", ""),
		OriginAgentCluster:        envStr("GOREST_SERVER_HELMET_ORIGIN_AGENT_CLUSTER", ""),
		XDNSPrefetchControl:       envStr("GOREST_SERVER_HELMET_X_DNS_PREFETCH_CONTROL", ""),
		XDownloadOptions:          envStr("GOREST_SERVER_HELMET_X_DOWNLOAD_OPTIONS", ""),
		XPermittedCrossDomain:     envStr("GOREST_SERVER_HELMET_X_PERMITTED_CROSS_DOMAIN", ""),
		HSTSMaxAge:                envStrToInt("GOREST_SERVER_HELMET_HSTS_MAX_AGE", 0),
		HSTSExcludeSubdomains:     hSTSExcludeSubdomains,
		CSPReportOnly:             cspReportOnly,
		HSTSPreloadEnabled:        hSTSPreloadEnabled,
	}, nil
}
