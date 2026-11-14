# Benchmark plan and measurement matrix

> **Draft — no benchmark implementation or measurements are included yet.** This plan separates GOREST framework overhead from the starter's middleware, business logic, ScyllaDB, and password-hashing costs. It defines what to measure and report; it makes no performance claim.

## Planned package layout

The benchmark code will live under this directory. The initial scaffold only establishes package boundaries; it contains no benchmark workload, runner, external-service connection, or result.

```text
benchmark/
  README.md
  doc.go
  framework/              F0-F1 first; F3-F5 and F7-F8 in later phases
  middleware/             F2 and F6, including individual and combined profiles
  starter/                A0-A3; A4-A5 only if the starter adds those endpoints
  scylla/                  D0-D2, independent of the GOREST HTTP path
  internal/
    fixtures/             deterministic request/response and dataset fixtures
    harness/              shared setup, load orchestration, and correctness checks
    manifest/             version, machine, source, and run metadata
    report/               raw-result parsing, summaries, and chart inputs
  testdata/
    framework/            checked-in, non-sensitive framework payload fixtures
    starter/              checked-in, synthetic starter payload fixtures
  results/
    README.md             raw-run retention and sanitization rules
```

Keep scenario implementations in their workload package and reusable mechanics in `internal/`. Framework benchmarks must not import starter or ScyllaDB packages. Starter packages may use the application; database-reference benchmarks use the ScyllaDB driver directly. Generated reports must identify the source commit and retain raw runs, but must never contain `.env`, secrets, cookies, or private-key/certificate material. Empty Go package files are intentional scaffolding, not completed benchmarks.

## Questions this benchmark should answer

1. What is the relative HTTP overhead of GOREST compared with using the same Fiber version directly?
2. How do common request/response work and selected middleware affect that overhead?
3. What throughput and latency does the starter deliver for its actual user endpoints under a documented ScyllaDB setup?
4. Which costs come from the framework, application code, password hashing, or the database?

These are separate questions. Results from a database-backed or password-hashing endpoint must not be presented as framework-only performance.

## Workload matrix

