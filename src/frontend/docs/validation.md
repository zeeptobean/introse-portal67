# Prototype validation

Validated on 26 October 2026 against the repository's Go API and an isolated PostgreSQL 16.4 database. All eight migrations and four seed files were applied. The user's normal database and environment files were not changed.

## Automated checks

The migration adds Zod contract tests and a reproducible browser suite. Run `pnpm check`, `pnpm lint`, `pnpm test`, `pnpm test:browser`, `pnpm build`, and `go test ./...` from their respective frontend/backend directories. The 22 unit tests cover term selection, timetable lanes, zero versus blank scores, Unicode limits, VND and ID validation, cookie isolation/relay, backend response validation and uncertain writes. Go tests include trusted-proxy forwarding and spoof rejection.

The browser suite uses a running seeded Go API and SvelteKit server, with no mock success data. It checks ordinary HTML forms with JavaScript disabled, cookie rotation/logout, unauthorized POST actions, the enhanced cross-role workflows, no browser `/api` requests, narrow layouts, and a deliberately aborted form response that preserves its draft and disables immediate resubmission.

Final migration results: Svelte/TypeScript and lint passed with zero warnings/errors; all 22 unit tests and all 3 browser tests passed; the Node production build and Go test suite passed. The browser suite was run against the built Node server at `http://localhost:4173` with Microsoft Edge headless, the real Go API, and the isolated sample database. The ordinary development server remains `http://localhost:5173`.

## Browser journeys exercised

| Journey                                              | Observed outcome                                                                                                                                |
| ---------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| Sign in and sign out as student, lecturer, staff     | Correct role home and navigation; Go session is relayed through the SvelteKit server                                                            |
| Preview and register CSC15001, then cancel           | Confirmed registration reference, seats and 10-credit summary updated; cancellation released the seat and retained history                      |
| Student document submission → staff approval         | Durable request receipt; final comment, actor, and both history events returned by the API                                                      |
| Staff lookup and simulated payment                   | Correct student ID retained; 400,000 VND payment changed balance from 5,750,000 to 5,350,000 VND                                                |
| Correct payment to 500,000 VND, then void            | Review/confirmed balances 5,250,000 then 5,750,000 VND; voided row remained excluded from totals                                                |
| Add 100,000 VND charge                               | Saved charge returned to the confirmed student ledger                                                                                           |
| Lecturer invalid grade, zero grade, then 8.50        | Invalid 12 blocked with a focused error summary; zero saved as a grade; 8.50 appeared in student Results as A / 3.50 with 4 of 7 credits graded |
| 320 px tuition ledger, payment review, and gradebook | No page-level horizontal overflow; identity, input, review, and save controls remained available                                                |
| Mobile navigation                                    | Dialog opened on Close menu; Escape closed it and restored focus to Menu                                                                        |

The original browser observations above were recorded before migration. The automated migration suite repeats the document/decision, registration/cancellation, payment/correction/void, charge creation and grade-to-results journeys against the server boundary. Grade values toggle between 8.25 and 8.50 across repeated smoke runs, and each request/payment/charge has a distinct test reference. Tests operate only on the isolated synthetic database.

The grade editor waits for hydration before enabling its custom inputs: an early typed value must not be replaced by the initial saved score. Ordinary login and read pages remain operable with JavaScript disabled. Node production checks use a fixed build-time public origin while preserving SvelteKit's CSRF checks.

## Limits to keep visible

- This is a sample-data course prototype; no official-system connection, money transfer, or document issuance.
- GPA descriptions reflect the seeded “Sample 10-point scale.” The API does not expose rule metadata. If backend bands change, update the explanation; a versioned rule endpoint remains a follow-up (O09).
- Grade and tuition preflight reads do not provide atomic version checks. Concurrent edit protection remains backend work (O06). Coordinate editors during demonstrations.
- New payment/charge/request writes lack durable idempotency keys (O10). An uncertain response blocks immediate resubmission and directs the user to current records. Do not infer a successful creation merely from a similar row.
- Unsaved edits are guarded during navigation/reload/sign-out. Session expiry redirects to sign-in and discards drafts; same-account draft restoration from U03 is not implemented.
- The English copy follows the supplied reference. Vietnamese-first localization remains a team policy choice (O03); names and text retain Unicode.
- Small prototype lists are shown in full and filtered locally. Server pagination, user research, load/concurrency testing, and a complete accessibility audit are not claimed.
- The Node adapter is configured. Deployment still needs a correct public origin, HTTPS and trusted ingress settings; see the README. There is no development `/api` proxy.
