package config

import (
	"fmt"

	"github.com/elsyahtech/gorest/filter"
	goresthelper "github.com/elsyahtech/gorest/helper"
)

// FilterConfig loads JWT and authentication-cookie settings.
func FilterConfig() (*filter.Config, error) {
	cfg := filter.ConfigDefault

	var err error

	accessTokenDuration, err := goresthelper.ReadEnvironmentDurationWithDefault(
		"GOREST_FILTER_ACCESS_TOKEN_DURATION",
		filter.DefaultAccessTokenDuration,
	)
	if err != nil {
		return &filter.Config{}, fmt.Errorf("%w", err)
	}

	refreshTokenDuration, err := goresthelper.ReadEnvironmentDurationWithDefault(
		"GOREST_FILTER_REFRESH_TOKEN_DURATION",
		filter.DefaultRefreshTokenDuration,
	)
	if err != nil {
		return &filter.Config{}, fmt.Errorf("%w", err)
	}

	cfg.SignatureKey = goresthelper.ReadEnvironmentStringPointerWithDefault(
		"GOREST_FILTER_SIGNATURE_KEY",
		"",
	)

	cfg.SigningMethod = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_FILTER_SIGNING_METHOD",
		filter.DefaultSigningMethod,
	)

	cfg.AccessTokenDuration = accessTokenDuration
	cfg.RefreshTokenDuration = refreshTokenDuration

	cfg.CookieConfig.CookieSameSite = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_FILTER_COOKIE_SAME_SITE",
		filter.DefaultCookieSameSite,
	)

	cfg.CookieConfig.CookieDomain = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_FILTER_COOKIE_DOMAIN",
		filter.DefaultCookieDomain,
	)

	cfg.CookieConfig.CookiePath = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_FILTER_COOKIE_PATH",
		filter.DefaultCookiePath,
	)

	cfg.CookieConfig.CookieSecure, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_FILTER_COOKIE_SECURE",
		filter.DefaultCookieSecure,
	)
	if err != nil {
		return &filter.Config{}, fmt.Errorf("failed to parse GOREST_FILTER_COOKIE_SECURE: %w", err)
	}

	cfg.CookieConfig.CookieHTTPOnly, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_FILTER_COOKIE_HTTP_ONLY",
		filter.DefaultCookieHTTPOnly,
	)
	if err != nil {
		return &filter.Config{}, fmt.Errorf("failed to parse GOREST_FILTER_COOKIE_HTTP_ONLY: %w", err)
	}

	cfg.CookieConfig.CookieSessionOnly, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_FILTER_COOKIE_SESSION_ONLY",
		filter.DefaultCookieSessionOnly,
	)
	if err != nil {
		return &filter.Config{}, fmt.Errorf("failed to parse GOREST_FILTER_COOKIE_SESSION_ONLY: %w", err)
	}

	return &cfg, nil
}
