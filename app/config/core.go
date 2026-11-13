package config

import (
	"fmt"

	"github.com/elsyahtech/gorest"
	goresthelper "github.com/elsyahtech/gorest/helper"
	"github.com/elsyahtech/gorest/server"
	"github.com/gofiber/fiber/v3"
)

func coreConfig() (*fiber.Config, error) {
	strictRouting, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_STRICT_ROUTING",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_STRICT_ROUTING: %w", err)
	}

	caseSensitive, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_CASE_SENSITIVE",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_CASE_SENSITIVE: %w", err)
	}

	skipUnmatchedRoutes, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_SKIP_UNMATCHED_ROUTES",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_SKIP_UNMATCHED_ROUTES: %w", err)
	}

	disableHeadAutoRegister, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_DISABLE_HEAD_AUTO_REGISTER",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_HEAD_AUTO_REGISTER: %w", err)
	}

	disableFlashMessages, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_DISABLE_FLASH_MESSAGES",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_FLASH_MESSAGES: %w", err)
	}

	immutable, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_IMMUTABLE",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_IMMUTABLE: %w", err)
	}

	unescapePath, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_UNESCAPE_PATH",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_UNESCAPE_PATH: %w", err)
	}

	passLocalsToViews, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_PASS_LOCALS_TO_VIEWS",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_PASS_LOCALS_TO_VIEWS: %w", err)
	}

	passLocalsToContext, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_PASS_LOCALS_TO_CONTEXT",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_PASS_LOCALS_TO_CONTEXT: %w", err)
	}

	getOnly, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_GET_ONLY",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_GET_ONLY: %w", err)
	}

	disableKeepalive, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_DISABLE_KEEPALIVE",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_KEEPALIVE: %w", err)
	}

	disableDefaultDate, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_DISABLE_DEFAULT_DATE",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_DEFAULT_DATE: %w", err)
	}

	disableDefaultContentType, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_DISABLE_DEFAULT_CONTENT_TYPE",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_DEFAULT_CONTENT_TYPE: %w", err)
	}

	disableHeaderNormalizing, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_DISABLE_HEADER_NORMALIZING",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_HEADER_NORMALIZING: %w", err)
	}

	streamRequestBody, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_STREAM_REQUEST_BODY",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_STREAM_REQUEST_BODY: %w", err)
	}

	disablePreParseMultipartForm, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_DISABLE_PRE_PARSE_MULTIPART_FORM",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_PRE_PARSE_MULTIPART_FORM: %w", err)
	}

	reduceMemoryUsage, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_REDUCE_MEMORY_USAGE",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_REDUCE_MEMORY_USAGE: %w", err)
	}

	trustProxy, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_TRUST_PROXY",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_TRUST_PROXY: %w", err)
	}

	linkLocal, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_TRUST_PROXY_LINK_LOCAL",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_TRUST_PROXY_LINK_LOCAL: %w", err)
	}

	loopback, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_TRUST_PROXY_LOOPBACK",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_TRUST_PROXY_LOOPBACK: %w", err)
	}

	private, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_TRUST_PROXY_PRIVATE",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_TRUST_PROXY_PRIVATE: %w", err)
	}

	unixSocket, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_TRUST_PROXY_UNIX_SOCKET",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_TRUST_PROXY_UNIX_SOCKET: %w", err)
	}

	enableIPValidation, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_ENABLE_IP_VALIDATION",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_ENABLE_IP_VALIDATION: %w", err)
	}

	enableSplittingOnParsers, err := goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_SERVER_ENABLE_SPLITTING_ON_PARSERS",
		false,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_ENABLE_SPLITTING_ON_PARSERS: %w", err)
	}

	readTimeout, err := goresthelper.ReadEnvironmentDurationWithDefault(
		"GOREST_SERVER_READ_TIMEOUT",
		server.DefaultServerReadTimeout,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("%w", err)
	}

	writeTimeout, err := goresthelper.ReadEnvironmentDurationWithDefault(
		"GOREST_SERVER_WRITE_TIMEOUT",
		server.DefaultServerWriteTimeout,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("%w", err)
	}

	idleTimeout, err := goresthelper.ReadEnvironmentDurationWithDefault(
		"GOREST_SERVER_IDLE_TIMEOUT",
		server.DefaultServerIdleTimeout,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("%w", err)
	}

	bodyLimit, err := goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_SERVER_BODY_LIMIT",
		server.DefaultServerBodyLimit,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_BODY_LIMIT: %w", err)
	}

	maxRanges, err := goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_SERVER_MAX_RANGES",
		fiber.DefaultMaxRanges,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_MAX_RANGES: %w", err)
	}

	concurrency, err := goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_SERVER_CONCURRENCY",
		fiber.DefaultConcurrency,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_CONCURRENCY: %w", err)
	}

	readBufferSize, err := goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_SERVER_READ_BUFFER_SIZE",
		fiber.DefaultReadBufferSize,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_READ_BUFFER_SIZE: %w", err)
	}

	writeBufferSize, err := goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_SERVER_WRITE_BUFFER_SIZE",
		fiber.DefaultWriteBufferSize,
	)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_WRITE_BUFFER_SIZE: %w", err)
	}

	return &fiber.Config{
		ServerHeader:                 goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_HEADER", gorest.DefaultAppName),
		StrictRouting:                strictRouting,
		CaseSensitive:                caseSensitive,
		SkipUnmatchedRoutes:          skipUnmatchedRoutes,
		DisableHeadAutoRegister:      disableHeadAutoRegister,
		DisableFlashMessages:         disableFlashMessages,
		Immutable:                    immutable,
		UnescapePath:                 unescapePath,
		BodyLimit:                    bodyLimit,
		MaxRanges:                    maxRanges,
		Concurrency:                  concurrency,
		ViewsLayout:                  goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_VIEWS_LAYOUT", ""),
		PassLocalsToViews:            passLocalsToViews,
		PassLocalsToContext:          passLocalsToContext,
		ReadTimeout:                  readTimeout,
		WriteTimeout:                 writeTimeout,
		IdleTimeout:                  idleTimeout,
		ReadBufferSize:               readBufferSize,
		WriteBufferSize:              writeBufferSize,
		ProxyHeader:                  goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_PROXY_HEADER", ""),
		GETOnly:                      getOnly,
		DisableKeepalive:             disableKeepalive,
		DisableDefaultDate:           disableDefaultDate,
		DisableDefaultContentType:    disableDefaultContentType,
		DisableHeaderNormalizing:     disableHeaderNormalizing,
		AppName:                      goresthelper.ReadEnvironmentStringWithDefault("GOREST_APP_NAME", gorest.DefaultAppName),
		SharedStatePrefix:            goresthelper.ReadEnvironmentStringWithDefault("GOREST_SERVER_SHARED_STATE_PREFIX", gorest.DefaultAppName),
		StreamRequestBody:            streamRequestBody,
		DisablePreParseMultipartForm: disablePreParseMultipartForm,
		ReduceMemoryUsage:            reduceMemoryUsage,
		TrustProxy:                   trustProxy,
		TrustProxyConfig: fiber.TrustProxyConfig{
			Proxies:    goresthelper.ReadEnvironmentStringSliceWithDefault("GOREST_SERVER_TRUST_PROXY_PROXIES", make([]string, 0)),
			LinkLocal:  linkLocal,
			Loopback:   loopback,
			Private:    private,
			UnixSocket: unixSocket,
		},
		EnableIPValidation:       enableIPValidation,
		RequestMethods:           goresthelper.ReadEnvironmentStringSliceWithDefault("GOREST_SERVER_REQUEST_METHODS", make([]string, 0)),
		EnableSplittingOnParsers: enableSplittingOnParsers,
	}, nil
}
