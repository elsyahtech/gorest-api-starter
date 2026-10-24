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

- The `/examples` routes do not implement authentication or authorization. Do not expose them to untrusted clients without adding appropriate access controls.
- The checked-in application configuration defaults to SQLite and contains no production database credentials. The optional driver-test Compose file uses disposable local-only test credentials; do not reuse them. Runtime connection settings are read from `GOREST_DB_*` and optional Redis settings from `GOREST_REDIS_*`; the application does not load `.env` files. Inject production secrets through a protected environment or secret manager. If a real secret has ever been committed or used in a deployment, rotate it.
- `Migration=true` / `Seeder=true` are local-example defaults only. With `APP_ENV=production`, both default to false and startup rejects either being enabled. Apply schema changes through a controlled deployment step; do not run sample seeders against production data.
- Database TLS is driver-specific. GOREST emits an informational production warning for configurations where transport encryption is disabled or not guaranteed, but the warning does not block startup or configure TLS. Review the [per-driver TLS guidance](./docs/database.md#production-transport-encryption-by-driver) and apply the policy appropriate to your database and network.
- CORS, security headers, request limits, and the in-memory rate limiter are not substitutes for authentication, authorization, TLS, or deployment-level controls. The in-memory limiter applies per process.
- Native SQL examples are driver-specific; ORM and transaction behavior differs between database drivers. Passing the test suite does not establish that every driver or production topology is secure.

See [Configuration, security, and deployment](./docs/configuration.md) for deployment guidance and [Testing and E2E coverage](./docs/testing.md) for the test scope and driver limitations.

If a report concerns a vulnerability in GOREST or another dependency rather than this starter's integration, please also report it to the relevant upstream project through its private reporting channel.
