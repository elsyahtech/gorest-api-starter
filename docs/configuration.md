# Configuration and deployment

The active runtime configuration is defined in `app/config/`. The executable loads `.env` using `godotenv` when the file is present; environment-only startup is supported, and process environment variables take precedence. Keep `.env`, test credentials, private keys, and certificate files out of source control.

## Application and HTTP server

The sample uses GOREST v1.1.20. Application settings use `GOREST_APP_*`; the HTTP listener and middleware use the `GOREST_SERVER_*` variables in `.env.example`. The default local bind is `127.0.0.1:3000`. Set the bind address intentionally when deploying behind a proxy or container network.

Server listener options include TCP/Unix network, TLS certificate files, TLS minimum version, shutdown timeout, Unix socket mode, and prefork settings. Blank listener values use the defaults provided by the GOREST version selected in `go.mod`. See `.env.example` and the [GOREST v1.1.20 configuration guide](https://github.com/elsyahtech/gorest/blob/v1.1.20/docs/configuration.md).

## MySQL

MySQL connection variables are documented in the [MySQL configuration guide](./database.md). GOREST v1.1.20 supports verified server TLS, custom CA files, and client certificates for mTLS through the database configuration. Follow the guide to set the required CA and client key pair.

## Redis (optional)

Redis is optional, disabled by default by both `app/config/redis.go` and `.env.example`. The current user endpoints do not use Redis, so the starter runs without a Redis server. Enable Redis only when your application needs a Redis-backed feature.

### Connect to Redis

1. Provision a Redis instance reachable from the application and create a dedicated ACL user where authentication is required. Keep passwords in environment variables or a secret manager.
2. Set `GOREST_REDIS_DISABLE=false`, `GOREST_REDIS_HOST`, and `GOREST_REDIS_PORT`. Configure `GOREST_REDIS_DATABASE` if needed and `GOREST_REDIS_USER` / `GOREST_REDIS_PASSWORD` when the server requires authentication.
3. If Redis supports TLS and its certificate chains to a CA trusted by the operating system, set `GOREST_REDIS_TLS=true`. GOREST verifies the server certificate and hostname; no CA file or client certificate is needed for server-authenticated TLS.
4. Start the application. GOREST connects and sends `PING` during service initialization, so an unreachable server, invalid credentials, or failed certificate verification prevents startup.

### Use a private CA or mutual TLS

If the Redis server certificate is signed by a private CA, set `GOREST_REDIS_CA_PATH` to a readable PEM file containing that CA certificate. Include any intermediate certificates needed to verify the server.

If the Redis server requires mutual TLS (mTLS), set `GOREST_REDIS_TLS=true` and provide both `GOREST_REDIS_CERT_PATH` and `GOREST_REDIS_KEY_PATH` for the client certificate and its private key. The server must trust the client certificate's issuer and require client authentication. For a private server CA, also set `GOREST_REDIS_CA_PATH`; the client certificate/key authenticate the client, while the CA path verifies the Redis server. The starter rejects a client certificate or key configured without its matching pair.

For example, with a private CA and an mTLS-enabled Redis server:

```dotenv
GOREST_REDIS_DISABLE=false
GOREST_REDIS_TLS=true
GOREST_REDIS_CA_PATH=/run/secrets/redis-ca.pem
GOREST_REDIS_CERT_PATH=/run/secrets/redis-client.pem
GOREST_REDIS_KEY_PATH=/run/secrets/redis-client-key.pem
```

Mount certificate files as secrets with read access limited to the application account. Do not commit credentials, private keys, or certificates. The starter passes these settings to GOREST using `redis.Run(redisConfig)`. See the [GOREST v1.1.20 configuration guide](https://github.com/elsyahtech/gorest/blob/v1.1.20/docs/configuration.md) for pool, timeout, retry, and TLS options. See [Testing and coverage](./testing.md) for optional connectivity checks.

## Authentication and sample limits

Replace the example `GOREST_FILTER_SIGNATURE_KEY` with a unique, high-entropy secret. Review cookie security settings and use HTTPS in deployments where credentials or cookies cross an untrusted network.

The create and login endpoints are public examples. Only `GET /users/` uses authentication; the example does not implement authorization roles or protect business resources. Add application-specific access control before exposing it.

## Monitoring and logs

Monitoring endpoints are managed by GOREST and configured through `GOREST_MONITORING_*`. The starter defaults to loopback access. Restrict monitoring at the network or proxy boundary as well. Choose a log directory and retention policy appropriate to the deployment, and avoid logging credentials or tokens.

For repeatable tests and coverage evidence, see [Testing and coverage](./testing.md).
