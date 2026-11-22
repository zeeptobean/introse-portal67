# AI Usage Log — Prototype v2 enrollment events

* **Member Name:** [Full Name]
* **Student ID:** [MSSV]
* **Role:** Developer
* **Phase Gate:** PA5 - Construction
* **Date:** 2026-11-22

## Session details

### Task

Start the `prototype-v2` implementation using the microservices design and the team's hosted Redis.

* **Tool Used:** Codex
* **Objective:** Add transactional enrollment events and publish them to Redis Streams from the existing Go API.
* **Input Prompt:** Read the microservices design thread, start the new prototype branch, and connect to the shared Redis endpoint.
* **AI Output Summary:** Added the PostgreSQL migration, versioned enrollment events, Redis Streams relay, configuration, and integration tests for seat capacity, transaction rollback, and publication retry.
* **Human Review & Modifications:** Reviewed, the current design are suitable for extendability and since this is prototype, this is all subject to changes
* **Verification Steps:** Ran the Go suite with isolated PostgreSQL and repeated the concurrent one-seat test ten times. Verified that HTTP registration and cancellation produced both events in a temporary stream on the hosted Redis instance. All three Playwright tests passed against the production frontend build.
