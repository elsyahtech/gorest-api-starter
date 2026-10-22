package config

import (
	"time"

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
func Server() server.Config {
	return server.Config{
		// Keep request timeouts finite and cap the body size. These timeout values
		// match GOREST defaults; adjust them for the application's workload.
		Core: &fiber.Config{
			BodyLimit:    4 * 1024 * 1024, // 4 MiB
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		},

		// Replace this example with the exact browser origins allowed to call the API.
		// CORS is enforced by browsers; it does not authenticate or authorize requests.
		Cors: &cors.Config{
			AllowOrigins: []string{"https://app.example.com"},
			AllowMethods: []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
			MaxAge:       300,
		},

		// Recover converts handler panics into HTTP errors. It does not replace
		// input validation or application error handling.
		Recover: &recover.ConfigDefault,

		// Attach a request ID to make logs and request troubleshooting easier.
		Requestid: &requestid.ConfigDefault,

		// Compress eligible responses to reduce response size.
		Compress: &compress.ConfigDefault,

		// Add common HTTP security headers. These headers do not replace
		// authentication, authorization, or secure deployment configuration.
		Helmet: &helmet.ConfigDefault,

		// The default limiter uses in-memory storage, so limits apply per process.
		// Use shared storage or an upstream gateway if limits must span instances.
		// Review proxy trust settings before relying on client IP behind a proxy.
		Limiter: &limiter.Config{
			Max:        50,
			Expiration: time.Minute,
		},

		// Local-only by default. In a container, bind to 0.0.0.0 only when the
		// deployment network/firewall controls which clients can reach the service.
		Host: "127.0.0.1",
		Port: 3000,
	}
}
