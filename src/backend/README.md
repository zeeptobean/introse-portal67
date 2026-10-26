# Portal67 backend

Go HTTP API using `net/http` and a PostgreSQL connection pool from `pgx`.

Copy `.env.example` to `.env`, then run from this directory:

```sh
go mod download
go run ./cmd/api
```

The API listens at http://127.0.0.1:8080 by default.

- `GET /api/health`: process liveness.
- `GET /api/ready`: pings PostgreSQL with a three-second timeout.

The SvelteKit server now calls the API on behalf of the browser and relays the existing SCS session. For login rate limiting behind it, set `TRUSTED_PROXY_CIDRS` to the immediate SvelteKit server addresses (for example `127.0.0.1/32,::1/128` on one machine). Go uses a single `X-Forwarded-For` address only from those peers. Leave the setting empty for direct connections, and never configure an unrestricted trust range. Frontend HTTPS/origin and private-backend setup is documented in [the frontend README](../frontend/README.md).
