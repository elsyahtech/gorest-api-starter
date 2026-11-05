# Configuration and deployment

The active runtime configuration is defined in `app/config/`. The executable loads `.env` using `godotenv` when the file is present; environment-only startup is supported, and process environment variables take precedence. Keep `.env`, test credentials, private keys, and certificate files out of source control.

## Application and HTTP server

The sample uses GOREST v1.1.19. Application settings use `GOREST_APP_*`; the HTTP listener and middleware use the `GOREST_SERVER_*` variables in `.env.example`. The default local bind is `127.0.0.1:3000`. Set the bind address intentionally when deploying behind a proxy or container network.

Server listener options include TCP/Unix network, TLS certificate files, TLS minimum version, shutdown timeout, Unix socket mode, and prefork settings. Blank listener values use the defaults provided by the GOREST version selected in `go.mod`. See `.env.example` and the [GOREST configuration guide](https://github.com/elsyahtech/gorest/blob/main/docs/configuration.md).

## MySQL

MySQL connection variables are documented in the [MySQL configuration guide](./database.md). Use TLS where required by the network and verify what the selected GOREST/MySQL driver configuration actually validates. In particular, do not claim client-certificate authentication unless the connector demonstrably loads and presents that certificate.

## Redis

Redis wiring is retained as an optional GOREST service. `GOREST_REDIS_DISABLE=true` is the starter default, both in the loader and `.env.example`. Set it to `false` only after configuring a reachable Redis service. No current user endpoint depends on Redis, and unit or MySQL integration tests do not require it.

## Authentication and sample limits

Replace the example `GOREST_FILTER_SIGNATURE_KEY` with a unique, high-entropy secret. Review cookie security settings and use HTTPS in deployments where credentials or cookies cross an untrusted network.

The create and login endpoints are public examples. Only `GET /users/` uses authentication; the example does not implement authorization roles or protect business resources. Add application-specific access control before exposing it.

## Monitoring and logs

Monitoring endpoints are managed by GOREST and configured through `GOREST_MONITORING_*`. The starter defaults to loopback access. Restrict monitoring at the network or proxy boundary as well. Choose a log directory and retention policy appropriate to the deployment, and avoid logging credentials or tokens.

For repeatable tests and coverage evidence, see [Testing and coverage](./testing.md).
