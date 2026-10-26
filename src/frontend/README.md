# Portal67 frontend

SvelteKit 3, Svelte 5, TypeScript, Vite+, Zod 4, and Superforms 3. Pages render on the server, and all Go API calls originate from server-only modules. The [design package](../../docs/design/portal67-experience-design.md) defines the product flows and visual tokens. Start with the [architecture guide](docs/architecture.md) when contributing.

## Run locally

Prerequisites: Node matching `package.json`, pnpm 10, Go matching `../backend/go.mod`, PostgreSQL, and [goose](https://github.com/pressly/goose) for migrations. The repository's Compose file provides PostgreSQL if Docker is available.

From the repository root:

```sh
docker compose up -d db
```

From `src/backend`, copy `.env.example` to `.env`. Set `DATABASE_URL` to your local database. For the default Compose database, run these commands once:

```sh
goose -dir db/migrations postgres 'postgres://portal67:portal67_dev@127.0.0.1:5432/portal67?sslmode=disable' up
goose -dir db/seeds -no-versioning postgres 'postgres://portal67:portal67_dev@127.0.0.1:5432/portal67?sslmode=disable' up
go run ./cmd/api
```

The seed files contain synthetic demo accounts and their shared password. See [`00001_demo_users.sql`](../backend/db/seeds/00001_demo_users.sql). Use a separate sample database; the prototype writes actual enrollments, grades, requests, charges, and payments to that database.

In another terminal, from `src/frontend`:

```sh
pnpm install --frozen-lockfile
pnpm dev
```

The frontend runs at http://localhost:5173. Copy `.env.example` to `.env` to change the private `BACKEND_URL` (default `http://127.0.0.1:8080`). There is no browser `/api` proxy. Browser forms submit to SvelteKit actions, which validate input and call Go. `BACKEND_URL` must be an HTTP(S) origin without credentials, a path, or a query.

The Go session remains authoritative. SvelteKit forwards only its `session` cookie to Go and relays rotation/expiry to the browser as HttpOnly, Secure, SameSite=Strict. Use `localhost` for local HTTP demonstrations; deployment requires HTTPS. Existing sample data and database migrations are unchanged.

For per-client login rate limits, configure Go's `TRUSTED_PROXY_CIDRS` with only the SvelteKit server's immediate address. The example uses `127.0.0.1/32,::1/128` for same-machine development. Leave it empty for direct API use. SvelteKit overwrites `X-Forwarded-For`; Go ignores this header from untrusted peers or when it contains multiple addresses.

```sh
pnpm check       # Svelte and TypeScript diagnostics
pnpm lint        # Vite+ formatting, lint and type checks
pnpm format      # Vite+ formatting, including Svelte files
pnpm test        # Domain and API recovery contract tests
pnpm test:browser # Browser smoke tests; running seeded Go + frontend required
pnpm build       # Production build
pnpm preview     # Preview the production build locally
```

## Production Node server

The project uses `adapter-node`. Set the canonical public `ORIGIN` **when building** (SvelteKit 3 uses `paths.origin`, not adapter-node's former runtime `ORIGIN` setting). For example, in PowerShell:

```powershell
$env:ORIGIN='https://portal.example.com'
pnpm build
$env:BACKEND_URL='http://127.0.0.1:8080'
$env:HOST='127.0.0.1'
$env:PORT='3000'
node build/index.js
```

Put an HTTPS reverse proxy in front of the Node server and keep Go reachable from Node. The browser no longer needs access to Go. Production Node does not load `.env` automatically; export runtime variables or use Node's `--env-file=.env` option.

Without a fixed build-time `ORIGIN`, adapter-node assumes HTTPS and derives the host from the request. Behind a trusted ingress, configure `PROTOCOL_HEADER`, `HOST_HEADER`, and any client-address settings according to the [Node adapter documentation](https://svelte.dev/docs/kit/adapter-node). Never disable CSRF checks to work around an origin mismatch.

For a local production-build smoke test, build with `ORIGIN=http://localhost:4173`, then run Node on port 4173. The browser suite accepts `PLAYWRIGHT_BASE_URL` and optional `PLAYWRIGHT_CHANNEL=msedge`. Otherwise install Chromium once with `pnpm exec playwright install chromium`. Browser tests write synthetic requests, payments, charges, and grades; run them only against an isolated seeded database.

## Workspaces

- Students: Today, registration and section previews, classes/grid and exams, estimated GPA, tuition ledger, and document submission/history.
- Lecturers: assigned sections, validated batch grade entry, changed/hidden row counts, and saved grade history.
- Staff: request queue and final decisions; student lookup; charge creation; payment creation, correction, voiding, and history.

Accounts have one role. The API enforces authorization; frontend role guards provide navigation and recovery. See [validation evidence and remaining limitations](docs/validation.md).
