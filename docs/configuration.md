# Configuration, security, and deployment

Runtime configuration is defined under `app/config/` and loaded by `main.go`. The starter reads environment variables directly; it does not load `.env` files. For local development, unset `GOREST_APP_ENV` (or set it to `development`) to use the example configuration. Supply deployment values through the shell, container/orchestrator, or a secret manager.

`GOREST_APP_ENV` accepts exactly three values from GOREST: `development`, `production`, and `testing`. If it is unset, the application uses `development`. Values such as `test` and `staging` are rejected during configuration loading.

## GOREST release notifications

This starter enables GOREST's `CheckForUpdates` option in `main.go`. At startup, GOREST asynchronously checks the official GitHub latest-release endpoint and logs a notice if a newer GOREST release is available. It does not update dependencies or modify the application, and the check does not block or fail startup. The request has a three-second timeout. Enabling this feature means the running application makes an outbound request to GitHub at each startup; account for that in environments with restricted egress. To disable the check, change the application configuration to `CheckForUpdates: gorest.DisableAutoUpdate`. There is no environment-variable override in this starter.

## HTTP server

`app/config/server.go` sets:

- Local host `127.0.0.1` and port `3000`.
- A configurable request body limit and read, write, and idle timeouts (`GOREST_SERVER_BODY_LIMIT`, `GOREST_SERVER_READ_TIMEOUT`, `GOREST_SERVER_WRITE_TIMEOUT`, and `GOREST_SERVER_IDLE_TIMEOUT`). Durations use Go duration syntax, such as `10s` or `2m`; blank values use the local defaults.
- CORS origins from `GOREST_SERVER_CORS_ORIGINS` (comma-separated); local defaults are restricted to localhost.
- Recover, request ID, compression, helmet, and rate limiter middleware.
- GOREST's built-in operational endpoints are registered by default through `monitoring.Run`. The framework default allowlist is empty, but this starter defaults to loopback-only access (`127.0.0.1` and `::1`). Configure `GOREST_MONITORING_ALLOWED_IPS` to allow additional trusted health-check clients, or set `GOREST_MONITORING_DISABLE=true` to disable the routes entirely. An explicitly empty `GOREST_MONITORING_ALLOWED_IPS` denies all callers. The allowlist checks the TCP peer IP (the reverse proxy IP when proxied), not forwarded headers; it is not user authentication. See [operational monitoring endpoints](#operational-monitoring-endpoints).

Set `GOREST_SERVER_HOST` and `GOREST_SERVER_PORT` for deployments. Both are required when `GOREST_APP_ENV=production`, preventing accidental use of the loopback/local bind defaults. Set CORS to the exact browser origins you trust. CORS is a browser policy, not authentication. The default rate limiter uses in-memory storage, so its limit applies to each process rather than across all instances. If you run the application in a container, bind to an address reachable by the proxy and restrict access at the network boundary.

### Listener configuration

The starter maps listener environment variables to GOREST v1.1.17's `server.Config.Listener`. Blank values use GOREST defaults.

| Variable | Purpose | Default when blank |
|---|---|---|
| `GOREST_SERVER_LISTENER_NETWORK` | Listener network: `tcp`, `tcp4`, `tcp6`, or `unix` | `tcp4` |
| `GOREST_SERVER_UNIX_SOCKET_PATH` | Unix socket path; required for `unix` | Not applicable |
| `GOREST_SERVER_CERT_FILE` and `GOREST_SERVER_CERT_KEY_FILE` | Server TLS certificate and private key; configure both to enable HTTPS | TLS disabled |
| `GOREST_SERVER_CERT_CLIENT_FILE` | Client CA bundle for mutual TLS; requires server TLS | Client certificates not required |
| `GOREST_SERVER_TLS_MIN_VERSION` | Minimum TLS version: `1.2` or `1.3` | TLS 1.2 |
| `GOREST_SERVER_SHUTDOWN_TIMEOUT` | Grace period for HTTP server shutdown | 30s |
| `GOREST_SERVER_UNIX_SOCKET_FILE_MODE` | Unix socket permission mode in octal, e.g. `0770` | `0770` |
| `GOREST_SERVER_PREFORK_RECOVER_INTERVAL` | Delay before restarting a crashed prefork worker | `0s` (immediate restart) |
| `GOREST_SERVER_PREFORK_SHUTDOWN_GRACE_PERIOD` | Time to wait for prefork workers before force-stopping them | 5s |
| `GOREST_SERVER_PREFORK_RECOVER_THRESHOLD` | Worker restart threshold | GOREST default (at least 1) |
| `GOREST_SERVER_ENABLE_PREFORK` | Run GOREST-managed worker processes | Disabled |
| `GOREST_SERVER_ENABLE_PRINT_ROUTES` | Include route list in startup output | Disabled |

