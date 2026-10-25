-- +goose Up
-- Append-only history of every score saved for an enrollment. The grades row
-- holds the current score; its graded_by and updated_at match the newest
-- revision. Grades that existed before this table get one revision each.
CREATE TABLE grade_revisions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    enrollment_id BIGINT NOT NULL REFERENCES enrollments (id),
    score NUMERIC(4,2) NOT NULL CHECK (score BETWEEN 0 AND 10),
    graded_by BIGINT NOT NULL REFERENCES users (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX grade_revisions_history_idx ON grade_revisions (enrollment_id, created_at, id);

INSERT INTO grade_revisions (enrollment_id, score, graded_by, created_at)
SELECT enrollment_id, score, graded_by, updated_at
FROM grades
ORDER BY updated_at, enrollment_id;

-- +goose Down
DROP TABLE grade_revisions;
