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

- The active example exposes `POST /users/create` and `POST /users/login` without authentication; `GET /users/` requires a valid access cookie. This is demonstration authentication, not a complete authorization system. Add application-specific access controls before exposing user data to untrusted clients.
- This starter is configured for Oracle. `.env.example` contains placeholders, not production credentials, and `.env` is ignored by Git. The application loads `.env` when present; for production, inject secrets through a protected environment or secret manager. If a real secret has ever been committed or used in a deployment, rotate it.
- Oracle migrations and sample seeders are disabled by default in `.env.example`. Enabling `GOREST_DB_MIGRATION` or `GOREST_DB_SEEDER` writes to the configured Oracle schema; use a controlled deployment step and never run sample seeders against production data.
- Oracle TLS uses TCPS and an Oracle wallet. The example defaults to unencrypted TCP for local testing only. Configure and verify the wallet for remote production connections; see the [Oracle configuration guide](./docs/database.md).
- Redis is optional and disabled by default in `.env.example`. If enabled, configure its connection and transport security separately; the user endpoints do not depend on Redis.
- CORS, security headers, request limits, and the in-memory rate limiter are not substitutes for authentication, authorization, TLS, or deployment-level controls. The in-memory limiter applies per process.
- Passing the test suite does not establish that every Oracle deployment topology is secure.

See [Configuration, security, and deployment](./docs/configuration.md) and the [testing guide](./docs/testing.md) for more detail.

If a report concerns a vulnerability in GOREST or another dependency rather than this starter's integration, please also report it to the relevant upstream project through its private reporting channel.

##