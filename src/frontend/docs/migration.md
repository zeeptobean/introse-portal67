# Server and schema migration

The original browser API modules were replaced by Zod response/input contracts, Superforms, server loads and actions. Existing Go routes, permissions, database schemas and sample records remain the source of truth. No data migration or reset is required.

## Completed changes

- Zod schemas own response and form types, including string bigint IDs, nullable grades, Go nil slices, Unicode text limits and campus-time payment input.
- All Go HTTP communication is in `src/lib/server/`, behind a private validated `BACKEND_URL`. The Vite `/api` proxy and browser API helpers are removed.
- SvelteKit relays only Go's session cookie. Login, logout and protected actions are checked on the server; session tokens never enter returned page data.
- Superforms manages input validation, errors and submission state. Decision/payment review remains explicit; the gradebook retains changed and filtered-out rows.
- Semester preferences use non-sensitive per-account/task cookies. Draft state uses component context so server-rendered requests do not share mutable state.
- `adapter-node` produces a runnable server. `ORIGIN` sets SvelteKit 3's build-time canonical origin; runtime `BACKEND_URL` controls Go connectivity.
- Optional Go trusted-proxy CIDRs retain per-client login rate limiting without accepting arbitrary forwarded headers.

## Compatibility and rollback

A pre-migration frontend archive was saved outside this worktree before editing. Its local validation-session path was `C:\Users\BBQ\AppData\Local\Temp\portal67-before-server-migration-20261026-152206.tar`. This temporary artifact is not a committed backup or a deployment mechanism.

To roll back, restore the prior frontend source, manifests, lockfile and Vite configuration as a consistent set. Keep the original browser client's `/api` routing available if serving that version. The Go API remains compatible, and the optional trusted-proxy setting can be left empty for direct API requests. Do not roll back or delete database records to undo frontend code.

## Remaining backend decisions

The frontend cannot provide durable creation idempotency, atomic edit versioning, or a versioned GPA-rule description without new Go contracts. An uncertain new request/payment/charge is therefore never automatically retried. Existing-record actions can verify exact current records. Same-account draft restoration after session expiry remains unimplemented.

See [architecture](architecture.md), [validation](validation.md), and the [setup/deployment instructions](../README.md) for the maintained workflow.
