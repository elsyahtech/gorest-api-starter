package config

import (
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

// Server configures the HTTP server and its middleware.
// These middleware provide request handling controls; applications must still
// implement authentication, authorization, and any required CSRF protections.
// Configure TLS at a trusted ingress/reverse proxy or in the server deployment.
func ServerConfig() (server.Config, error) {
	cfg := server.ConfigDefault

	cfg.Host = str(os.Getenv("GOREST_SERVER_HOST"), server.DefaultServerHost)
	cfg.Port = strToInt(os.Getenv("GOREST_SERVER_PORT"), server.DefaultServerPort)
	cfg.Recover = &recover.Config{
		EnableStackTrace: strToBool(os.Getenv("GOREST_SERVER_RECOVER_ENABLE_STACK_TRACE"), false),
	}

	cfg.Requestid = &requestid.Config{
		Header: os.Getenv("GOREST_SERVER_REQUEST_ID_HEADER"),
	}

	cfg.Compress = &compress.Config{
		Level: compress.Level(strToInt(os.Getenv("GOREST_SERVER_COMPRESS_LEVEL"), int(compress.LevelDefault))),
	}

	cfg.Core = fiberCoreConfig()

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

	cfg.Cors = &cors.Config{
		AllowOrigins:          origins,
		AllowMethods:          methods,
		AllowHeaders:          headers,
		ExposeHeaders:         toStrArray("GOREST_SERVER_CORS_EXPOSE_HEADERS", []string{}),
		MaxAge:                strToInt(os.Getenv("GOREST_SERVER_CORS_MAX_AGE"), 300),
		AllowCredentials:      strToBool(os.Getenv("GOREST_SERVER_CORS_ALLOW_CREDENTIALS"), false),
		AllowPrivateNetwork:   strToBool(os.Getenv("GOREST_SERVER_CORS_ALLOW_PRIVATE_NETWORK"), false),
		DisableValueRedaction: strToBool(os.Getenv("GOREST_SERVER_CORS_DISABLE_VALUE_REDACTION"), false),
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
		HSTSMaxAge:                strToInt(os.Getenv("GOREST_SERVER_HELMET_HSTS_MAX_AGE"), 0),
		HSTSExcludeSubdomains:     strToBool(os.Getenv("GOREST_SERVER_HELMET_HSTS_EXCLUDE_SUBDOMAINS"), false),
		CSPReportOnly:             strToBool(os.Getenv("GOREST_SERVER_HELMET_CSP_REPORT_ONLY"), false),
		HSTSPreloadEnabled:        strToBool(os.Getenv("GOREST_SERVER_HELMET_HSTS_PRELOAD_ENABLED"), false),
	}

	cfg.Limiter = &limiter.Config{
		Max:                    strToInt(os.Getenv("GOREST_SERVER_RATE_LIMIT_MAX"), 50),
		Expiration:             strToTimeDuration(os.Getenv("GOREST_SERVER_RATE_LIMIT_EXPIRATION"), time.Minute),
		SkipFailedRequests:     strToBool(os.Getenv("GOREST_SERVER_RATE_LIMIT_SKIP_FAILED_REQUESTS"), false),
		SkipSuccessfulRequests: strToBool(os.Getenv("GOREST_SERVER_RATE_LIMIT_SKIP_SUCCESSFUL_REQUESTS"), false),
		DisableHeaders:         strToBool(os.Getenv("GOREST_SERVER_RATE_LIMIT_DISABLE_HEADERS"), false),
		DisableValueRedaction:  strToBool(os.Getenv("GOREST_SERVER_RATE_LIMIT_DISABLE_VALUE_REDACTION"), false),
	}

	listener, err := loadServerListenerConfig()
	if err != nil {
		return server.Config{}, err
	}
	cfg.Listener = listener

	return cfg, nil
}

func fiberCoreConfig() *fiber.Config {
	return &fiber.Config{
		ServerHeader:                 os.Getenv("GOREST_SERVER_HEADER"),
		StrictRouting:                strToBool(os.Getenv("GOREST_SERVER_STRICT_ROUTING"), false),
		CaseSensitive:                strToBool(os.Getenv("GOREST_SERVER_CASE_SENSITIVE"), false),
		SkipUnmatchedRoutes:          strToBool(os.Getenv("GOREST_SERVER_SKIP_UNMATCHED_ROUTES"), false),
		DisableHeadAutoRegister:      strToBool(os.Getenv("GOREST_SERVER_DISABLE_HEAD_AUTO_REGISTER"), false),
		DisableFlashMessages:         strToBool(os.Getenv("GOREST_SERVER_DISABLE_FLASH_MESSAGES"), false),
		Immutable:                    strToBool(os.Getenv("GOREST_SERVER_IMMUTABLE"), false),
		UnescapePath:                 strToBool(os.Getenv("GOREST_SERVER_UNESCAPE_PATH"), false),
		BodyLimit:                    strToInt(os.Getenv("GOREST_SERVER_BODY_LIMIT"), server.DefaultServerBodyLimit),
		MaxRanges:                    strToInt(os.Getenv("GOREST_SERVER_MAX_RANGES"), fiber.DefaultMaxRanges),
		Concurrency:                  strToInt(os.Getenv("GOREST_SERVER_CONCURRENCY"), fiber.DefaultConcurrency),
		ViewsLayout:                  os.Getenv("GOREST_SERVER_VIEWS_LAYOUT"),
		PassLocalsToViews:            strToBool(os.Getenv("GOREST_SERVER_PASS_LOCALS_TO_VIEWS"), false),
		PassLocalsToContext:          strToBool(os.Getenv("GOREST_SERVER_PASS_LOCALS_TO_CONTEXT"), false),
		ReadTimeout:                  strToTimeDuration(os.Getenv("GOREST_SERVER_READ_TIMEOUT"), server.DefaultServerReadTimeout),
		WriteTimeout:                 strToTimeDuration(os.Getenv("GOREST_SERVER_WRITE_TIMEOUT"), server.DefaultServerWriteTimeout),
		IdleTimeout:                  strToTimeDuration(os.Getenv("GOREST_SERVER_IDLE_TIMEOUT"), server.DefaultServerIdleTimeout),
		ReadBufferSize:               strToInt(os.Getenv("GOREST_SERVER_READ_BUFFER_SIZE"), fiber.DefaultReadBufferSize),
		WriteBufferSize:              strToInt(os.Getenv("GOREST_SERVER_WRITE_BUFFER_SIZE"), fiber.DefaultWriteBufferSize),
		ProxyHeader:                  os.Getenv("GOREST_SERVER_PROXY_HEADER"),
		GETOnly:                      strToBool(os.Getenv("GOREST_SERVER_GET_ONLY"), false),
		DisableKeepalive:             strToBool(os.Getenv("GOREST_SERVER_DISABLE_KEEPALIVE"), false),
		DisableDefaultDate:           strToBool(os.Getenv("GOREST_SERVER_DISABLE_DEFAULT_DATE"), false),
		DisableDefaultContentType:    strToBool(os.Getenv("GOREST_SERVER_DISABLE_DEFAULT_CONTENT_TYPE"), false),
		DisableHeaderNormalizing:     strToBool(os.Getenv("GOREST_SERVER_DISABLE_HEADER_NORMALIZING"), false),
		AppName:                      str(os.Getenv("GOREST_APP_NAME"), gorest.DefaultAppName),
		SharedStatePrefix:            os.Getenv("GOREST_SERVER_SHARED_STATE_PREFIX"),
		StreamRequestBody:            strToBool(os.Getenv("GOREST_SERVER_STREAM_REQUEST_BODY"), false),
		DisablePreParseMultipartForm: strToBool(os.Getenv("GOREST_SERVER_DISABLE_PRE_PARSE_MULTIPART_FORM"), false),
		ReduceMemoryUsage:            strToBool(os.Getenv("GOREST_SERVER_REDUCE_MEMORY_USAGE"), false),
		TrustProxy:                   strToBool(os.Getenv("GOREST_SERVER_TRUST_PROXY"), false),
		TrustProxyConfig: fiber.TrustProxyConfig{
			Proxies:   toStrArray("GOREST_SERVER_TRUST_PROXY_PROXIES", []string{}),
			LinkLocal: strToBool(os.Getenv("GOREST_SERVER_TRUST_PROXY_LINK_LOCAL"), false),
			Loopback:  strToBool(os.Getenv("GOREST_SERVER_TRUST_PROXY_LOOPBACK"), false),
			Private:   strToBool(os.Getenv("GOREST_SERVER_TRUST_PROXY_PRIVATE"), false),
			UnixSocket: strToBool(
				os.Getenv("GOREST_SERVER_TRUST_PROXY_UNIX_SOCKET"),
				false,
			),
		},
		EnableIPValidation:       strToBool(os.Getenv("GOREST_SERVER_ENABLE_IP_VALIDATION"), false),
		RequestMethods:           toStrArray("GOREST_SERVER_REQUEST_METHODS", []string{}),
		EnableSplittingOnParsers: strToBool(os.Getenv("GOREST_SERVER_ENABLE_SPLITTING_ON_PARSERS"), false),
	}
}
