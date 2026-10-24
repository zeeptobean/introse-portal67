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