| ID | Scope | Request/workload | Comparison | Included work | Explicitly excluded or reported separately | Purpose |
|---|---|---|---|---|---|---|
| F0 | Framework baseline | `GET /bench/plain`, fixed small JSON response | GOREST vs Fiber v3 directly | Router, handler dispatch, response writing; same Go version and equivalent server settings | Database, authentication, logging I/O, compression, TLS | Estimate minimum framework/server-path overhead |
| F1 | Framework request path | `POST /bench/echo`, fixed JSON body decoded and equivalent JSON response encoded | GOREST vs Fiber v3 directly | Route matching, request parsing, response serialization | Database, auth, password hashing, unrelated middleware | Compare a representative non-trivial HTTP handler path |
| F2 | Middleware profiles | Same F0/F1 requests with middleware enabled | GOREST vs equivalent Fiber middleware; measure each middleware individually and a separately named combined profile | Only middleware enabled in the profile; report exact config | Database and user business logic | Measure middleware cost without attributing it to routing |
| F3 | Larger JSON response | `GET /bench/json`, deterministic nested JSON response at documented payload sizes | GOREST vs direct Fiber with identical response bytes | Serialization and response writing | Database and user business logic | Measure larger-payload serialization and writes |
| F4 | Route/query handling | `GET /bench/items/:id?limit=...`, fixed valid and invalid inputs | GOREST vs direct Fiber with equivalent route and handler | Route matching, parameter/query parsing, response | Database and auth unless a separate profile is declared | Measure realistic route dispatch and request handling |
| F5 | Error/recovery path | Deterministic 4xx response and separately controlled 5xx/recovery response | GOREST vs direct Fiber with equivalent error behavior | Error response path; recovery only in its named profile | Normal success path, database | Measure error handling without mixing it with successful throughput |
| F6 | Request ID and recovery | F0/F1 with request ID and recovery middleware enabled | GOREST vs Fiber equivalents; each individually and as a combined profile | Request-ID generation/propagation and recovery middleware | Other middleware unless identified in the profile | Isolate selected middleware costs |
| F7 | Large request body | `POST /bench/echo`, same deterministic JSON structure at each payload profile | GOREST vs direct Fiber with identical request and response bytes | Body read, parse, and response serialization | Database, auth, unrelated middleware | Measure request handling as body size grows |
| F8 | Route-table cardinality | Same request/handler registered in route tables of increasing size; measure routes near beginning, middle, and end, plus a miss | GOREST vs direct Fiber with equivalent route sets and handler behavior | Route registration/matching and dispatch | Database, response-size changes, middleware unless named | Measure whether registered route count/order affects dispatch |
| A0 | Starter HTTP path, no external I/O | A benchmark-only route wired through the starter's GOREST server/config and response envelope | GOREST starter stack vs F1/F2 baseline | Actual starter configuration and selected middleware, route, response envelope | ScyllaDB, Redis, user library, Argon2 | Estimate overhead added by starter wiring/configuration |
| A1 | User listing | `GET /users/`, valid access cookie required; fixed dataset sizes (100 and 1,000 rows, with 10,000 as an optional scale profile if service limits permit); same ScyllaDB keyspace/TLS/session-pool settings as declared | Actual starter endpoint | Authentication filter, handler/library/model, ScyllaDB read, DTO mapping, response | Password hashing; Redis | Measure the protected read path; report DB contribution and response growth |
| A2 | User creation | `POST /users/create`, public route; unique synthetic ID/email/password per request; dedicated ScyllaDB test DB | Actual starter endpoint | JSON parsing, Argon2id hash, insert, response | Do not compare directly with F0–F6 as framework performance | Measure real write behavior and resource costs |
| A3 | Login | `POST /users/login`, public route; pre-created synthetic account with fixed valid credentials and hash | Actual starter endpoint | JSON parsing, ScyllaDB lookup, Argon2id verification, token/cookie creation | Do not compare directly with framework-only rows | Measure authentication behavior; report hash parameters and DB contribution |
| A4 | Update user | Not currently supported by the starter's registered API | Not applicable | — | Do not implement as part of this benchmark until an actual endpoint exists | Prevent benchmarking invented behavior |
| A5 | Delete user | Not currently supported by the starter's registered API | Not applicable | — | Do not implement as part of this benchmark until an actual endpoint exists | Prevent benchmarking invented behavior |
| D0 | Direct DB read/insert | Fixed parameterized CQL read and insert using the native ScyllaDB driver, outside GOREST HTTP/application layers | Same ScyllaDB version/host/keyspace/dataset/TLS and session-pool settings as relevant A scenarios | ScyllaDB round-trip and CQL operation | GOREST HTTP routing/middleware | Establish an independent DB reference |
| D1 | Session-pool reuse | Repeated fixed read with declared session/connection-pool settings using the native ScyllaDB driver | Compare declared pool/reuse profiles; keep query and DB target fixed | Driver and session-pool behavior plus DB read | GOREST HTTP path | Evaluate connection reuse and pool effects |
| D2 | Concurrent read/write | Bounded concurrent parameterized CQL reads and unique-row writes using the native ScyllaDB driver | Same keyspace and declared pool settings; run as a separate DB workload | ScyllaDB concurrency and driver behavior | GOREST HTTP path | Evaluate DB behavior under parallel load |

For A1–A3, record the exact fixture/dataset size, sanitized database placement/version/TLS/pool configuration, request payload size, authentication state, expected status and response schema, and cleanup method. A2 requires a unique record per successful request and independent row-count verification. A3 uses a pre-created account; never include its credential or access cookie in artifacts. A2/A3 use the current Argon2id settings (`m=65536 KiB, t=3, p=4`, 32-byte output, 16-byte salt) only if those values are confirmed in the exact tested source/dependency version; record any change instead of assuming them.

### Proposed payload and scale profiles

These are starting profiles for Auditor review, not mandatory targets. Measure fixture sizes in bytes, keep them identical between GOREST and Fiber, and report actual sizes with every result. Adjust only when documented server limits, available memory, or test-service capacity require it; do not silently change a profile between implementations.

| Profile | F1/F7 request body | F3 response body | A1 ScyllaDB dataset | Notes |
|---|---:|---:|---:|---|
| Small / baseline | 1 KiB | 1 KiB | 100 rows | Establish the low-payload baseline |
| Medium | 64 KiB | 64 KiB | 1,000 rows | Show behavior under a moderate increase |
| Large / scale | 1 MiB | 1 MiB | 10,000 rows (optional) | Run only after verifying body/response limits, memory headroom, DB capacity, and cleanup; otherwise mark not run with the reason |

F1 remains a small, fixed-body request for the initial Phase 1 comparison. F3 varies response size; F7 varies request size while holding response behavior constant or explicitly documenting the response-size rule. A1 dataset-size results are application-plus-database measurements, not framework-only results. Dataset growth also changes response bytes because the current list endpoint returns users; report row count and response bytes together, and do not treat it as a pure database-size comparison.

