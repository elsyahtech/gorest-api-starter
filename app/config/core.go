package config

import (
	"fmt"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest/server"
	"github.com/gofiber/fiber/v3"
)

func coreConfig() (*fiber.Config, error) {
	strictRouting, err := envStrToBool("GOREST_SERVER_STRICT_ROUTING", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_STRICT_ROUTING: %w", err)
	}

	caseSensitive, err := envStrToBool("GOREST_SERVER_CASE_SENSITIVE", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_CASE_SENSITIVE: %w", err)
	}

	skipUnmatchedRoutes, err := envStrToBool("GOREST_SERVER_SKIP_UNMATCHED_ROUTES", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_SKIP_UNMATCHED_ROUTES: %w", err)
	}

	disableHeadAutoRegister, err := envStrToBool("GOREST_SERVER_DISABLE_HEAD_AUTO_REGISTER", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_HEAD_AUTO_REGISTER: %w", err)
	}

	disableFlashMessages, err := envStrToBool("GOREST_SERVER_DISABLE_FLASH_MESSAGES", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_FLASH_MESSAGES: %w", err)
	}

	immutable, err := envStrToBool("GOREST_SERVER_IMMUTABLE", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_IMMUTABLE: %w", err)
	}

	unescapePath, err := envStrToBool("GOREST_SERVER_UNESCAPE_PATH", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_UNESCAPE_PATH: %w", err)
	}

	passLocalsToViews, err := envStrToBool("GOREST_SERVER_PASS_LOCALS_TO_VIEWS", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_PASS_LOCALS_TO_VIEWS: %w", err)
	}

	passLocalsToContext, err := envStrToBool("GOREST_SERVER_PASS_LOCALS_TO_CONTEXT", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_PASS_LOCALS_TO_CONTEXT: %w", err)
	}

	getOnly, err := envStrToBool("GOREST_SERVER_GET_ONLY", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_GET_ONLY: %w", err)
	}

	disableKeepalive, err := envStrToBool("GOREST_SERVER_DISABLE_KEEPALIVE", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_KEEPALIVE: %w", err)
	}

	disableDefaultDate, err := envStrToBool("GOREST_SERVER_DISABLE_DEFAULT_DATE", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_DEFAULT_DATE: %w", err)
	}

	disableDefaultContentType, err := envStrToBool("GOREST_SERVER_DISABLE_DEFAULT_CONTENT_TYPE", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_DEFAULT_CONTENT_TYPE: %w", err)
	}

	disableHeaderNormalizing, err := envStrToBool("GOREST_SERVER_DISABLE_HEADER_NORMALIZING", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_HEADER_NORMALIZING: %w", err)
	}

	streamRequestBody, err := envStrToBool("GOREST_SERVER_STREAM_REQUEST_BODY", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_STREAM_REQUEST_BODY: %w", err)
	}

	disablePreParseMultipartForm, err := envStrToBool("GOREST_SERVER_DISABLE_PRE_PARSE_MULTIPART_FORM", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_PRE_PARSE_MULTIPART_FORM: %w", err)
	}

	reduceMemoryUsage, err := envStrToBool("GOREST_SERVER_REDUCE_MEMORY_USAGE", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_REDUCE_MEMORY_USAGE: %w", err)
	}

	trustProxy, err := envStrToBool("GOREST_SERVER_TRUST_PROXY", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_TRUST_PROXY: %w", err)
	}

	linkLocal, err := envStrToBool("GOREST_SERVER_TRUST_PROXY_LINK_LOCAL", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_TRUST_PROXY_LINK_LOCAL: %w", err)
	}

	loopback, err := envStrToBool("GOREST_SERVER_TRUST_PROXY_LOOPBACK", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_TRUST_PROXY_LOOPBACK: %w", err)
	}

	private, err := envStrToBool("GOREST_SERVER_TRUST_PROXY_PRIVATE", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_TRUST_PROXY_PRIVATE: %w", err)
	}

	unixSocket, err := envStrToBool("GOREST_SERVER_TRUST_PROXY_UNIX_SOCKET", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_TRUST_PROXY_UNIX_SOCKET: %w", err)
	}

	enableIPValidation, err := envStrToBool("GOREST_SERVER_ENABLE_IP_VALIDATION", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_ENABLE_IP_VALIDATION: %w", err)
	}

	enableSplittingOnParsers, err := envStrToBool("GOREST_SERVER_ENABLE_SPLITTING_ON_PARSERS", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_ENABLE_SPLITTING_ON_PARSERS: %w", err)
	}

	readTimeout, err := envStrToTimeDuration("GOREST_SERVER_READ_TIMEOUT", server.DefaultServerReadTimeout)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("%w", err)
	}

	writeTimeout, err := envStrToTimeDuration("GOREST_SERVER_WRITE_TIMEOUT", server.DefaultServerWriteTimeout)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("%w", err)
	}

	idleTimeout, err := envStrToTimeDuration("GOREST_SERVER_IDLE_TIMEOUT", server.DefaultServerIdleTimeout)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("%w", err)
	}

	return &fiber.Config{
		ServerHeader:                 envStr("GOREST_SERVER_HEADER", gorest.DefaultAppName),
		StrictRouting:                strictRouting,
		CaseSensitive:                caseSensitive,
		SkipUnmatchedRoutes:          skipUnmatchedRoutes,
		DisableHeadAutoRegister:      disableHeadAutoRegister,
		DisableFlashMessages:         disableFlashMessages,
		Immutable:                    immutable,
		UnescapePath:                 unescapePath,
		BodyLimit:                    envStrToInt("GOREST_SERVER_BODY_LIMIT", server.DefaultServerBodyLimit),
		MaxRanges:                    envStrToInt("GOREST_SERVER_MAX_RANGES", fiber.DefaultMaxRanges),
		Concurrency:                  envStrToInt("GOREST_SERVER_CONCURRENCY", fiber.DefaultConcurrency),
		ViewsLayout:                  envStr("GOREST_SERVER_VIEWS_LAYOUT", ""),
		PassLocalsToViews:            passLocalsToViews,
		PassLocalsToContext:          passLocalsToContext,
		ReadTimeout:                  readTimeout,
		WriteTimeout:                 writeTimeout,
		IdleTimeout:                  idleTimeout,
		ReadBufferSize:               envStrToInt("GOREST_SERVER_READ_BUFFER_SIZE", fiber.DefaultReadBufferSize),
		WriteBufferSize:              envStrToInt("GOREST_SERVER_WRITE_BUFFER_SIZE", fiber.DefaultWriteBufferSize),
		ProxyHeader:                  envStr("GOREST_SERVER_PROXY_HEADER", ""),
		GETOnly:                      getOnly,
		DisableKeepalive:             disableKeepalive,
		DisableDefaultDate:           disableDefaultDate,
		DisableDefaultContentType:    disableDefaultContentType,
		DisableHeaderNormalizing:     disableHeaderNormalizing,
		AppName:                      envStr("GOREST_APP_NAME", gorest.DefaultAppName),
		SharedStatePrefix:            envStr("GOREST_SERVER_SHARED_STATE_PREFIX", gorest.DefaultAppName),
		StreamRequestBody:            streamRequestBody,
		DisablePreParseMultipartForm: disablePreParseMultipartForm,
		ReduceMemoryUsage:            reduceMemoryUsage,
		TrustProxy:                   trustProxy,
		TrustProxyConfig: fiber.TrustProxyConfig{
			Proxies:    envToStrArray("GOREST_SERVER_TRUST_PROXY_PROXIES", make([]string, 0)),
			LinkLocal:  linkLocal,
			Loopback:   loopback,
			Private:    private,
			UnixSocket: unixSocket,
		},
		EnableIPValidation:       enableIPValidation,
		RequestMethods:           envToStrArray("GOREST_SERVER_REQUEST_METHODS", make([]string, 0)),
		EnableSplittingOnParsers: enableSplittingOnParsers,
	}, nil
}
