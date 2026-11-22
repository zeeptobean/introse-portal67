# Portal67 backend

Go HTTP API using `net/http` and a PostgreSQL connection pool from `pgx`.

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

## Prototype v2 enrollment events

Set `GOOSE_DBSTRING` to the same database as `DATABASE_URL` and run `goose -dir ./db/migrations up` from this directory before starting the v2 API. Migration `00009_registration_outbox.sql` adds the outbox and preserves existing enrollments. Stop earlier API writers before switching over if every enrollment change needs an event. When rolling back the application, keep the migration applied to preserve queued events.

Set `REDIS_URL` in the ignored `.env` to enable publication. The hosted instance is `redis.ancon-pickerel.ts.net:6379`. Supply an authenticated Redis URL and URL-encode special characters in the credential. Set `REDIS_KEY_PREFIX` to a distinct value for each environment sharing Redis. The default prefix is `p67`.

The API commits registrations and cancellations together with their events in PostgreSQL. The relay publishes them to `p67:ev:registration` by default. Unsent events stay queued if Redis is unavailable or unset. A failure after publication can cause redelivery with the same event ID. Existing enrollment rows receive a version during migration. Events start with writes made by the v2 API.

To run the registration integration tests, set `PORTAL67_TEST_DATABASE_URL` to a development PostgreSQL database where the test role can create schemas, then run `go test ./internal/academics -count=1`. The tests create and remove an isolated schema.