### Starter request contracts

These are correctness expectations for the current routes, not measurements. Confirm them against the exact tested commit before implementing the harness.

| Scenario | Authentication | Dataset and database | Request | Expected response and independent check |
|---|---|---|---|---|
| A0 | None unless the benchmark route explicitly tests auth | No external I/O | Fixed deterministic request defined by the harness | Fixed status, response schema, and bytes shared with its F1/F2 comparison |
| A1 | Valid access cookie from a prepared synthetic account | Fixed ScyllaDB keyspace dataset; recommend 100 users for baseline and optionally 1,000 as a separate scale profile; report exact count, TLS, and session-pool settings | `GET /users/` | HTTP 200; response contains expected rows and public fields, no password hashes; independently compare row count/fixture IDs |
| A2 | Public route | Dedicated ScyllaDB test database; one unique synthetic ID/email per offered operation | `POST /users/create` with fixed-size JSON user/password payload | HTTP 200 and success envelope; independently verify exactly one row with expected fields and a non-plaintext password hash; clean only that test row |
| A3 | Public route; valid synthetic account created before measurement | Fixed account/hash in the same declared ScyllaDB configuration | `POST /users/login` with fixed-size JSON credentials | HTTP 200, success envelope, and non-empty access cookie; do not record cookie or credentials |

For A2, a request not confirmed as successful must not be entered as a successful write. For A3, verify the cookie's presence without retaining its value. Record request/response sizes and exact correctness assertions for all rows.

### Comparison rules

- **Framework comparison:** GOREST v1.1.20 uses Fiber v3.6.1 in this repository. Pin and report both versions. Use the same Go toolchain, handler logic, payloads, response bytes, server settings, and middleware semantics for F0/F1. The direct-Fiber implementation is a measurement control, not a production alternative claim.
- **Middleware comparison:** keep an “off” baseline and add middleware in named profiles. Do not compare different middleware behavior as if it were equivalent. Report middleware names and settings, including compression, request ID, CORS, Helmet, recovery, logging, and authentication where applicable.
- **Starter endpoint comparisons:** A1–A3 describe end-to-end application behavior. Label them as application results, not GOREST-only results. Keep A2/A3 separate because Argon2 is deliberately compute-intensive.
- **Database comparisons:** keep ScyllaDB on the same host/network, version, keyspace/table layout, dataset, TLS mode, and session-pool configuration for compared runs. Report whether it is local or remote. Do not mix DB results into framework-only throughput.
- Do not add a competing framework, change production defaults, or make performance optimizations solely to improve benchmark numbers without a separately reviewed scope.

## Load and run matrix

| Dimension | Draft levels / rule | Record in report |
|---|---|---|
| Execution mode | Go microbenchmarks for isolated in-process operations; HTTP load tests for complete server requests. Do not combine their numbers. | Tool, exact command, tool version, benchmark/test name |
| Concurrency-based levels | 1 (serial baseline), 8, 32, then 64 or 128 only when hardware and service limits allow | Configured and achieved concurrency |
| Rate-based levels | Gradually increase offered request rate from below saturation until throughput stops scaling or latency/errors exceed the declared stopping rule | Offered, completed, and successful request rates at each rate; latency and errors |
| Duration | Warm-up before measurement; target at least 30 seconds per HTTP measurement and at least 5 independent repetitions | Warm-up, measured duration, repetition count |
| Connection behavior | HTTP keep-alive enabled for the primary profile; run a separate no-reuse profile only if relevant | Keep-alive, connection count, protocol |
| Payload | Fixed, checked-in or generated deterministic body and response for each scenario; same bytes across framework controls | Request/response size and body fixture identifier |
| Dataset | Fixed documented row count and distribution for A1; unique IDs for A2; fixed account/hash settings for A3 | Dataset size, setup method, cleanup method; never include credentials |
| Transport | Loopback/plain HTTP for core comparison; TLS as a separate profile if deployment uses it | TLS version, certificate verification mode, endpoint placement |
| Runtime | Release build, no debugger/profiler in primary run; profiling runs separate | Go version, OS/architecture, build flags, GOMAXPROCS, GC settings |
| Repetitions | At least 5 independent runs per scenario/profile; randomize or alternate comparison order to reduce drift | Per-run values, median and spread; do not retain only the best run |
| Isolation | Load generator should not contend with the server; prefer separate hosts/processors for HTTP load, and avoid unrelated jobs | Server/load-generator host details and background-load notes |

