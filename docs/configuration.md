# Configuration, security, and deployment

The starter's runtime configuration is defined in Go under `app/config/`, with some application settings in `main.go`. The application does not currently load these settings automatically from environment variables or a `.env` file. Set up an appropriate configuration method before using the application across multiple environments.

## HTTP server

`app/config/server.go` sets:

- Host `127.0.0.1` and port `3000`.
- A 4 MiB body limit and read, write, and idle timeouts.
- The example CORS origin `https://app.example.com`, allowed methods and headers, and max age.
- Recover, request ID, compression, helmet, and rate limiter middleware.

Change CORS to the exact origins you trust. CORS is a browser policy, not authentication. The default rate limiter uses in-memory storage, so its limit applies to each process rather than across all instances. If you run the application in a container, change the bind host only when your network and firewall restrict access.

## Database and Redis

- PostgreSQL is enabled by default. See [`database.go`](../app/config/database.go) for TLS, credentials, automatic migrations and seeders, and connection pool settings.
- Redis is registered as a service but has `Enabled` set to `false`. If you enable it, configure a reachable host, authentication, and TLS as required by your network.
- Do not use the sample credentials in source code. Use a secrets manager or another suitable deployment secret store, and do not write secrets to logs.

## Logging

The logger writes to `./logs/app.log` and has rotation settings in `app/config/logger.go`. Set the log location, file permissions, retention, and storage limits for your deployment. Do not include passwords, tokens, or other sensitive data in logs.

## Before deployment

- Replace the sample credentials and store secrets in your deployment's secret management system.
- Add authentication and authorization before exposing routes. None of the `/examples` endpoints currently has access controls.
- Do not run sample seeders or automatic migrations on every production replica. Run migrations in a controlled step.
- Restrict CORS to the origins you need. CORS does not prevent non-browser clients from calling the API.
- Configure TLS at a trusted ingress or reverse proxy, and configure database TLS and certificate verification.
- Review the bind address, proxy trust, request limits, rate limiting across instances, logging, health checks, and backup and restore policies.
- Replace or remove demonstration routes that your business application does not need.

See the [database guide](./database.md) for migrations and the [API reference](./api-reference.md) for the endpoints that need protection.
