# Configuration, security, and deployment

Runtime configuration is defined under `app/config/` and loaded by `main.go`. The starter reads environment variables directly; it does not load `.env` files. For local development, unset `APP_ENV` (or set it to `development`) to use the SQLite example configuration. Supply deployment values through the shell, container/orchestrator, or a secret manager.

`APP_ENV` accepts exactly three values from GOREST: `development`, `production`, and `testing`. If it is unset, the application uses `development`. Values such as `test` and `staging` are rejected during configuration loading.

## HTTP server

`app/config/server.go` sets:

- Local host `127.0.0.1` and port `3000`.
- A 4 MiB body limit and read, write, and idle timeouts.
- CORS origins from `GOREST_CORS_ORIGINS` (comma-separated); local defaults are restricted to localhost.
- Recover, request ID, compression, helmet, and rate limiter middleware.

Set `GOREST_SERVER_HOST` and `GOREST_SERVER_PORT` for deployments. Both are required when `APP_ENV=production`, preventing accidental use of the loopback/local bind defaults. Set CORS to the exact browser origins you trust. CORS is a browser policy, not authentication. The default rate limiter uses in-memory storage, so its limit applies to each process rather than across all instances. If you run the application in a container, bind to an address reachable by the proxy and restrict access at the network boundary.

## Database and Redis

- SQLite is the local default through `database.ConfigDefault`. The starter points it to `./app/database` and enables migrations and sample seeders for local examples. No external service or database credentials are needed initially.
- Set `GOREST_DB_DRIVER`, `GOREST_DB_HOST`, `GOREST_DB_PORT`, `GOREST_DB_NAME`, `GOREST_DB_USERNAME`, and `GOREST_DB_PASSWORD` to select and connect to a server database. Use `GOREST_DB_TLS`, `GOREST_DB_SSL_CA_PATH`, `GOREST_DB_SSL_CERT_PATH`, `GOREST_DB_SSL_KEY_PATH`, and `GOREST_DB_LIB_DIR` where required by the driver. `GOREST_DB_MIGRATION`, `GOREST_DB_SEEDER`, and `GOREST_DB_DIRECTORY` control startup tasks and file locations. See the [database guide](./database.md) for details.
- Redis is disabled by default. Set `GOREST_REDIS_DISABLE=false` to enable Redis for the Redis query example routes; connection and secret settings use the `GOREST_REDIS_*` variables described below.
- When `APP_ENV=production`, database driver and connection fields (or an explicit SQLite directory) must be supplied, migration and seeder flags must be false, and the local server host and port must be overridden. The starter requires an explicit `GOREST_DB_TLS` mode for network drivers except ScyllaDB, but does not universally require that mode to enable encryption. GOREST warns when production database transport encryption is disabled or not guaranteed; the warning does not block startup. See the [per-driver TLS explanation](./database.md#production-transport-encryption-by-driver). Invalid boolean and port values fail before services start. The application does not read `.env` files; inject secrets through a protected environment or secret manager, and do not commit them to source.

## Logging

The logger writes to `./logs/app.log` and has rotation settings in `app/config/logger.go`. Set the log location, file permissions, retention, and storage limits for your deployment. Do not include passwords, tokens, or other sensitive data in logs.

## Environment variable reference

| Variable | Purpose | Local default |
| :-- | :-- | :-- |
| `APP_ENV` | GOREST runtime environment (`development`, `production`, or `testing`) | `development` |
| `GOREST_SERVER_HOST`, `GOREST_SERVER_PORT` | HTTP bind address and port | `127.0.0.1`, `3000` |
| `GOREST_CORS_ORIGINS` | Comma-separated allowed browser origins | `http://localhost:3000`, `http://127.0.0.1:3000` |
| `GOREST_DB_DRIVER`, `GOREST_DB_HOST`, `GOREST_DB_PORT`, `GOREST_DB_NAME` | Database selection and endpoint | SQLite, GOREST defaults |
| `GOREST_DB_USERNAME`, `GOREST_DB_PASSWORD` | Database credentials | Empty |
| `GOREST_DB_TLS` | Optional driver-specific TLS mode; required explicitly by this starter for production MySQL, PostgreSQL, Oracle, SQL Server, and MongoDB (the mode may disable TLS) | GOREST default; ScyllaDB does not use this setting; SQLite has no network TLS |
| `GOREST_DB_SSL_CA_PATH`, `GOREST_DB_SSL_CERT_PATH`, `GOREST_DB_SSL_KEY_PATH`, `GOREST_DB_LIB_DIR` | Driver-specific CA/client certificate, wallet, or Oracle client paths | Empty |
| `GOREST_DB_DIRECTORY` | SQLite file and migration/seeder root | `./app/database` |
| `GOREST_DB_MIGRATION`, `GOREST_DB_SEEDER` | Run migrations/seeders during application startup | `true` / `true` in development; `false` / `false` in production |
| `GOREST_REDIS_DISABLE` | Skip Redis service initialization | `true` (Redis disabled); set to `false` to enable |
| `GOREST_REDIS_HOST`, `GOREST_REDIS_PORT`, `GOREST_REDIS_DATABASE`, `GOREST_REDIS_TLS` | Redis connection settings | `127.0.0.1:6379`, DB `0`, TLS off |
| `GOREST_REDIS_USERNAME`, `GOREST_REDIS_PASSWORD` | Optional Redis credentials | Empty |
| `GOREST_REDIS_SSL_CA_PATH`, `GOREST_REDIS_SSL_CERT_PATH`, `GOREST_REDIS_SSL_KEY_PATH` | Optional Redis TLS files | Empty |

## Before deployment

- Configure database credentials for the target environment, and ensure real production credentials are not committed to the repository.
- Add authentication and authorization before exposing routes. None of the `/examples` endpoints currently has access controls.
- `Migration=true` / `Seeder=true` are local-example defaults only. Under `APP_ENV=production`, they default to false, and startup rejects either flag if explicitly enabled. Run schema migrations once in a controlled deployment step; do not run sample seeders against production data.
- Restrict CORS to the origins you need. CORS does not prevent non-browser clients from calling the API.
- Choose and verify transport encryption for each database and network path according to your deployment policy. GOREST's production warning is informational and does not block startup; see the [per-driver TLS explanation](./database.md#production-transport-encryption-by-driver).
- Review the bind address, proxy trust, request limits, rate limiting across instances, logging, health checks, and backup and restore policies.
- Replace or remove demonstration routes that your business application does not need.

See the [database guide](./database.md) for migrations and the [API reference](./api-reference.md) for the endpoints that need protection.
