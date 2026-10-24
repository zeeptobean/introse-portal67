-- +goose Up
CREATE TABLE document_requests (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    student_id BIGINT NOT NULL REFERENCES users (id),
    request_type TEXT NOT NULL CHECK (request_type IN ('enrollment_confirmation', 'transcript', 'other')),
    purpose TEXT NOT NULL CHECK (btrim(purpose) <> ''),
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX document_requests_student_idx ON document_requests (student_id);
CREATE INDEX document_requests_pending_idx ON document_requests (submitted_at) WHERE status = 'pending';

-- Append-only history. Submission records NULL -> pending; a staff decision
-- records pending -> approved/rejected and must carry a comment.
CREATE TABLE document_request_events (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    document_request_id BIGINT NOT NULL REFERENCES document_requests (id),
    actor_id BIGINT NOT NULL REFERENCES users (id),
    from_status TEXT,
    to_status TEXT NOT NULL,
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (
        (from_status IS NULL AND to_status = 'pending')
        OR (from_status IS NOT NULL AND from_status = 'pending' AND to_status IN ('approved', 'rejected'))
    ),
    CHECK (to_status = 'pending' OR (comment IS NOT NULL AND btrim(comment) <> '')),
    UNIQUE (document_request_id, to_status)
);

CREATE INDEX document_request_events_history_idx
    ON document_request_events (document_request_id, created_at, id);
CREATE UNIQUE INDEX document_request_events_one_decision_idx
    ON document_request_events (document_request_id) WHERE from_status = 'pending';

-- +goose Down
DROP TABLE document_request_events;
DROP TABLE document_requests;