`GOREST_SERVER_HOST` and `GOREST_SERVER_PORT` choose the address for TCP listeners. For Unix sockets, set the network to `unix` and provide `GOREST_SERVER_UNIX_SOCKET_PATH`. TLS certificate and key must be provided together. Durations use Go syntax such as `30s` or `2m`; negative values are invalid. The socket mode must be an octal permission from `0000` through `0777`.

GOREST-managed prefork is supported on macOS and Linux and only with `tcp4` or `tcp6`. Workers initialize the application independently, so in-memory state—including rate limits and HTTP metrics—is not shared. Set the worker shutdown grace period in relation to the server shutdown timeout and the longest request that must finish. Enable prefork only after validating the application's dependencies and deployment process model.

Example HTTPS listener settings:

```sh
export GOREST_SERVER_CERT_FILE=/etc/example/server.crt
export GOREST_SERVER_CERT_KEY_FILE=/etc/example/server.key
export GOREST_SERVER_TLS_MIN_VERSION=1.2
export GOREST_SERVER_SHUTDOWN_TIMEOUT=30s
```

TLS can instead terminate at a trusted reverse proxy. In that case, secure the proxy-to-application network and configure proxy trust only when the application requires client addresses from forwarding headers.

### Operational monitoring endpoints

GOREST registers these built-in routes under `/gorest` by default:

- `/live` and `/ready` for liveness and readiness checks.
- `/metrics` for Go runtime metrics.
- `/http` for application HTTP metrics.
- `/storage` for estimated storage used by the configured database.
- `/database` for database and connection-pool metrics.

The framework's default `MonitoringAllowedIPs` is empty, which denies access. This starter overrides that default to allow only IPv4 and IPv6 loopback addresses when `GOREST_MONITORING_ALLOWED_IPS` is unset. An explicitly empty environment variable restores deny-all behavior. The allowlist is not authentication. GOREST checks the TCP connection's peer IP; behind a reverse proxy, allow only the proxy's trusted address and restrict access at the proxy/network layer as well. Do not allow public caller ranges.

For deployments with health checks running on a trusted non-loopback address, add the exact IP or CIDR:

```sh
export GOREST_MONITORING_ALLOWED_IPS=127.0.0.1/32,::1/128
```

