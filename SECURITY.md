# Security Policy

## Supported versions

Security fixes are targeted at the latest code on the `main` branch. No maintenance or security-fix commitment is made for older commits or release artifacts unless a release explicitly says otherwise. Update to the latest `main` revision before reporting an issue.

## Reporting a vulnerability

Please report suspected vulnerabilities privately. Do not disclose vulnerability details in a public GitHub issue, pull request, discussion, or social media post before the maintainers have had a reasonable opportunity to investigate and coordinate a fix.

The preferred channel is GitHub's private vulnerability reporting for this repository:

[Report a vulnerability privately](https://github.com/elsyahtech/gorest-api-starter/security/advisories/new)

If private reporting is unavailable, open a public issue asking the maintainers for a private reporting channel. Do not include exploit details, sensitive information, or a proof of concept in that request.

Include as much of the following as you can safely provide:

- A concise description and potential impact.
- The affected commit, branch, or dependency version.
- The relevant configuration and environment, with credentials and personal data removed.
- Reproduction steps or a minimal proof of concept, if available.
- Any mitigation or workaround you have identified.

Please use a test environment you control. Do not access, modify, or expose another person's data; disrupt services; or test against infrastructure without authorization. Avoid sending real secrets or production data in a report.

Maintainers will acknowledge and investigate reports on a best-effort basis. There is no guaranteed response or remediation timeline, and the project does not currently promise a bug bounty. We will coordinate public disclosure with the reporter when practical.

## Scope and important context

This repository is a learning and implementation starter, not a production-secured service. In particular:

- `POST /users/create` and `POST /users/login` are public examples; `GET /users/` requires an access cookie. This demonstration does not provide authorization roles or protect other business resources. Add application-specific access controls before exposing it to untrusted clients.
- The starter is configured for MongoDB. `.env.example` contains placeholders, `.env` is ignored by Git, and the application loads `.env` when present. Inject production secrets through a protected environment or secret manager. If a real secret has ever been committed or used in a deployment, rotate it.
- MongoDB migrations and sample seeders are disabled by default in `.env.example`. They write to the configured database; run them only as a controlled deployment step and never run the sample seeder against production data.
- MongoDB TLS is disabled by default for local development. Enable TLS and validate the server certificate for remote deployments. Client certificate and key paths are only for mTLS; protect the private key. See the [MongoDB configuration guide](./docs/database.md).
- The integration suite performs migrations, inserts three sample documents, and writes test documents. It requires a dedicated non-production test database and exact database-name confirmation; never point it at production or data that must be preserved.
- Redis is optional and disabled by default. CORS, security headers, request limits, and the in-memory rate limiter are not substitutes for authentication, authorization, TLS, or deployment-level controls. The in-memory limiter applies per process.
- MongoDB transaction support depends on server topology; transactions require a replica set or sharded cluster. Passing the test suite does not establish that every production topology is secure.

See [Configuration and deployment](./docs/configuration.md) and [Testing](./docs/testing.md) for additional guidance.

If a report concerns a vulnerability in GOREST or another dependency rather than this starter's integration, please also report it to the relevant upstream project through its private reporting channel.
