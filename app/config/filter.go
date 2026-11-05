package config

import (
	"fmt"

	"github.com/elsyahtech/gorest/filter"
)

// FilterConfig loads JWT and authentication-cookie settings.
func FilterConfig() (filter.Config, error) {
	cfg := filter.ConfigDefault

	var err error

	cfg.SignatureKey = strPtr("GOREST_FILTER_SIGNATURE_KEY", "")
	cfg.SigningMethod = str("GOREST_FILTER_SIGNING_METHOD", filter.DefaultSigningMethod)
	cfg.AccessTokenDuration = strToTimeDuration("GOREST_FILTER_ACCESS_TOKEN_DURATION", filter.DefaultAccessTokenDuration)
	cfg.RefreshTokenDuration = strToTimeDuration("GOREST_FILTER_REFRESH_TOKEN_DURATION", filter.DefaultRefreshTokenDuration)
	cfg.CookieSameSite = str("GOREST_FILTER_COOKIE_SAME_SITE", filter.DefaultCookieSameSite)
	cfg.CookieDomain = str("GOREST_FILTER_COOKIE_DOMAIN", filter.DefaultCookieDomain)
	cfg.CookiePath = str("GOREST_FILTER_COOKIE_PATH", filter.DefaultCookiePath)
	cfg.CookieMaxAge = strToInt("GOREST_FILTER_COOKIE_MAX_AGE", filter.DefaultCookieMaxAge) //nolint:staticcheck

	cfg.CookieSecure, err = strToBool("GOREST_FILTER_COOKIE_SECURE", filter.DefaultCookieSecure)
	if err != nil {
		return filter.Config{}, fmt.Errorf("failed to parse GOREST_FILTER_COOKIE_SECURE: %w", err)
	}

	cfg.CookieHTTPOnly, err = strToBool("GOREST_FILTER_COOKIE_HTTP_ONLY", filter.DefaultCookieHTTPOnly)
	if err != nil {
		return filter.Config{}, fmt.Errorf("failed to parse GOREST_FILTER_COOKIE_HTTP_ONLY: %w", err)
	}

	cfg.CookieSessionOnly, err = strToBool("GOREST_FILTER_COOKIE_SESSION_ONLY", filter.DefaultCookieSessionOnly)
	if err != nil {
		return filter.Config{}, fmt.Errorf("failed to parse GOREST_FILTER_COOKIE_SESSION_ONLY: %w", err)
	}

	return cfg, nil
}
