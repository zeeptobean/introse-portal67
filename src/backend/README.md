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

## Backend layout

- `cmd/api/main.go`: loads configuration, starts the database pool, and manages
  the HTTP server lifecycle.
- `internal/config/config.go`: loads `.env`, reads settings, and applies defaults.
- `internal/app/routes.go`: constructs feature stores/services/handlers and
  registers all endpoints and their intended access boundaries. It passes the
  shared session manager to authentication and wraps the router with SCS.
- `internal/app/sessions.go`: configures the SCS PostgreSQL store and cookies.
- `internal/httpx/json.go`: shared JSON responses and the placeholder response.
- `internal/auth`: account types, login/logout/current-user handlers, and
  authentication/role middleware placeholders.
- `internal/academics`: course sections, enrollment, scheduling, grades, and GPA.
- `internal/tuition`: semester charges and simulated payment records.
- `internal/documents`: student requests, staff decisions, and request history.
- `db/migrations` and `db/seeds`: versioned schema changes and separate demo data.

Each feature has a `handler.go`, `service.go`, and `store.go`. Models belong in
`model.go`, or `user.go` for authentication. Constructors wire concrete
dependencies; feature operations and queries are left as TODOs.

Handlers will decode HTTP input and format responses. Services will enforce
business rules and own transactions. Stores will run parameterized SQL. Related
queries within an operation must share the service's `pgx.Tx`.

## Scaffold endpoints

Only health and readiness are implemented. Every feature endpoint currently
returns HTTP `501` with `{"error":"not_implemented"}`. The authentication and
role middleware also return `501` without invoking the next handler, so their
placeholders do not grant access.

- Authentication: `POST /api/auth/login`, `POST /api/auth/logout`, `GET /api/me`.
- Academics: `GET /api/sections`, `POST /api/enrollments`,
  `DELETE /api/enrollments/{id}`, `GET /api/me/timetable`, `GET /api/me/exams`,
  `GET /api/me/results`, `PUT /api/sections/{id}/grades`.
- Tuition: `GET /api/me/tuition`, `POST /api/staff/payments`,
  `PATCH /api/staff/payments/{id}`.
- Documents: `POST /api/document-requests`, `GET /api/document-requests`,
  `PATCH /api/document-requests/{id}`.

Implement authentication and its middleware before enabling the protected
workflows. Role checks do not replace record ownership or lecturer-assignment
checks inside the feature services. Feature model files are placeholders except
for the shared account identity type.

Build all packages from this directory with `go build ./...`. No feature test
suite, additional runtime dependency, automatic migration, or seed execution is
part of this scaffold.
