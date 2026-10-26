# Frontend structure

The application follows the [page and flow specifications](../../../docs/design/portal67-pages-and-flows.md). SvelteKit owns page loading and form handling; Go owns authentication, permissions, business rules, and transactions.

## Ownership

| Location                                                       | Responsibility                                                                                                    |
| -------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| `src/lib/schemas/{shared,academics,documents,tuition,auth}.ts` | Zod response/input contracts and inferred types; IDs remain strings                                               |
| `src/env.ts`                                                   | Validated private runtime `BACKEND_URL`                                                                           |
| `src/hooks.server.ts`                                          | Construct a separate API client for each incoming request; private no-store responses                             |
| `src/lib/server/client.ts`                                     | Trusted backend origin, session-cookie forwarding/relay, response parsing, timeout and uncertainty classification |
| `src/lib/server/api.ts`                                        | Typed Go endpoint calls; server-only import boundary                                                              |
| `src/lib/server/auth.ts`                                       | Per-request identity lookup, protected loads/actions, safe role-scoped login destination                          |
| `src/lib/server/forms.ts`                                      | Translate API failures into safe form messages or page errors; verify uncertain writes against current records    |
| `src/lib/server/{registration,grades,tuition}.ts`              | Domain action orchestration, authoritative read-back and tuition preflight checks                                 |
| `src/lib/server/terms.ts`                                      | Non-sensitive semester preference cookies, partitioned by account and task                                        |
| `src/lib/form.ts`                                              | Shared Superforms defaults; no blind write replay                                                                 |
| `src/lib/draft.svelte.ts`                                      | Component-context draft guard; isolated during SSR                                                                |
| `src/lib/components/`                                          | Display and interaction: review steps, filtered roster, accessible errors, agendas and ledgers                    |
| `src/lib/styles.css`                                           | Design tokens, element defaults and primitives shared by several components                                       |
| `src/routes/**/+page.server.ts`                                | Authoritative page data and form actions                                                                          |
| `src/routes/**/+layout.server.ts`                              | Identity, semester data and role-specific navigation checks                                                       |

Import application modules through `#lib/...` with their file extension, and import types from the schema module that owns them. `+page.svelte` files compose views and consume server data. Components submit to SvelteKit actions using Superforms `enhance`; they do not import Go HTTP helpers or store credentials. Read-only list filtering and timetable interactions remain in the browser.

## Add a workflow

1. Check the actual Go model/handler and design scope. Define the Zod response and input schema in the appropriate domain module, and infer its TypeScript type.
2. Add a typed endpoint call to `server/api.ts`. All responses pass through Zod before they reach page data. Declare collections with `list()`: Go nil or omitted slices become empty arrays there, so views never need a fallback.
3. Add a `+page.server.ts` load. Wait for the protected parent before reading private records; use `readPage` for safe missing/access/service states. Initialize forms with `superValidate` and `zod4`.
4. Add an action that calls `requireUser(event, role)` independently of its layout, validates submitted input, calls Go, and returns a form or confirmed redirect. Layout checks alone do not protect POST actions.
5. Compose the form with `portalForm`, the Zod client adapter, field labels/errors and the required review step. Superforms owns fields, validation and submission state. Keep domain-specific comparison and changed-row logic local.
6. After a successful mutation, refresh authoritative data or redirect to a newly loaded receipt/ledger. Do not replay writes on transport failure; use `verifyUncertain` to check current records for the intended change before reporting the write as unconfirmed.
7. Style a component in its own `<style>` block, with media queries beside the rules they adjust. Add to `styles.css` only what several components share.
8. Run the checks in the README and exercise the real seeded backend journey.

Server loads use `portal:session`, `portal:academics`, `portal:documents`, `portal:grades`, and `portal:tuition` dependencies. Enhanced successful forms invalidate page data; redirects load current records. Grade edits include only changed rows, including rows hidden by search; untouched blanks never become zero grades.

## Sessions and state

Go's PostgreSQL-backed SCS session is the only authentication session. Node forwards only the `session` cookie to the configured private origin. It never follows backend redirects with that cookie. Only the backend session cookie is relayed to the frontend host, preserving expiry/clearing and enforcing HttpOnly, Secure, SameSite=Strict. Tokens never appear in page data. Login passwords are cleared from failed server form responses.

Identity caching lives in `event.locals`, not a module singleton. Editable draft state lives in Svelte component context, not shared server state. Non-sensitive semester IDs live in HttpOnly cookies and are cleared at sign-out; private records and drafts never enter cookies or browser storage. Session expiry sends the user to sign-in and does not replay their pending action.

Go's optional trusted-proxy configuration preserves per-client login limits through Node. Do not trust arbitrary browser forwarding headers or all network ranges. SvelteKit's action origin checks stay enabled. Configure the actual public origin and trusted ingress correctly when deploying.

## Recovery and limits

Transport loss after a write is an unknown outcome. Registration, cancellation and grade actions try exact authoritative read-back. Decisions load the current terminal decision. New requests, charges and payments cannot be safely identified by similar values; an uncertain result blocks immediate repeat and links to current records. A preflight tuition read detects some stale edits but does not provide atomic concurrency protection.

Zod mirrors Go's input and response contracts; it does not replace backend permissions or business validation. GPA and ledger totals come from Go. Current sample grading policy, concurrent edit protection, and durable creation idempotency remain the backend decisions recorded in the design package.

## Deployment and rollback

The frontend requires Node (`adapter-node`), rather than static hosting. The private backend URL is a runtime setting. The public origin is a SvelteKit 3 build-time `paths.origin` value when `ORIGIN` is supplied. See the README for local and HTTPS deployment commands.

The migration changes no database schemas or existing API routes. The previous frontend was archived outside the checkout before editing. Rollback means restoring that frontend snapshot and its dependencies/configuration, while leaving database records and Go's established authorization intact. No reverse database migration is needed. Old browser clients still require their original same-origin `/api` routing if kept in service during a rollout.
