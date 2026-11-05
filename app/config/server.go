package config

import (
	"fmt"
	"os"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest/server"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

// ServerConfig loads HTTP server, middleware, and listener settings.
// Authentication and authorization must be configured separately for protected routes.
func ServerConfig() (server.Config, error) {
	cfg := server.ConfigDefault

	cfg.Host = str("GOREST_SERVER_HOST", server.DefaultServerHost)
	cfg.Port = strToInt("GOREST_SERVER_PORT", server.DefaultServerPort)

	stackTrace, err := strToBool("GOREST_SERVER_RECOVER_ENABLE_STACK_TRACE", false)
	if err != nil {
		return server.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_RECOVER_ENABLE_STACK_TRACE: %w", err)
	}

	cfg.Recover = &recover.Config{
		EnableStackTrace: stackTrace,
	}

	cfg.Requestid = &requestid.Config{
		Header: os.Getenv("GOREST_SERVER_REQUEST_ID_HEADER"),
	}

	cfg.Compress = &compress.Config{
		Level: compress.Level(strToInt("GOREST_SERVER_COMPRESS_LEVEL", int(compress.LevelDefault))),
	}

	cfg.Core, err = fiberCoreConfig()
	if err != nil {
		return server.Config{}, fmt.Errorf("failed to configure fiber core: %w", err)
	}

	origins := toStrArray("GOREST_SERVER_CORS_ORIGINS", []string{"http://localhost:3000", "http://127.0.0.1:3000"})
	if len(origins) == 0 {
		origins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}

	methods := toStrArray("GOREST_SERVER_CORS_METHODS", []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	if len(methods) == 0 {
		methods = []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}

	headers := toStrArray("GOREST_SERVER_CORS_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"})
	if len(headers) == 0 {
		headers = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"}
	}

	allowCredentials, err := strToBool("GOREST_SERVER_CORS_ALLOW_CREDENTIALS", false)
	if err != nil {
		return server.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_CORS_ALLOW_CREDENTIALS: %w", err)
	}

	allowPrivateNetwork, err := strToBool("GOREST_SERVER_CORS_ALLOW_PRIVATE_NETWORK", false)
	if err != nil {
		return server.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_CORS_ALLOW_PRIVATE_NETWORK: %w", err)
	}

	disableCORSValueRedaction, err := strToBool("GOREST_SERVER_CORS_DISABLE_VALUE_REDACTION", false)
	if err != nil {
		return server.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_CORS_DISABLE_VALUE_REDACTION: %w", err)
	}

	cfg.Cors = &cors.Config{
		AllowOrigins:          origins,
		AllowMethods:          methods,
		AllowHeaders:          headers,
		ExposeHeaders:         toStrArray("GOREST_SERVER_CORS_EXPOSE_HEADERS", make([]string, 0)),
		MaxAge:                strToInt("GOREST_SERVER_CORS_MAX_AGE", 300),
		AllowCredentials:      allowCredentials,
		AllowPrivateNetwork:   allowPrivateNetwork,
		DisableValueRedaction: disableCORSValueRedaction,
	}

	hSTSExcludeSubdomains, err := strToBool("GOREST_SERVER_HELMET_HSTS_EXCLUDE_SUBDOMAINS", false)
	if err != nil {
		return server.Config{}, err
	}

	cspReportOnly, err := strToBool("GOREST_SERVER_HELMET_CSP_REPORT_ONLY", false)
	if err != nil {
		return server.Config{}, err
	}

	hSTSPreloadEnabled, err := strToBool("GOREST_SERVER_HELMET_HSTS_PRELOAD_ENABLED", false)
	if err != nil {
		return server.Config{}, err
	}

	cfg.Helmet = &helmet.Config{
		XSSProtection:             os.Getenv("GOREST_SERVER_HELMET_XSS_PROTECTION"),
		ContentTypeNosniff:        os.Getenv("GOREST_SERVER_HELMET_CONTENT_TYPE_NOSNIFF"),
		XFrameOptions:             os.Getenv("GOREST_SERVER_HELMET_X_FRAME_OPTIONS"),
		ContentSecurityPolicy:     os.Getenv("GOREST_SERVER_HELMET_CONTENT_SECURITY_POLICY"),
		ReferrerPolicy:            os.Getenv("GOREST_SERVER_HELMET_REFERRER_POLICY"),
		PermissionPolicy:          os.Getenv("GOREST_SERVER_HELMET_PERMISSION_POLICY"),
		CrossOriginEmbedderPolicy: os.Getenv("GOREST_SERVER_HELMET_CROSS_ORIGIN_EMBEDDER_POLICY"),
		CrossOriginOpenerPolicy:   os.Getenv("GOREST_SERVER_HELMET_CROSS_ORIGIN_OPENER_POLICY"),
		CrossOriginResourcePolicy: os.Getenv("GOREST_SERVER_HELMET_CROSS_ORIGIN_RESOURCE_POLICY"),
		OriginAgentCluster:        os.Getenv("GOREST_SERVER_HELMET_ORIGIN_AGENT_CLUSTER"),
		XDNSPrefetchControl:       os.Getenv("GOREST_SERVER_HELMET_X_DNS_PREFETCH_CONTROL"),
		XDownloadOptions:          os.Getenv("GOREST_SERVER_HELMET_X_DOWNLOAD_OPTIONS"),
		XPermittedCrossDomain:     os.Getenv("GOREST_SERVER_HELMET_X_PERMITTED_CROSS_DOMAIN"),
		HSTSMaxAge:                strToInt("GOREST_SERVER_HELMET_HSTS_MAX_AGE", 0),
		HSTSExcludeSubdomains:     hSTSExcludeSubdomains,
		CSPReportOnly:             cspReportOnly,
		HSTSPreloadEnabled:        hSTSPreloadEnabled,
	}

	skipFailedRequests, err := strToBool("GOREST_SERVER_RATE_LIMIT_SKIP_FAILED_REQUESTS", false)
	if err != nil {
		return server.Config{}, err
	}

	skipSuccessfulRequests, err := strToBool("GOREST_SERVER_RATE_LIMIT_SKIP_SUCCESSFUL_REQUESTS", false)
	if err != nil {
		return server.Config{}, err
	}

	disableHeaders, err := strToBool("GOREST_SERVER_RATE_LIMIT_DISABLE_HEADERS", false)
	if err != nil {
		return server.Config{}, err
	}

	disableValueRedaction, err := strToBool("GOREST_SERVER_RATE_LIMIT_DISABLE_VALUE_REDACTION", false)
	if err != nil {
		return server.Config{}, err
	}

	cfg.Limiter = &limiter.Config{
		Max:                    strToInt("GOREST_SERVER_RATE_LIMIT_MAX", 50),
		Expiration:             strToTimeDuration("GOREST_SERVER_RATE_LIMIT_EXPIRATION", time.Minute),
		SkipFailedRequests:     skipFailedRequests,
		SkipSuccessfulRequests: skipSuccessfulRequests,
		DisableHeaders:         disableHeaders,
		DisableValueRedaction:  disableValueRedaction,
	}

	listener, err := loadServerListenerConfig()
	if err != nil {
		return server.Config{}, err
	}

	cfg.Listener = listener

	return cfg, nil
}

func fiberCoreConfig() (*fiber.Config, error) {
	strictRouting, err := strToBool("GOREST_SERVER_STRICT_ROUTING", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_STRICT_ROUTING: %w", err)
	}

	caseSensitive, err := strToBool("GOREST_SERVER_CASE_SENSITIVE", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_CASE_SENSITIVE: %w", err)
	}

	skipUnmatchedRoutes, err := strToBool("GOREST_SERVER_SKIP_UNMATCHED_ROUTES", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_SKIP_UNMATCHED_ROUTES: %w", err)
	}

	disableHeadAutoRegister, err := strToBool("GOREST_SERVER_DISABLE_HEAD_AUTO_REGISTER", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_HEAD_AUTO_REGISTER: %w", err)
	}

	disableFlashMessages, err := strToBool("GOREST_SERVER_DISABLE_FLASH_MESSAGES", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_FLASH_MESSAGES: %w", err)
	}

	immutable, err := strToBool("GOREST_SERVER_IMMUTABLE", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_IMMUTABLE: %w", err)
	}

	unescapePath, err := strToBool("GOREST_SERVER_UNESCAPE_PATH", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_UNESCAPE_PATH: %w", err)
	}

	passLocalsToViews, err := strToBool("GOREST_SERVER_PASS_LOCALS_TO_VIEWS", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_PASS_LOCALS_TO_VIEWS: %w", err)
	}

	passLocalsToContext, err := strToBool("GOREST_SERVER_PASS_LOCALS_TO_CONTEXT", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_PASS_LOCALS_TO_CONTEXT: %w", err)
	}

	getOnly, err := strToBool("GOREST_SERVER_GET_ONLY", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_GET_ONLY: %w", err)
	}

	disableKeepalive, err := strToBool("GOREST_SERVER_DISABLE_KEEPALIVE", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_KEEPALIVE: %w", err)
	}

	disableDefaultDate, err := strToBool("GOREST_SERVER_DISABLE_DEFAULT_DATE", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_DEFAULT_DATE: %w", err)
	}

	disableDefaultContentType, err := strToBool("GOREST_SERVER_DISABLE_DEFAULT_CONTENT_TYPE", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_DEFAULT_CONTENT_TYPE: %w", err)
	}

	disableHeaderNormalizing, err := strToBool("GOREST_SERVER_DISABLE_HEADER_NORMALIZING", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_HEADER_NORMALIZING: %w", err)
	}

	streamRequestBody, err := strToBool("GOREST_SERVER_STREAM_REQUEST_BODY", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_STREAM_REQUEST_BODY: %w", err)
	}

	disablePreParseMultipartForm, err := strToBool("GOREST_SERVER_DISABLE_PRE_PARSE_MULTIPART_FORM", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_DISABLE_PRE_PARSE_MULTIPART_FORM: %w", err)
	}

	reduceMemoryUsage, err := strToBool("GOREST_SERVER_REDUCE_MEMORY_USAGE", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_REDUCE_MEMORY_USAGE: %w", err)
	}

	trustProxy, err := strToBool("GOREST_SERVER_TRUST_PROXY", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_TRUST_PROXY: %w", err)
	}

	linkLocal, err := strToBool("GOREST_SERVER_TRUST_PROXY_LINK_LOCAL", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_TRUST_PROXY_LINK_LOCAL: %w", err)
	}

	loopback, err := strToBool("GOREST_SERVER_TRUST_PROXY_LOOPBACK", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_TRUST_PROXY_LOOPBACK: %w", err)
	}

	private, err := strToBool("GOREST_SERVER_TRUST_PROXY_PRIVATE", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_TRUST_PROXY_PRIVATE: %w", err)
	}

	unixSocket, err := strToBool("GOREST_SERVER_TRUST_PROXY_UNIX_SOCKET", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_TRUST_PROXY_UNIX_SOCKET: %w", err)
	}

	enableIPValidation, err := strToBool("GOREST_SERVER_ENABLE_IP_VALIDATION", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_ENABLE_IP_VALIDATION: %w", err)
	}

	enableSplittingOnParsers, err := strToBool("GOREST_SERVER_ENABLE_SPLITTING_ON_PARSERS", false)
	if err != nil {
		return &fiber.Config{}, fmt.Errorf("failed to parse GOREST_SERVER_ENABLE_SPLITTING_ON_PARSERS: %w", err)
	}

	return &fiber.Config{
		ServerHeader:                 os.Getenv("GOREST_SERVER_HEADER"),
		StrictRouting:                strictRouting,
		CaseSensitive:                caseSensitive,
		SkipUnmatchedRoutes:          skipUnmatchedRoutes,
		DisableHeadAutoRegister:      disableHeadAutoRegister,
		DisableFlashMessages:         disableFlashMessages,
		Immutable:                    immutable,
		UnescapePath:                 unescapePath,
		BodyLimit:                    strToInt("GOREST_SERVER_BODY_LIMIT", server.DefaultServerBodyLimit),
		MaxRanges:                    strToInt("GOREST_SERVER_MAX_RANGES", fiber.DefaultMaxRanges),
		Concurrency:                  strToInt("GOREST_SERVER_CONCURRENCY", fiber.DefaultConcurrency),
		ViewsLayout:                  os.Getenv("GOREST_SERVER_VIEWS_LAYOUT"),
		PassLocalsToViews:            passLocalsToViews,
		PassLocalsToContext:          passLocalsToContext,
		ReadTimeout:                  strToTimeDuration("GOREST_SERVER_READ_TIMEOUT", server.DefaultServerReadTimeout),
		WriteTimeout:                 strToTimeDuration("GOREST_SERVER_WRITE_TIMEOUT", server.DefaultServerWriteTimeout),
		IdleTimeout:                  strToTimeDuration("GOREST_SERVER_IDLE_TIMEOUT", server.DefaultServerIdleTimeout),
		ReadBufferSize:               strToInt("GOREST_SERVER_READ_BUFFER_SIZE", fiber.DefaultReadBufferSize),
		WriteBufferSize:              strToInt("GOREST_SERVER_WRITE_BUFFER_SIZE", fiber.DefaultWriteBufferSize),
		ProxyHeader:                  os.Getenv("GOREST_SERVER_PROXY_HEADER"),
		GETOnly:                      getOnly,
		DisableKeepalive:             disableKeepalive,
		DisableDefaultDate:           disableDefaultDate,
		DisableDefaultContentType:    disableDefaultContentType,
		DisableHeaderNormalizing:     disableHeaderNormalizing,
		AppName:                      str("GOREST_APP_NAME", gorest.DefaultAppName),
		SharedStatePrefix:            os.Getenv("GOREST_SERVER_SHARED_STATE_PREFIX"),
		StreamRequestBody:            streamRequestBody,
		DisablePreParseMultipartForm: disablePreParseMultipartForm,
		ReduceMemoryUsage:            reduceMemoryUsage,
		TrustProxy:                   trustProxy,
		TrustProxyConfig: fiber.TrustProxyConfig{
			Proxies:    toStrArray("GOREST_SERVER_TRUST_PROXY_PROXIES", make([]string, 0)),
			LinkLocal:  linkLocal,
			Loopback:   loopback,
			Private:    private,
			UnixSocket: unixSocket,
		},
		EnableIPValidation:       enableIPValidation,
		RequestMethods:           toStrArray("GOREST_SERVER_REQUEST_METHODS", make([]string, 0)),
		EnableSplittingOnParsers: enableSplittingOnParsers,
	}, nil
}