Do not run framework alternatives simultaneously on a constrained host. Keep warm-up and measurement windows identical. If the machine or external ScyllaDB/Redis service is shared or noisy, record that and repeat rather than hiding outliers.

For rate-based tests, define offered rate as requests scheduled by the load generator, completed rate as requests receiving a complete response, and successful rate as completed responses that pass the scenario's status/body correctness checks. Report all three. Report errors as both a count and a percentage of offered requests; classify HTTP status failures, timeouts, connection errors, and invalid responses separately.

## Metrics matrix

| Metric | Unit / reporting | Applies to | Interpretation |
|---|---|---|---|
| Throughput | Offered, completed, and successful requests/second; report all three per run and median | HTTP workloads | Distinguishes offered load, work completed, and valid responses |
| Latency | p50, p95, p99, and max in milliseconds; per run and median across runs | HTTP workloads | Report client-observed latency and identify the load generator; avoid averages alone |
| Errors | Count and percentage by status/transport error; expected and unexpected failures separated | All HTTP workloads | A fast run that drops or rejects requests is not a successful result |
| CPU | CPU utilization and/or CPU time over the measured interval; where possible CPU time per 1,000 successful requests | HTTP and DB workloads | Include host/process attribution and normalization method |
| Memory | Peak and steady-state process RSS | HTTP and DB workloads | Note sampling method and whether it is process or container memory |
| Allocations | `ns/op`, `B/op`, `allocs/op` for Go microbenchmarks; server allocation metrics only with validated low-overhead method | F0–F6 and targeted code paths | Record instrumentation overhead; not a substitute for end-to-end latency |
| Go runtime | Goroutines, GC cycles, total GC pause, heap allocation, heap system memory, runtime memory allocation, GC CPU fraction | Server workloads | Capture at documented intervals; do not profile during primary runs |
| Database | Query and transaction latency, pool settings, DB CPU/resource use where measurable, query/index details | A1–A3, D0–D2 | Helps distinguish remote DB and pool bottlenecks from application work |
| Hashing | Argon2 parameters and isolated hash/verify duration, reported separately | A2/A3 | Dominant by design in create/login; do not ascribe it to GOREST |
| Correctness | Status, response schema/body checks, created-row verification and cleanup | Every scenario | Invalid responses or data make performance numbers unusable |

Use the Go benchmark comparison tool (for example `benchstat`) for repeated microbenchmark output. For HTTP load, retain raw tool output and calculate percentiles with a documented, pinned load generator. The implementation phase must choose the actual tools and pin their versions before collecting results.

Do not run CPU profiling, memory profiling, or tracing during primary measurements. Collect profiles and traces in separate diagnostic runs and label them as such; their instrumentation can change performance.

## Correctness gates

Do not accept performance data unless every scenario validates the expected HTTP status, response body/schema, and number of successful operations. For reads, verify returned data/count against the fixed dataset. For writes, independently query ScyllaDB after the run to verify expected row counts, uniqueness, and content, then verify test-data cleanup. Report timeouts, connection failures, invalid responses, duplicate records, and aborted runs explicitly; never count them as successful operations or silently discard failed runs.

## Environment and controls to capture

| Category | Required record |
|---|---|
| Source | Repository, branch, exact commit SHA, dirty/clean worktree |
| Software | Verify and record exact Go, GOREST, Fiber, ScyllaDB server/driver, load generator, and profiler versions against the tested commit and `go.mod` where applicable |
| Hardware | CPU model/core count, memory, architecture, OS/kernel, VM/container limits if any |
| Process | Release build mode/flags, `GOMAXPROCS`, `GOGC`, GC settings, CPU quota, process/container limits |
| Server | Bind/port, protocol, keep-alive, listener settings, middleware profile and relevant GOREST config |
| ScyllaDB | Local/remote placement, version, TLS mode, sanitized keyspace/table data description, session-pool settings, query/index details |
| Redis | Whether disabled or enabled; if enabled, endpoint placement and TLS mode. Do not include hostnames, usernames, or secrets in public artifacts unless approved |
| Test data | Dataset generation, row counts, request and response payload sizes/fixtures, reset/cleanup steps |
| Run conditions | Time/date, warm-up, duration, repetitions, order, background load, failed/aborted runs and reason |

Never put `.env`, credentials, access cookies, private keys, or certificate contents in benchmark output. Use sanitized configuration and synthetic test identities. Preserve every raw run, including failed or aborted runs, rather than replacing them with summary numbers.

## Results table template

Fill one row per scenario, implementation, concurrency, and repetition. Keep raw per-run data; summarize only after validating correctness and error rate.

