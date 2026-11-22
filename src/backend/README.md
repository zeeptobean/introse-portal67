# Portal67 backend

## Run prototype v2

The gateway exposes the existing Portal67 API. Identity owns sessions and accounts, academic owns grades and the course catalog, registration owns enrollment and seat counts, affairs owns tuition and document requests, and notification exposes `GET /api/me/notifications`. Each service runs with its own PostgreSQL role and schema. Services exchange changes through Redis Streams.

1. Copy the repository-root `.env.example` to `.env`. Run `go run ./cmd/portalctl keygen` from `src/backend` and put the two generated values in the root `.env`. The gateway gets the signing key. Services get the corresponding public key.
2. Set `REDIS_URL` in that file to the authenticated hosted endpoint `redis.ancon-pickerel.ts.net:6379`, with the password URL-encoded. Leave it blank to use the private Compose Redis. Choose a separate `REDIS_KEY_PREFIX` for each stack on shared Redis.
3. From the repository root, run `docker compose --profile prototype-v2 up --build -d`. The migration job creates the five service roles and schemas. The gateway publishes only on loopback. Check `http://127.0.0.1:8080/api/ready` before starting the frontend.
4. Set the SvelteKit server's `BACKEND_URL` to `http://127.0.0.1:8080` and start it normally. The sample login accounts are defined in [the demo seed](db/seeds/00001_demo_users.sql).

Use PostgreSQL 18 for the split. Compose supplies demo database credentials for local use. For a shared deployment, configure the admin URL and distinct passwords for the five service roles in the ignored environment file. Set `GATEWAY_TRUSTED_PROXY_CIDRS` to the immediate BFF addresses if the gateway should accept the BFF's forwarded client IP.

### Run the services directly

From `src/backend`, run `go run ./cmd/portalctl migrate` with `PORTAL67_ADMIN_DATABASE_URL` and the five `<SERVICE>_DATABASE_URL` values. Set `PORTAL67_SEED_DEMO=1` for an empty demo database. Start `go run ./cmd/identity-service`, `./cmd/academic-service`, `./cmd/affairs-service`, `./cmd/notification-service`, and `./cmd/registration-service` in separate processes. For each process, set `DATABASE_URL` to its own role URL and set `REDIS_URL`, `REDIS_KEY_PREFIX`, and `IDENTITY_PUBLIC_KEYS`. The services bind to loopback ports 8081 through 8085.

Start `go run ./cmd/gateway` with `IDENTITY_SIGNING_KEY`, the same Redis settings, and `IDENTITY_URL`, `ACADEMIC_URL`, `AFFAIRS_URL`, `NOTIFICATION_URL`, and `REGISTRATION_URL` pointing to those services. Set `SERVICE_ADDR` per process to override its bind address.

### Cutover and Redis recovery

Stop legacy API writers before the first split migration. The migration copies the current data, including sessions and enrollment history, into the split schemas. The old `public` schema is a pre-cutover checkpoint. Preserve the split schemas and reconcile their newer writes before rolling back to the legacy API.

Outbox rows remain available for replay. The relay retries Redis delivery, and consumers apply an event and its receipt in a single database transaction. Seat accounting stays in PostgreSQL. If a cancellation and grade cross in transit, academic preserves the grade, removes it from active results, and notifies staff for review. The event wire format is defined in the [v1 contract](contracts/events/v1.schema.json).

- Requeue an outbox window with `go run ./cmd/portalctl replay <RFC3339>`, using the owning `SERVICE` and its `DATABASE_URL`.
- Republish current owner state with `go run ./cmd/portalctl republish` for `identity`, `academic`, or `registration`. Use this to seed newly created projections.
- Retry an entry from `<prefix>:dlq:<group>` with `go run ./cmd/portalctl dlq-retry <stream> <group> <entry-id>`, using the environment's Redis settings. This publishes a new Redis entry with the original event ID.

### Verify

Run `go test ./...` from `src/backend`. Set `PORTAL67_TEST_DATABASE_URL` and `PORTAL67_TEST_REDIS_URL` to run the outbox and Redis integration tests. Set `PORTAL67_SPLIT_TEST_ADMIN_URL` to an isolated PostgreSQL 18 cluster to run the split migration, permission, seat-capacity, and grade/cancellation checks. The split test creates and drops a disposable database. Run the frontend Playwright suite against a production SvelteKit build backed by the gateway and an isolated seeded stack.

## Legacy single-process API

The legacy API is available for the pre-cutover schema. It uses `net/http` and a PostgreSQL connection pool from `pgx`.

Copy `.env.example` to `.env`. Once PostgreSQL is running, run these commands from this directory.

```sh
go mod download
goose -dir ./db/migrations up
go run ./cmd/api
```

The API listens at http://127.0.0.1:8080 by default.

- `GET /api/health`: process liveness.
- `GET /api/ready` checks PostgreSQL and, when configured, Redis within a shared three-second timeout.

The SvelteKit server now calls the API on behalf of the browser and relays the existing SCS session. For login rate limiting behind it, set `TRUSTED_PROXY_CIDRS` to the immediate SvelteKit server addresses (for example `127.0.0.1/32,::1/128` on one machine). Go uses a single `X-Forwarded-For` address only from those peers. Leave the setting empty for direct connections, and never configure an unrestricted trust range. Frontend HTTPS/origin and private-backend setup is documented in [the frontend README](../frontend/README.md).

### Legacy enrollment events

Set `GOOSE_DBSTRING` to the same database as `DATABASE_URL` and run `goose -dir ./db/migrations up` from this directory before starting the v2 API. Migration `00009_registration_outbox.sql` adds the outbox and preserves existing enrollments. Stop earlier API writers before switching over if every enrollment change needs an event. When rolling back the application, keep the migration applied to preserve queued events.

Set `REDIS_URL` in the ignored `.env` to enable publication. The hosted instance is `redis.ancon-pickerel.ts.net:6379`. Supply an authenticated Redis URL and URL-encode special characters in the credential. Set `REDIS_KEY_PREFIX` to a distinct value for each environment sharing Redis. The default prefix is `p67`.

The API commits registrations and cancellations together with their events in PostgreSQL. The relay publishes them to `p67:ev:registration` by default. Unsent events stay queued if Redis is unavailable or unset. A failure after publication can cause redelivery with the same event ID. Existing enrollment rows receive a version during migration. Events start with writes made by the v2 API.

To run the registration integration tests, set `PORTAL67_TEST_DATABASE_URL` to a development PostgreSQL database where the test role can create schemas, then run `go test ./internal/academics -count=1`. The tests create and remove an isolated schema.
