# Go, Microservices, and DevOps Architect Journey

A 24-week, project-driven path for a frontend developer to build backend depth,
operational skill, and software-architecture judgment using Didar.

> The goal is not to collect tools or split everything into services. The goal is
> to make sound decisions, implement them, operate the result, and explain the
> trade-offs. Six months can build strong evidence for backend ownership and
> technical-lead scope; an architect title still requires sustained production
> impact and organizational trust.

## Commitment and weekly rhythm

Plan for **12–15 focused hours each week**:

- **Study — 4 hours:** official documentation, a focused chapter, and repository reading.
- **Build — 7 hours:** one small, working change in Didar.
- **Verify and document — 2 hours:** tests, diagrams, ADRs, or a runbook.
- **Review — 1–2 hours:** explain the work, inspect trade-offs, and plan the next week.

Use one branch per weekly deliverable. Keep changes reviewable. Do not move on
until the weekly quality gate passes.

## Assigned YouTube track

Videos count toward the weekly **4 study hours**; they are not extra work. Watch
only the assigned chapters, normally for **45–90 minutes per week**, and spend
the remaining study time reading official documentation and Didar's code. Pause
to take notes and reproduce examples instead of watching passively. Long courses
are intentionally divided across several weeks.

| Week | Assigned video | Watch for |
|---|---|---|
| 1 | [Go Programming – Golang Course with Bonus Projects](https://www.youtube.com/watch?v=un6ZyFkqFKo) | Intro, setup, variables, functions, packages, and basic HTTP |
| 2 | [Go Programming – Golang Course with Bonus Projects](https://www.youtube.com/watch?v=un6ZyFkqFKo) | Structs, interfaces, pointers, errors, and generics; stop after the relevant chapters |
| 3 | [Learn Go by Building 11 Projects](https://www.youtube.com/watch?v=jFfo23yIWac) | Web server and API project chapters; compare their layering with Didar |
| 4 | [Hussein Nasser's HTTP playlist](https://www.youtube.com/playlist?list=PLQnljOFTspQU6zO0drAYHFtkkyfNJw1IO) | Authentication headers, request lifecycle, timeouts, and status codes |
| 5 | [Hussein Nasser's PostgreSQL playlist](https://www.youtube.com/playlist?list=PLQnljOFTspQWGrOqslniFlRcwxyY94cjj) | Choose the indexing and query-plan videos relevant to the city query |
| 6 | [Go Programming – Golang Course with Bonus Projects](https://www.youtube.com/watch?v=un6ZyFkqFKo) | Goroutines, channels, mutexes, and concurrency exercises |
| 7 | [Hussein Nasser's Backend Engineering path](https://www.husseinnasser.com/2020/04/system-design-and-backend-engineering.html) | Process lifecycle, connections, timeouts, and graceful failure |
| 8 | [Go Testing Tutorials](https://www.youtube.com/results?search_query=Go+testing+table+driven+tests+httptest+official) | Select a recent focused lesson on table-driven tests, `httptest`, benchmarks, and `pprof` |
| 9 | [Simon Brown on mapping software architecture](https://www.youtube.com/watch?v=NhPa1t7M_4E) | C4 abstraction levels, audience, notation, and diagram scope |
| 10 | [Hussein Nasser's Load Balancing and Proxies playlist](https://www.youtube.com/playlist?list=PLQnljOFTspQVMeBmWI2AhxULWEeo7AaMC) | Reverse proxies, connection failure, timeout budgets, and retries |
| 11 | [Saga and Transactional Outbox lessons](https://www.youtube.com/results?search_query=saga+pattern+transactional+outbox+microservices) | Compare choreography, orchestration, outbox, duplicate delivery, and compensation |
| 12 | [OWASP API Security lessons](https://www.youtube.com/results?search_query=OWASP+API+Security+Top+10+official) | Broken authorization, authentication abuse, rate limits, and unsafe consumption |
| 13 | [Docker Tutorial for Beginners](https://www.youtube.com/watch?v=3c-iBn73dDE) | Containers, layers, Dockerfiles, Compose, volumes, registry, and deployment chapters |
| 14 | [GitHub Actions Tutorial](https://www.youtube.com/watch?v=R8_veQiYBjI) | Workflow events, jobs, runners, tests, image builds, and secrets |
| 15 | [Docker Tutorial for Beginners](https://www.youtube.com/watch?v=3c-iBn73dDE) | Rewatch registry and deployment chapters; map them to immutable releases and rollback |
| 16 | [Terraform explained in 15 minutes](https://www.youtube.com/watch?v=l5k1ai_GBDE) | Providers, resources, state, plan/apply, and declarative infrastructure |
| 17 | [Kubernetes Crash Course](https://www.youtube.com/watch?v=s_o8dwzRlu4) | Core components, YAML, Deployments, Services, ConfigMaps, Secrets, and local setup |
| 18 | [Kubernetes Tutorial for Beginners](https://www.youtube.com/watch?v=X48VuDVv0do) | Namespaces, Ingress, volumes, StatefulSets, Services, and production concerns |
| 19 | [OpenTelemetry distributed tracing tutorial](https://www.youtube.com/watch?v=62eQ5BG8M3g) and [Prometheus tutorial](https://www.youtube.com/watch?v=h4Sl21AKiDg) | Trace propagation, RED metrics, labels, histograms, and dashboards |
| 20 | [Google SRE: SLI, SLO, and error-budget lessons](https://www.youtube.com/results?search_query=Google+SRE+SLI+SLO+error+budget) | User-focused indicators, error budgets, burn rates, and actionable alerts |
| 21 | [Hussein Nasser's Software Architecture playlist](https://www.youtube.com/playlist?list=PLQnljOFTspQXNP6mQchJVP3S-3oKGEuw9) | Select videos on architectural drivers, quality attributes, and trade-offs |
| 22 | [Hussein Nasser's Database Engineering playlist](https://www.youtube.com/playlist?list=PLQnljOFTspQXjD0HOzN7P2tgzu7scWpl2) | Capacity, caching, partitioning, consistency, and database bottlenecks |
| 23 | [Hussein Nasser's Software Architecture playlist](https://www.youtube.com/playlist?list=PLQnljOFTspQXNP6mQchJVP3S-3oKGEuw9) | Evolution, coupling, modular monoliths, microservices, and migration cost |
| 24 | [Software Architecture Design Review talks](https://www.youtube.com/results?search_query=software+architecture+design+review+presentation) | Watch one strong review and critique its structure, evidence, risks, and delivery |

For search-result assignments, prefer recent videos from the official project,
conference, or a recognized practitioner. Record the title and link you chose in
your weekly notes so the curriculum remains reproducible.

## Didar map

```text
frontend/              Vue 3 SPA
backend/               Main API: experiences, bookings, hosts, brands
services/auth/         Auth service: OTP, JWT, users, sessions
pkg/authkit/           Shared JWT and auth middleware
pkg/platform/          Shared HTTP helpers
backend/migrations/    Main database migrations
services/auth/migrations/
docker-compose.yml     Local multi-container environment
deploy/                Render all-in-one production image
.github/workflows/     Existing frontend deployment workflow
```

Local endpoints:

- Frontend: `http://localhost:5173`
- Main API: `http://localhost:8080/api/v1`
- Auth API: `http://localhost:8081/api/v1`
- PostgreSQL: `localhost:5432`

## Rules for the journey

1. Read the referenced code before changing it.
2. Every behavior change needs an appropriate test.
3. Record important decisions as: problem, constraints, options, choice,
   consequences, and a review date.
4. Prefer a modular monolith until an independently deployable service has a
   clear business or operational reason.
5. Never place real credentials in Git. Use examples or secret references.
6. At the end of every fourth week, give a short demo and explain one failure
   mode without notes.

---

## Phase 1 — Go and backend foundations (Weeks 1–4)

### Week 1 — Trace the auth service

**Study (4h)**

- Complete the relevant sections of [A Tour of Go](https://go.dev/tour/).
- Read `services/auth/cmd/server/main.go`,
  `services/auth/internal/httpserver/router.go`, and
  `services/auth/internal/config/config.go`.
- Trace `POST /auth/otp/request` from `frontend/src/views/LoginView.vue` to the
  auth handler, service, repository, and database.

**Build (7h)**

- Run `docker compose up --build`, log in, and locate the development OTP.
- Add `GET /api/v1/auth/ping` returning the service name and status.
- List all auth routes and classify them as public, protected, or internal.

**Verify and document (2h)**

- Add a handler test with `httptest`.
- Draw the request path and record configuration inputs.

**Quality gate:** the endpoint returns the documented JSON, its test passes, and
you can explain bootstrap and dependency wiring from memory.

**Evidence:** ping change, handler test, and `docs/diagrams/auth-request-flow.md`.

### Week 2 — Go types, interfaces, errors, and context

**Study (4h)**

- Read `services/auth/internal/user/model.go`,
  `services/auth/internal/otp/service.go`, and both OTP sender implementations.
- Learn structs, JSON tags, interfaces, pointer semantics, error wrapping,
  `errors.Is`, and `context.Context`.

**Build (7h)**

- Trace every error path in OTP request and verification.
- Replace one weak error path with a wrapped, classifiable error.
- Add request logging middleware with method, path, status, duration, and
  request ID. Extend the existing middleware stack; do not duplicate its output.

**Verify and document (2h)**

- Add table-driven tests for success, invalid input, and service failure.
- Explain why nullable fields use pointers and why request context must reach
  database and external-service calls.

**Quality gate:** errors remain useful across layers, cancellation is preserved,
and logs contain no OTP codes, tokens, or personal data.

**Evidence:** middleware, table-driven tests, and an error-handling note.

### Week 3 — Handler, repository, and SQL flow

**Study (4h)**

- Read `backend/cmd/server/main.go`, `backend/internal/httpserver/router.go`,
  and the files under `backend/internal/experience/`.
- Learn manual dependency injection, parameterized SQL, pagination, sorting,
  and stable API response shapes.

**Build (7h)**

- Trace `GET /api/v1/experiences?city=تهران` end to end.
- Add a validated `sort=price_asc|price_desc` option. Do not insert raw query
  values into SQL.
- Add a bounded page size and stable secondary sort.

**Verify and document (2h)**

- Test allowed, missing, and malicious sort values.
- Record the API contract and query-performance assumptions.

**Quality gate:** invalid values receive a consistent client error, SQL remains
parameterized, and repeated pagination has deterministic ordering.

**Evidence:** API enhancement, tests, and a contract note.

### Week 4 — Authentication and booking invariants

**Study (4h)**

- Read `pkg/authkit/jwt.go`, `pkg/authkit/middleware.go`,
  `backend/internal/booking/`, and `backend/internal/authclient/client.go`.
- Learn authentication versus authorization, JWT validation, transactions,
  validation boundaries, and service-call timeouts.

**Build (7h)**

- Require `guests_count >= 1` and enforce any existing capacity rule at the
  correct layer.
- Add an explicit timeout to the auth client if it is not already bounded.
- Verify callers cannot act as another user by changing request data.

**Verify and document (2h)**

- Add booking handler tests for unauthenticated, invalid, and valid requests.
- Demo Weeks 1–4 and explain a booking failure path.

**Quality gate:** invariants are tested, authorization uses trusted token
identity, and all outbound calls have a deadline.

**Evidence:** validated booking flow, tests, and Month 1 demo notes.

---

## Phase 2 — Production-grade Go (Weeks 5–8)

### Week 5 — Database design and migrations

**Study (4h)**

- Read both migration directories and both `internal/dbx/migrate.go` files.
- Inspect `backend/migrations/0011_microservice_auth_split.sql`.
- Learn indexes, transactions, constraints, query plans, PostGIS basics, and
  backward-compatible migration patterns.

**Build (7h)**

- Measure the experience city query with `EXPLAIN (ANALYZE, BUFFERS)`.
- Add an appropriately named city index migration only if the plan justifies it.
- Write rollback instructions; do not assume startup migration is always safe.

**Verify and document (2h)**

- Test migration on an empty database and a populated local database.
- Draw ownership of user, booking, and experience data.

**Quality gate:** the migration is repeatably tested and its benefit and write
cost are supported by evidence.

**Evidence:** migration, before/after query plans, and data-ownership diagram.

### Week 6 — Transactions, concurrency, and cancellation

**Study (4h)**

- Learn goroutines, channels, mutexes, race detection, transaction isolation,
  lost updates, and graceful cancellation.
- Identify where concurrency is useful and where it would add risk.

**Build (7h)**

- Make one booking operation atomic if it spans multiple writes.
- Create a small bounded worker exercise for a non-critical task; do not launch
  unbounded goroutines from handlers.
- Propagate cancellation through the full request path.

**Verify and document (2h)**

- Run `go test -race ./...`.
- Test rollback and cancellation behavior.

**Quality gate:** no leaked work continues after cancellation, the race detector
passes, and transaction boundaries are explained.

**Evidence:** transactional test, race result, and concurrency decision note.

### Week 7 — Service lifecycle and consistent APIs

**Study (4h)**

- Learn graceful HTTP shutdown, readiness versus liveness, signal handling,
  structured errors, OpenAPI, and configuration validation.
- Compare error responses across auth and main API.

**Build (7h)**

- Add graceful shutdown to both Go services.
- Separate readiness from liveness where database availability matters.
- Define consistent error envelopes and document current APIs with OpenAPI.

**Verify and document (2h)**

- Test shutdown during an in-flight request and startup with invalid config.
- Validate the OpenAPI document.

**Quality gate:** services reject bad configuration, stop without abruptly
dropping normal requests, and expose meaningful health semantics.

**Evidence:** lifecycle changes, OpenAPI file, and health-endpoint tests.

### Week 8 — Testing strategy and performance

**Study (4h)**

- Learn unit, handler, repository integration, migration, contract, and end-to-end
  testing; use each only where it gives useful confidence.
- Learn Go benchmarks, `pprof`, allocation profiles, and basic load testing.

**Build (7h)**

- Create a test pyramid for OTP and booking critical paths.
- Add Postgres-backed repository tests and one API contract test.
- Benchmark or load-test the experience-list endpoint and profile one bottleneck.

**Verify and document (2h)**

- Run all Go tests in a clean environment.
- Demo Month 2 and compare measured results with the Week 3 assumptions.

**Quality gate:** critical paths have layered tests, results are reproducible,
and any optimization is based on a profile.

**Evidence:** testing strategy, integration/contract tests, benchmark, and profile.

---

## Phase 3 — Microservice design and reliability (Weeks 9–12)

### Week 9 — Boundaries and architecture views

**Study (4h)**

- Analyze `services/auth/`, `backend/internal/authclient/client.go`,
  `pkg/authkit/`, and the auth-split migration.
- Learn bounded contexts, coupling, cohesion, data ownership, C4 diagrams, and
  evolutionary architecture.

**Build (7h)**

- Create C4 context and container diagrams for the current system.
- Write ADR 0001 reconstructing the auth-split decision, including alternatives.
- Identify one candidate capability that should remain in the main API.

**Verify and document (2h)**

- Review diagrams against actual routes, databases, and deployments.
- List the runtime and organizational costs of each service boundary.

**Quality gate:** every shown relationship exists, ownership is unambiguous, and
the ADR includes consequences rather than only benefits.

**Evidence:** `docs/architecture/` diagrams and `docs/adr/0001-auth-split.md`.

### Week 10 — Resilient synchronous communication

**Study (4h)**

- Learn timeout budgets, retries with jitter, idempotency, circuit breakers,
  bulkheads, error budgets, and retry storms.
- Inspect internal auth calls and gateway routing in `deploy/nginx.conf`.

**Build (7h)**

- Define the failure policy for the auth internal API.
- Add only the timeout/retry behavior supported by request semantics.
- Inject auth latency and outage; ensure the API fails predictably.

**Verify and document (2h)**

- Test timeout, connection failure, and recovery.
- Produce a sequence diagram and failure-mode table.

**Quality gate:** no infinite waits or unsafe retry, errors are observable, and
degraded behavior is documented.

**Evidence:** resilience tests, sequence diagram, and failure analysis.

### Week 11 — Events, idempotency, and consistency

**Study (4h)**

- Learn at-least-once delivery, deduplication, ordering, outbox, dead-letter
  queues, sagas, and eventual consistency.
- Compare an in-process handler, synchronous HTTP, and an event consumer for
  booking notifications.

**Build (7h)**

- Write an ADR selecting the simplest justified notification design.
- Implement a small booking-confirmation path with an idempotency key.
- If using events, include a transactional outbox; do not introduce Kafka only
  to satisfy a tool checklist.

**Verify and document (2h)**

- Test duplicate delivery, consumer failure, and recovery.
- Document consistency guarantees and user-visible delay.

**Quality gate:** duplicate processing is safe, data cannot be silently lost,
and the chosen design is proportional to Didar’s scale.

**Evidence:** ADR, implementation, tests, and consistency note.

### Week 12 — API evolution, contracts, and security

**Study (4h)**

- Learn API compatibility, versioning, consumer-driven contracts, OWASP API Top
  10, threat modeling, least privilege, secret rotation, and rate limiting.

**Build (7h)**

- Threat-model login and booking flows.
- Add contract checks between the main API and auth service.
- Add or design OTP rate limiting with clear limits and privacy-safe keys.

**Verify and document (2h)**

- Test broken contracts, replay/abuse cases, and authorization boundaries.
- Demo Month 3 and lead a 15-minute review of the notification ADR.

**Quality gate:** high-risk threats have controls or explicit follow-up, and an
API change cannot silently break its consumer.

**Evidence:** threat model, contract tests, rate-limit decision, and review notes.

---

## Phase 4 — Containers, CI/CD, and infrastructure (Weeks 13–16)

### Week 13 — Containers and deployment models

**Study (4h)**

- Read `docker-compose.yml`, all service Dockerfiles, `deploy/Dockerfile`,
  `deploy/start.sh`, `deploy/nginx.conf`, and `render.yaml`.
- Learn layers, build cache, image security, non-root execution, health checks,
  and supply-chain basics.

**Build (7h)**

- Build and inspect each image; reduce avoidable size or privilege.
- Run `db` and `auth` alone, then the complete environment.
- Compare local split containers, Render’s all-in-one image, and a future split
  deployment by cost, isolation, scaling, and operational complexity.

**Verify and document (2h)**

- Scan images and record justified findings.
- Write the deployment-model ADR.

**Quality gate:** images build reproducibly, contain no secrets, and the current
production compromise is explained rather than hidden.

**Evidence:** hardened image changes, scan report, and deployment ADR.

### Week 14 — Continuous integration

**Study (4h)**

- Read `.github/workflows/deploy-pages.yml`.
- Learn CI isolation, caching, least-privilege permissions, artifact provenance,
  vulnerability scanning, and branch protection.

**Build (7h)**

- Add backend CI for formatting, static analysis, tests, race checks, migration
  validation, dependency scanning, and image builds.
- Pin actions appropriately and avoid exposing secrets to untrusted changes.

**Verify and document (2h)**

- Trigger and repair an intentional test failure.
- Record expected checks and approximate pipeline duration.

**Quality gate:** a clean checkout produces the same result, failures block merge,
and required permissions are minimal.

**Evidence:** backend workflow, successful run, and CI policy.

### Week 15 — Delivery, rollback, and migrations

**Study (4h)**

- Learn immutable image tags, deployment strategies, feature flags, smoke tests,
  expand/contract migrations, rollback limitations, and recovery objectives.

**Build (7h)**

- Design build-once/promote-many delivery for Didar.
- Add a post-deploy smoke test and write a rollback runbook.
- Practice a failed application rollout and a migration-compatible rollback.

**Verify and document (2h)**

- Time recovery and record manual steps.
- State RTO and RPO assumptions for the learning environment.

**Quality gate:** the previous version can be restored without guessing, and
database compatibility is considered explicitly.

**Evidence:** delivery workflow/design, smoke test, and rollback runbook.

### Week 16 — Infrastructure as code

**Study (4h)**

- Learn Terraform fundamentals: providers, state, modules, plans, drift, secrets,
  and environment separation.

**Build (7h)**

- Define one small, repeatable non-production environment using Terraform or the
  chosen provider’s supported IaC.
- Add cost tags/labels, outputs, and documented bootstrap steps.
- Keep state and credentials outside Git.

**Verify and document (2h)**

- Create, change, detect drift in, and destroy the environment.
- Demo Month 4 and explain the all-in-one versus split deployment decision.

**Quality gate:** another developer can reproduce the environment from docs and
destruction does not remove shared or production resources.

**Evidence:** IaC, redacted plan output, cost note, and demo notes.

---

## Phase 5 — Kubernetes and observability (Weeks 17–20)

### Week 17 — Local Kubernetes

**Study (4h)**

- Learn Pods, Deployments, Services, Ingress, ConfigMaps, Secrets, namespaces,
  labels, and declarative reconciliation using kind or minikube.

**Build (7h)**

- Create `k8s/` manifests for namespace, Postgres, auth, API, frontend, and ingress.
- Use ConfigMaps for non-secrets and secret references for sensitive values.
- Deploy auth first, then the complete request path.

**Verify and document (2h)**

- Recreate the cluster from scratch and test `/health`.
- Map each Compose service to its Kubernetes objects.

**Quality gate:** manifests are repeatable, no plaintext secret is committed, and
service discovery works without fixed IP addresses.

**Evidence:** `k8s/`, setup instructions, and verified health output.

### Week 18 — Kubernetes reliability

**Study (4h)**

- Learn startup/readiness/liveness probes, requests/limits, rolling updates,
  autoscaling, disruption budgets, affinity, graceful termination, and Jobs.

**Build (7h)**

- Add appropriate probes and resource settings from measurements.
- Configure rolling updates and graceful termination.
- Decide how migrations run safely instead of running concurrently in replicas.

**Verify and document (2h)**

- Kill pods, perform a bad rollout, and verify rollback.
- Observe behavior during a dependency outage.

**Quality gate:** traffic reaches only ready pods, rollout failure is visible and
recoverable, and migration execution is single-owner.

**Evidence:** production-readiness manifests and failure-drill notes.

### Week 19 — Logs, metrics, and traces

**Study (4h)**

- Learn the three observability signals, RED/USE methods, cardinality, sampling,
  correlation IDs, OpenTelemetry, Prometheus, and trace context propagation.

**Build (7h)**

- Emit structured logs with request IDs across frontend gateway, API, and auth.
- Add request rate, error, duration, and dependency metrics.
- Trace login and booking across service boundaries.

**Verify and document (2h)**

- Locate one failed request using logs, metrics, and its trace.
- Check that telemetry contains no tokens, OTPs, or unnecessary PII.

**Quality gate:** one correlation value connects signals, metric labels are
bounded, and telemetry answers where time and errors occur.

**Evidence:** instrumentation, dashboard, and end-to-end trace capture.

### Week 20 — SLOs, alerts, and incident response

**Study (4h)**

- Learn SLIs, SLOs, error budgets, burn-rate alerts, symptom-based alerting,
  incident roles, postmortems, capacity, and backup/restore.

**Build (7h)**

- Define initial availability and latency SLOs for login and booking.
- Add actionable alerts and a runbook for each.
- Run a game day: slow database, auth outage, or failed deployment.

**Verify and document (2h)**

- Measure detection and recovery time and restore a local database backup.
- Demo Month 5 using the dashboard and one incident.

**Quality gate:** alerts map to user impact and an action, recovery is practiced,
and the postmortem identifies system improvements without blame.

**Evidence:** SLOs, alerts, runbooks, backup proof, and postmortem.

---

## Phase 6 — Architecture capstone and promotion portfolio (Weeks 21–24)

### Week 21 — Requirements and quality attributes

**Study (4h)**

- Learn stakeholder analysis, quality-attribute scenarios, architectural drivers,
  constraints, and Architecture Tradeoff Analysis Method concepts.

**Build (7h)**

- Define Didar’s key use cases and stakeholders.
- Write measurable NFRs for availability, latency, scalability, security,
  privacy, recovery, retention, operability, and cost.
- Prioritize conflicting attributes with explicit business reasoning.

**Verify and document (2h)**

- Review each NFR for stimulus, environment, response, and measure.

**Quality gate:** NFRs are measurable and drive decisions rather than repeating
generic goals such as “fast” or “scalable.”

**Evidence:** `docs/architecture/nfrs.md` and stakeholder review notes.

### Week 22 — Capacity, data, and security

**Study (4h)**

- Learn workload modeling, Little’s Law basics, bottleneck estimation, partitioning,
  caching trade-offs, multi-tenancy, STRIDE, privacy, and disaster recovery.

**Build (7h)**

- Estimate normal and peak traffic, storage growth, bandwidth, and database load.
- Review white-label/multi-tenant behavior under `backend/internal/brand/`.
- Update the threat model and create a risk register with owners and priorities.

**Verify and document (2h)**

- Show formulas, assumptions, sensitivity, and where load tests are needed.

**Quality gate:** estimates use units and ranges, tenant isolation risks are
addressed, and mitigations are prioritized by likelihood and impact.

**Evidence:** capacity model, threat model, and risk register.

### Week 23 — Evolutionary target architecture

**Study (4h)**

- Learn fitness functions, strangler migrations, data migration, platform
  boundaries, build-versus-buy, and cost-aware architecture.

**Build (7h)**

- Design three stages: current, next constraint-driven step, and growth stage.
- Decide which modules stay together and what evidence would justify extraction.
- Add C4 component/deployment views and ADRs for major choices.

**Verify and document (2h)**

- Test the design against outage, traffic spike, secret compromise, and team-growth
  scenarios. Include operational and migration cost.

**Quality gate:** the architecture evolves from measured needs, avoids a
microservices-by-default answer, and includes rollback/migration paths.

**Evidence:** capstone architecture document, diagrams, and ADR set.

### Week 24 — Design review and promotion package

**Study (4h)**

- Learn concise architecture communication, design-review facilitation, handling
  objections, and separating facts from assumptions.

**Build (7h)**

- Prepare a 15-minute review: context, drivers, current evidence, options,
  decisions, risks, operations, migration, and open questions.
- Rehearse with a reviewer and revise weak decisions.
- Create the portfolio index below with links to actual artifacts.

**Verify and document (2h)**

- Record the presentation or deliver it live.
- Assess every readiness criterion honestly and create a 90-day follow-up plan.

**Quality gate:** a reviewer can challenge the design, trace decisions to
requirements and evidence, and understand both risks and next steps.

**Evidence:** presentation, review feedback, portfolio index, and follow-up plan.

---

## Monthly review gates

At Weeks 4, 8, 12, 16, 20, and 24:

1. Run the complete test suite from a clean checkout.
2. Demo the month’s working behavior.
3. Explain one failure mode and recovery path without notes.
4. Review code quality, tests, security, and operational impact.
5. Update assumptions and the next month’s scope from evidence.
6. Ask a peer or AI pair to challenge one important decision.

## Promotion readiness criteria

Use evidence, not elapsed time:

- **Go:** independently design, test, profile, and operate an idiomatic service.
- **Microservices:** justify boundaries and communication styles, including when
  not to create a service.
- **DevOps:** own build-to-production flow, rollback, infrastructure,
  observability, SLOs, and incident diagnosis.
- **Architecture:** turn business needs into measurable quality attributes,
  quantify assumptions, compare options, document consequences, and lead review.
- **Leadership:** make work reviewable, improve team standards, mentor others,
  communicate risk early, and follow through after incidents.

A realistic first promotion target may be backend ownership or technical-lead
scope. Use the architect portfolio to demonstrate impact; do not rely on a title
earned only through coursework.

## Portfolio and promotion evidence

Replace each placeholder with a repository link as work is completed.

| Capability | Required evidence | Planned repository artifact |
|---|---|---|
| Go HTTP and errors | API change plus table-driven tests | `[ ]` [Auth package](../services/auth/internal/auth/) |
| Data design | Migration, query plan, ownership diagram | `[ ]` [Main migrations](../backend/migrations/) and `docs/architecture/data-ownership.md` |
| Concurrency | Race-safe implementation and cancellation test | `[ ]` [Booking package](../backend/internal/booking/) |
| API quality | OpenAPI, consistent errors, contract test | `[ ]` `docs/api/openapi.yaml` and [platform helpers](../pkg/platform/) |
| Testing | Unit, handler, repository, and critical-path tests | `[ ]` Test files beside their Go packages |
| Service boundaries | C4 diagrams and auth-split ADR | `[ ]` `docs/architecture/` and `docs/adr/0001-auth-split.md` |
| Reliability | Timeout policy and failure-injection result | `[ ]` [Auth client](../backend/internal/authclient/) and `docs/reliability/` |
| Async design | Notification ADR, idempotency, recovery test | `[ ]` `docs/adr/0002-booking-notifications.md` and the chosen implementation package |
| Security | Threat model, rate limit, risk register | `[ ]` `docs/security/threat-model.md` and `docs/security/risk-register.md` |
| Containers | Hardened images and scan result | `[ ]` [Backend Dockerfile](../backend/Dockerfile), [auth Dockerfile](../services/auth/Dockerfile), and `docs/security/image-scan.md` |
| CI/CD | Required pipeline, smoke test, rollback proof | `[ ]` [.github workflows](../.github/workflows/) and `docs/runbooks/rollback.md` |
| Infrastructure | Reproducible IaC environment | `[ ]` `infra/` and `docs/runbooks/environment-bootstrap.md` |
| Kubernetes | Repeatable manifests and failure drills | `[ ]` `k8s/` and `docs/runbooks/kubernetes-failure-drills.md` |
| Observability | Correlated logs, metrics, trace, dashboard | `[ ]` `observability/` and `docs/architecture/observability.md` |
| Operations | SLOs, alerts, runbooks, postmortem | `[ ]` `docs/slo/`, `docs/runbooks/`, and `docs/postmortems/` |
| Architecture | NFRs, capacity model, staged target design | `[ ]` `docs/architecture/nfrs.md`, `capacity.md`, and `target-architecture.md` |
| Communication | Recorded design review and feedback | `[ ]` `docs/reviews/final-design-review.md` |

## 24-week progress tracker

```text
[ ] Week 01 — Auth request flow and first tested endpoint
[ ] Week 02 — Go types, errors, context, and logging
[ ] Week 03 — Safe API filtering, sorting, and pagination
[ ] Week 04 — Booking invariants and authorization
[ ] Week 05 — Database design and measured migration
[ ] Week 06 — Transactions, concurrency, and cancellation
[ ] Week 07 — Lifecycle, health, errors, and OpenAPI
[ ] Week 08 — Layered tests and performance profile
[ ] Week 09 — Boundaries, C4 diagrams, and first ADR
[ ] Week 10 — Synchronous resilience and failure injection
[ ] Week 11 — Notification consistency and idempotency
[ ] Week 12 — Contracts, security, and rate limiting
[ ] Week 13 — Containers and deployment-model trade-offs
[ ] Week 14 — Backend continuous integration
[ ] Week 15 — Delivery, smoke tests, and rollback
[ ] Week 16 — Reproducible infrastructure as code
[ ] Week 17 — Local Kubernetes deployment
[ ] Week 18 — Kubernetes reliability and safe migrations
[ ] Week 19 — Structured logs, metrics, and traces
[ ] Week 20 — SLOs, alerts, backup, and incident drill
[ ] Week 21 — Stakeholders and measurable NFRs
[ ] Week 22 — Capacity, tenancy, and threat model
[ ] Week 23 — Evolutionary target architecture
[ ] Week 24 — Design review and promotion portfolio
```

## Decision journal template

Create decisions under `docs/adr/` using this format:

```markdown
# ADR NNNN: Decision title

Status: proposed | accepted | superseded
Date:
Review date:

## Context and constraints
## Decision drivers
## Options considered
## Decision
## Positive and negative consequences
## Validation and rollback
```

## Working with your AI pair

When stuck:

```text
I'm on Week X of docs/learning-plan.md.
I read: [paths]
I tried: [what I did]
I expected: [result]
I observed: [error, logs, test, or metric]
My current hypothesis: [specific cause]
Help me test the hypothesis without replacing the learning task.
```

When a deliverable is complete:

```text
I finished Week X. Review the changed files for correctness, tests, security,
operability, and unnecessary complexity. Then ask me three questions that verify
I understand the design and trade-offs.
```

Useful primary references:

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [The Go Blog](https://go.dev/blog/)
- [Docker documentation](https://docs.docker.com/)
- [Kubernetes documentation](https://kubernetes.io/docs/)
- [OpenTelemetry documentation](https://opentelemetry.io/docs/)
- [OWASP API Security](https://owasp.org/www-project-api-security/)
- [C4 model](https://c4model.com/)
