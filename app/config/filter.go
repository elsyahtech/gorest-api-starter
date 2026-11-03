package config

import (
	"os"

	"github.com/elsyahtech/gorest/filter"
)

// xxx
func FilterConfig() (filter.Config, error) {
	cfg := filter.ConfigDefault

	cfg.SignatureKey = strPtr(os.Getenv("GOREST_FILTER_SIGNATURE_KEY"), "")
	cfg.SigningMethod = str(os.Getenv("GOREST_FILTER_SIGNING_METHOD"), filter.DefaultSigningMethod)
	cfg.AccessTokenDuration = strToTimeDuration(os.Getenv("GOREST_FILTER_ACCESS_TOKEN_DURATION"), filter.DefaultAccessTokenDuration)
	cfg.RefreshTokenDuration = strToTimeDuration(os.Getenv("GOREST_FILTER_REFRESH_TOKEN_DURATION"), filter.DefaultRefreshTokenDuration)
	cfg.CookieSameSite = str(os.Getenv("GOREST_FILTER_COOKIE_SAME_SITE"), filter.DefaultCookieSameSite)
	cfg.CookieDomain = str(os.Getenv("GOREST_FILTER_COOKIE_DOMAIN"), filter.DefaultCookieDomain)
	cfg.CookiePath = str(os.Getenv("GOREST_FILTER_COOKIE_PATH"), filter.DefaultCookiePath)
	cfg.CookieMaxAge = strToInt(os.Getenv("GOREST_FILTER_COOKIE_MAX_AGE"), filter.DefaultCookieMaxAge)
	cfg.CookieSecure = strToBool(os.Getenv("GOREST_FILTER_COOKIE_SECURE"), filter.DefaultCookieSecure)
	cfg.CookieHTTPOnly = strToBool(os.Getenv("GOREST_FILTER_COOKIE_HTTP_ONLY"), filter.DefaultCookieHTTPOnly)
	cfg.CookieSessionOnly = strToBool(os.Getenv("GOREST_FILTER_COOKIE_SESSION_ONLY"), filter.DefaultCookieSessionOnly)

	return cfg, nil
}