Set `GOREST_MONITORING_DISABLE=true` to disable route registration completely. Set `GOREST_MONITORING_ENABLE_HTTP_METRICS=true` only when route-level breakdowns from `/gorest/http` are needed; it adds per-request processing and memory use. `/gorest/storage` is an estimate of database storage, not free filesystem capacity or maximum disk size. For filesystem capacity, use operating-system or infrastructure monitoring. See the [GOREST v1.1.17 configuration guide](https://github.com/elsyahtech/gorest/blob/v1.1.17/docs/configuration.md#operational-endpoints-and-environment-configuration).

## Database and Redis

- SQLite is the local default through `database.ConfigDefault`. The starter points it to `./app/database` and enables migrations and sample seeders for local examples. No external service or database credentials are needed initially.
- Set `GOREST_DB_DRIVER`, `GOREST_DB_HOST`, `GOREST_DB_PORT`, `GOREST_DB_NAME`, `GOREST_DB_USERNAME`, and `GOREST_DB_PASSWORD` to select and connect to a server database. Use `GOREST_DB_TLS`, `GOREST_DB_SSL_CA_PATH`, `GOREST_DB_SSL_CERT_PATH`, `GOREST_DB_SSL_KEY_PATH`, and `GOREST_DB_LIB_DIR` where required by the driver. `GOREST_DB_MIGRATION`, `GOREST_DB_SEEDER`, and `GOREST_DB_DIRECTORY` control startup tasks and file locations. See the [database guide](./database.md) for details.
- Redis is disabled by default. Set `GOREST_REDIS_DISABLE=false` to enable Redis for the Redis query example routes; connection and secret settings use the `GOREST_REDIS_*` variables described below.
- When `GOREST_APP_ENV=production`, database driver and connection fields (or an explicit SQLite directory) must be supplied, migration and seeder flags must be false, and the local server host and port must be overridden. SQLite uses a local file; GOREST rejects non-local SQLite hosts in production and recommends a server database for remote or horizontally scaled workloads. The starter requires an explicit `GOREST_DB_TLS` mode for network drivers except ScyllaDB, but does not universally require that mode to enable encryption. GOREST warns when production database transport encryption is disabled or not guaranteed; the warning does not block startup. See the [database guide](./database.md#production-transport-encryption-by-driver). Invalid boolean and port values fail before services start. The application does not read `.env` files; inject secrets through a protected environment or secret manager, and do not commit them to source.

## Authentication example

The starter demonstrates token authentication on `GET /examples/orm-queries/users/with-authentication`; `POST /examples/login` issues its `access` cookie. Other `/examples` routes are not protected, and this example does not implement authorization.

`GOREST_FILTER_SIGNATURE_KEY` configures the key used to sign authentication tokens. If it is unset, the current starter configuration falls back to the literal `your-signature-key`. That value is a development placeholder, not a secret; the application does not reject it automatically. Set a unique, high-entropy value through a protected environment or secret manager before using authentication in any deployed environment. Do not commit the value to source control.

### Secure database connections and `AllowInsecureProduction`

For production, configure database connections to use TLS and verify the database server's identity. Encryption without certificate/hostname verification may not protect against connecting to an impostor server. TLS options differ by database driver; follow the [per-driver guidance](./database.md#production-transport-encryption-by-driver) and use the strongest mode supported by your database and deployment. Do not set `GOREST_DB_TLS=disable` for production network connections.

GOREST's `AllowInsecureProduction` option is an explicit exception to its TLS safeguard: it allows a connection even when GOREST considers the environment production-like and TLS is not enabled. Despite the word `Production` in its name, this is **not** a recommendation to use insecure connections in production. It does **not** enable TLS or encrypt traffic. The starter intentionally rejects `GOREST_DB_ALLOW_INSECURE_PRODUCTION=true` when `GOREST_APP_ENV=production`, so developers cannot use this option to bypass the production TLS policy.

#### Local tests using a database without TLS

Use GOREST's standard testing environment for tests: set `GOREST_APP_ENV=testing` when loading the application's config and use `gorest.Testing` for `gorest.Config.Environment`. Do not use the literal string `"Test"`; it is not the standard environment value and GOREST may treat an unrecognized environment as production-like. With the standard testing environment, tests do not need `AllowInsecureProduction` just because a local test database has TLS disabled.

Only if a test intentionally uses a production-like/custom environment with a trusted local database that has no TLS should you explicitly allow that non-TLS connection. Keep the exception in test configuration; when the database supports TLS, configure TLS instead and leave this option at its default (`false`).

#### If the application loads settings with `LoadDatabase`

If your test deliberately uses a production-like/custom GOREST environment, set this variable in the same shell as your test. Add the database name and credentials required by your test:

```bash
export GOREST_APP_ENV=testing
export GOREST_DB_DRIVER=mysql
export GOREST_DB_HOST=127.0.0.1
export GOREST_DB_PORT=3306
export GOREST_DB_TLS=disable
export GOREST_DB_ALLOW_INSECURE_PRODUCTION=true
```

Then run the test command for your project. These environment variables affect only code that loads its database configuration through `LoadDatabase`. If the test database supports TLS, configure TLS and leave `GOREST_DB_ALLOW_INSECURE_PRODUCTION` unset (its default is `false`).

#### If a test creates `database.Config` directly

Environment variables are only read by `LoadDatabase`. If a test builds the GOREST config itself, set the field on that test config instead:

```go
dbConfig := database.ConfigDefault
dbConfig.Driver = database.MYSQL
dbConfig.TLS = "disable"
dbConfig.AllowInsecureProduction = true
```

Invalid boolean values cause configuration loading to fail. See the [testing guide](./testing.md#non-tls-connections-in-local-database-tests) for how this project uses the test-only exception, and the [GOREST v1.1.17 configuration guide](https://github.com/elsyahtech/gorest/blob/v1.1.17/docs/configuration.md) for framework details.

## Logging

The logger writes to `./logs/app.log` and has rotation settings in `app/config/logger.go`. Set the log location, file permissions, retention, and storage limits for your deployment. Do not include passwords, tokens, or other sensitive data in logs.

## Environment variable reference

| Variable | Purpose | Local default |
| :-- | :-- | :-- |
| `GOREST_APP_ENV` | GOREST runtime environment (`development`, `production`, or `testing`) | `development` |
| `GOREST_SERVER_HOST`, `GOREST_SERVER_PORT` | HTTP bind address and port | `127.0.0.1`, `3000` |
| `GOREST_SERVER_LISTENER_NETWORK`, `GOREST_SERVER_UNIX_SOCKET_PATH` | TCP or Unix socket listener; socket path is required for `unix` | `tcp4`; unset |
| `GOREST_SERVER_CERT_FILE`, `GOREST_SERVER_CERT_KEY_FILE`, `GOREST_SERVER_CERT_CLIENT_FILE`, `GOREST_SERVER_TLS_MIN_VERSION` | Server TLS and optional mutual TLS settings | TLS off; TLS 1.2 minimum |
| `GOREST_SERVER_SHUTDOWN_TIMEOUT`, `GOREST_SERVER_PREFORK_SHUTDOWN_GRACE_PERIOD`, `GOREST_SERVER_PREFORK_RECOVER_INTERVAL`, `GOREST_SERVER_PREFORK_RECOVER_THRESHOLD` | Graceful shutdown and prefork worker recovery | `30s`, `5s`, immediate, GOREST default |
| `GOREST_SERVER_UNIX_SOCKET_FILE_MODE`, `GOREST_SERVER_ENABLE_PREFORK`, `GOREST_SERVER_ENABLE_PRINT_ROUTES` | Unix socket permissions, prefork, and route-list output | `0770`, `false`, `false` |
| `GOREST_SERVER_BODY_LIMIT` | Maximum request body size in bytes | `4194304` |
| `GOREST_SERVER_READ_TIMEOUT`, `GOREST_SERVER_WRITE_TIMEOUT`, `GOREST_SERVER_IDLE_TIMEOUT` | HTTP read, write, and idle timeouts; Go duration syntax | `10s`, `30s`, `120s` |
| `GOREST_SERVER_CORS_ORIGINS`, `GOREST_SERVER_CORS_METHODS`, `GOREST_SERVER_CORS_HEADERS`, `GOREST_SERVER_CORS_EXPOSE_HEADERS` | Comma-separated CORS origins, methods, allowed request headers, and exposed response headers | Local origins; `GET,HEAD,POST,PUT,PATCH,DELETE,OPTIONS`; `Origin,Content-Type,Accept,Authorization,X-Request-ID`; empty |
| `GOREST_SERVER_CORS_MAX_AGE`, `GOREST_SERVER_CORS_ALLOW_CREDENTIALS`, `GOREST_SERVER_CORS_ALLOW_PRIVATE_NETWORK`, `GOREST_SERVER_CORS_DISABLE_VALUE_REDACTION` | CORS preflight cache age in seconds and CORS behavior flags | `300`; `false` |
| `GOREST_SERVER_RECOVER_ENABLE_STACK_TRACE` | Include stack traces in recovered panic handling | `false` |
| `GOREST_SERVER_REQUEST_ID_HEADER` | Request ID header name | `X-Request-ID` |
| `GOREST_SERVER_COMPRESS_LEVEL` | Compression level: `-1` disabled, `0` default, `1` best speed, `2` best compression | `0` |
| `GOREST_SERVER_HELMET_*` | Helmet security-header values, HSTS max age, and CSP/HSTS flags; see `.env.example` for the individual field names | GOREST middleware defaults |
| `GOREST_SERVER_RATE_LIMIT_MAX`, `GOREST_SERVER_RATE_LIMIT_EXPIRATION` | Maximum requests per in-memory rate-limit window and its duration | `50`, `1m` |
| `GOREST_SERVER_RATE_LIMIT_SKIP_FAILED_REQUESTS`, `GOREST_SERVER_RATE_LIMIT_SKIP_SUCCESSFUL_REQUESTS`, `GOREST_SERVER_RATE_LIMIT_DISABLE_HEADERS`, `GOREST_SERVER_RATE_LIMIT_DISABLE_VALUE_REDACTION` | Rate-limiter behavior flags | `false` |
| `GOREST_MONITORING_ALLOWED_IPS` | Comma-separated IP addresses/CIDRs allowed to access built-in `/gorest` endpoints; an explicitly empty value denies all | `127.0.0.1,::1` |
| `GOREST_MONITORING_DISABLE` | Disable registration of built-in monitoring endpoints | `false` (routes are registered; loopback callers are allowed) |
| `GOREST_MONITORING_ENABLE_HTTP_METRICS` | Include route-level breakdown in `/gorest/http` metrics | `false` |
| `GOREST_DB_DRIVER`, `GOREST_DB_HOST`, `GOREST_DB_PORT`, `GOREST_DB_NAME` | Database selection and endpoint | SQLite, GOREST defaults |
| `GOREST_DB_USERNAME`, `GOREST_DB_PASSWORD` | Database credentials | Empty |
| `GOREST_DB_TLS` | Optional driver-specific TLS mode; required explicitly by this starter for production MySQL, PostgreSQL, Oracle, SQL Server, and MongoDB (the mode may disable TLS) | GOREST default; ScyllaDB does not use this setting; SQLite has no network TLS |
| `GOREST_DB_ALLOW_INSECURE_PRODUCTION` | Explicitly bypass GOREST's TLS safeguard; this starter rejects `true` in production and documents it only for local tests | `false` |
| `GOREST_FILTER_SIGNATURE_KEY` | Secret used to sign authentication tokens | `your-signature-key` (development placeholder; unsafe for deployment) |
| `GOREST_DB_SSL_CA_PATH`, `GOREST_DB_SSL_CERT_PATH`, `GOREST_DB_SSL_KEY_PATH`, `GOREST_DB_LIB_DIR` | Driver-specific CA/client certificate, wallet, or Oracle client paths | Empty |
| `GOREST_DB_DIRECTORY` | SQLite file and migration/seeder root | `./app/database` |
| `GOREST_DB_MIGRATION`, `GOREST_DB_SEEDER` | Run migrations/seeders during application startup | `true` / `true` in development; `false` / `false` in production |
| `GOREST_REDIS_DISABLE` | Skip Redis service initialization | `true` (Redis disabled); set to `false` to enable |
| `GOREST_REDIS_HOST`, `GOREST_REDIS_PORT`, `GOREST_REDIS_DATABASE`, `GOREST_REDIS_TLS` | Redis connection settings | `127.0.0.1:6379`, DB `0`, TLS off |
| `GOREST_REDIS_USERNAME`, `GOREST_REDIS_PASSWORD` | Optional Redis credentials | Empty |
| `GOREST_REDIS_SSL_CA_PATH`, `GOREST_REDIS_SSL_CERT_PATH`, `GOREST_REDIS_SSL_KEY_PATH` | Optional Redis TLS files | Empty |

`.env.example` lists the scalar and list settings exposed by this starter's server configuration. Middleware callbacks, custom storage, and custom handlers remain code-level configuration rather than environment variables.

## Before deployment

- Configure database credentials for the target environment, and ensure real production credentials are not committed to the repository.
- Most `/examples` routes are unauthenticated. The login and one protected-route example do not secure the remaining routes or provide production authorization; add appropriate authentication and authorization before exposing any business API.
- Replace the default `GOREST_FILTER_SIGNATURE_KEY` placeholder with a unique, high-entropy secret before deploying.
- `Migration=true` / `Seeder=true` are local-example defaults only. Under `GOREST_APP_ENV=production`, they default to false, and startup rejects either flag if explicitly enabled. Run schema migrations once in a controlled deployment step; do not run sample seeders against production data.
- Restrict CORS to the origins you need. CORS does not prevent non-browser clients from calling the API.
- Choose and verify transport encryption for each database and network path according to your deployment policy. GOREST's production warning is informational and does not block startup; see the [per-driver TLS explanation](./database.md#production-transport-encryption-by-driver).
- Review the bind address, proxy trust, request limits, rate limiting across instances, logging, health checks, and backup and restore policies.
- Replace or remove demonstration routes that your business application does not need.

See the [database guide](./database.md) for migrations and the [API reference](./api-reference.md) for the endpoints that need protection.
