package config

import (
	"os"
	"time"

	"github.com/elsyahtech/gorest/filter"
)

func Filters() filter.Config {
	cfg, err := LoadFilters()
	if err != nil {
		panic(err)
	}

	return cfg
}

func LoadFilters() (filter.Config, error) {
	signatureKey := "your-signature-key" // Replace with your actual signature key

	if value, ok := os.LookupEnv("GOREST_FILTER_SIGNATURE_KEY"); ok {
		signatureKey = value
	}

	return filter.Config{
		SignatureKey:         &signatureKey,      // Replace with your actual signature key
		SigningMethod:        "HS256",            // Replace with your desired signing method
		AccessTokenDuration:  15 * time.Minute,   // 15 minutes in seconds
		RefreshTokenDuration: 7 * 24 * time.Hour, // 7 days in seconds
		Issuer:               "your-issuer",      // Replace with your actual issuer
		CookieSameSite:       "Lax",              // Replace with your desired SameSite value
		CookieDomain:         "",                 // Replace with your actual cookie domain
		CookiePath:           "/",                // Replace with your desired cookie path
		CookieSecure:         true,               // Set to true if using HTTPS
		CookieHTTPOnly:       true,               // Set to true for security
		CookieSessionOnly:    false,              // Set to true if you want session-only cookies
		// Replace with your desired logger configuration. See the detail in
		// https://github.com/elsyahtech/gorest/blob/main/docs/guide/filter.md#authentication-failure-logging
		LoggerConfig: filter.LoggerConfig{
			DisableLogging: false,
			EnabledEvents:  []filter.AuthenticationLogEvent{filter.AuthenticationEventMissingCredentials},
		},
	}, nil
}
