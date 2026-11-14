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

- `POST /users/create` and `POST /users/login` are public examples; only `GET /users/` requires authentication. The starter does not implement authorization roles or protect other resources. Add application-specific access controls before exposing it to untrusted clients.
- The application configuration defaults to ScyllaDB and contains no production database credentials. Runtime settings use `GOREST_DB_*` and optional Redis settings use `GOREST_REDIS_*`. The executable loads `.env` when present; protect that file and inject production secrets through a secret manager or protected environment. If a real secret has ever been committed or used in a deployment, rotate it.
- Migrations and seeders are disabled by default in `.env.example`. Run schema changes as a controlled deployment step, and do not run sample seeders against production data.
- Database TLS is driver-specific. Review the [ScyllaDB TLS guidance](./docs/database.md#tls-and-client-certificates) and verify the configured endpoint and trust chain for your deployment.
- CORS, security headers, request limits, and the in-memory rate limiter are not substitutes for authentication, authorization, TLS, or deployment-level controls. The in-memory limiter applies per process.
- Native SQL examples are driver-specific; ORM and transaction behavior differs between database drivers. Passing the test suite does not establish that every driver or production topology is secure.

See [Configuration, security, and deployment](./docs/configuration.md) for deployment guidance and [Testing and coverage](./docs/testing.md) for the test scope and ScyllaDB integration safeguards.

If a report concerns a vulnerability in GOREST or another dependency rather than this starter's integration, please also report it to the relevant upstream project through its private reporting channel.