| Scenario | Implementation/profile | Commit | Concurrency/offered RPS | Run | Duration | Offered RPS | Completed RPS | Successful RPS | p50/p95/p99/max ms | Error count/% | CPU / CPU per 1k success | Steady/peak RSS | ns/op, B/op, allocs/op | Notes |
|---|---|---|---:|---:|---:|---:|---:|---:|---|---|---|---|---|---|
| F0 | GOREST minimal | — | — | — | — | — | — | — | — | — | — | — | — | Not run |
| F0 | Fiber direct | — | — | — | — | — | — | — | — | — | — | — | — | Not run |
| A1 | Starter GET `/users/` + ScyllaDB | — | — | — | — | — | — | — | — | — | — | — | N/A | Not run |

The table above summarizes measurements only; retain raw tool output for every run. Also publish a machine-readable environment/run manifest, deterministic request/response fixtures, exact commands and scripts/Makefile targets, scenario execution status, and a concise reproduction guide. Provide throughput-vs-latency, concurrency-vs-throughput, and CPU/memory charts generated from the retained raw data. Charts must identify units, scenario/profile, source commit, and whether values are individual runs or summaries.

For summarized comparisons, report median and range (or a confidence interval if the chosen method supports it), absolute difference, relative difference, and error-rate difference. State the denominator and formula for every percentage (for example, `(candidate - baseline) / baseline * 100%`, with the direction labeled). Do not claim a meaningful win from a single run or a difference smaller than observed run-to-run variation.

## Decision gates before implementation

| Decision | Recommended default | Why it matters |
|---|---|---|
| Framework comparator | Fiber v3.6.1, the Fiber version already used by this project | Controls for the underlying HTTP framework and avoids introducing another dependency |
| Starter scope | Include one no-DB route plus actual ScyllaDB-backed read/write/login scenarios | Separates framework/server overhead from real application costs |
| Load generator | Choose and pin a maintained HTTP tool during implementation | Reproducible client-side latency and throughput; tool choice can affect results |
| Test database | Dedicated non-production schema with documented reset/cleanup and least-privilege account | Real user endpoint tests migrate/write data |
| Redis | Excluded from the primary user API matrix; optional separate service workload | Current endpoints do not depend on Redis |
| Performance threshold | No pass/fail threshold in the initial draft | First collect reproducible baseline data; avoid arbitrary targets or pre-judging results |
| Publication | Commit the benchmark harness, exact commands, raw outputs, environment manifest, and concise summary | Lets reviewers rerun and assess the evidence |
| Rate-based load | Add a controlled offered-rate sweep alongside concurrency-based tests | Reveals saturation and separates offered from completed/successful throughput |
| Starter write expansion | Do not implement A4/A5 until actual update/delete endpoints exist | The current starter registers create, login, and list only |
| Payload/dataset scale | Start with 1 KiB, 64 KiB, and 1 MiB HTTP fixtures; use 100 and 1,000 A1 rows, with 10,000 rows optional | Makes large-request and large-data workloads explicit while respecting server and DB limits |
| Route-table size | Compare 10, 100, and 1,000 equivalent routes for F8; if route count or registration exceeds a framework limit, record the limit and stop rather than changing semantics | Tests whether GOREST route-table size affects dispatch |

## Acceptance criteria

The benchmark is ready for auditor review when framework/direct-Fiber runs have equivalent workloads; middleware is independently measurable; starter endpoints and independent ScyllaDB references are reported separately; methodology and correctness are validated; CPU, memory, latency, throughput, and errors are collected; at least five independent runs and all raw data are retained; environment details and exact reproduction commands are available; and no unsupported performance claims are made.

## Implementation sequence

1. **Framework baseline:** implement F0/F1 for GOREST and direct Fiber using the small profile; validate the measurement method and submit initial results for review before expanding.
2. **Middleware and framework scale:** implement F2, F5, and F6, then F3/F7 payload sizes and F8 route cardinality as separately reported workloads.
3. **Starter:** implement A0–A3. A4/A5 remain out of scope unless actual update/delete endpoints are added and reviewed.
4. **Database:** implement D0–D2 outside the GOREST HTTP path.
5. **Runtime diagnostics and reporting:** add resource metrics, charts, raw-data retention, documentation, and reproduction checks. Use separate diagnostic runs for profiles/traces.

Do not implement every scenario at once. The next approval gate is review of Phase 1 methodology and initial F0/F1 results. Do not claim a meaningful advantage when the difference is smaller than observed run-to-run variation, and do not describe the highest throughput as framework capacity without latency and error rate at that load.