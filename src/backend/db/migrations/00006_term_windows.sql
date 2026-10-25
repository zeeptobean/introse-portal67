-- +goose Up
-- Students may register for or cancel sections only while
-- registration_opens_at <= now() < registration_closes_at. Lecturers may enter
-- or change grades only before grading_closes_at; after that grades are final.
ALTER TABLE terms
    ADD COLUMN registration_opens_at TIMESTAMPTZ,
    ADD COLUMN registration_closes_at TIMESTAMPTZ,
    ADD COLUMN grading_closes_at TIMESTAMPTZ;

-- Existing terms keep their previous behaviour as closely as possible:
-- registration stays open until the end of the term's last day (campus time),
-- and grading closes 30 days after the term ends.
UPDATE terms SET
    registration_opens_at = (starts_on - 90)::timestamp AT TIME ZONE 'Asia/Ho_Chi_Minh',
    registration_closes_at = (ends_on + 1)::timestamp AT TIME ZONE 'Asia/Ho_Chi_Minh',
    grading_closes_at = (ends_on + 31)::timestamp AT TIME ZONE 'Asia/Ho_Chi_Minh';

ALTER TABLE terms
    ALTER COLUMN registration_opens_at SET NOT NULL,
    ALTER COLUMN registration_closes_at SET NOT NULL,
    ALTER COLUMN grading_closes_at SET NOT NULL,
    ADD CONSTRAINT terms_registration_window_check CHECK (registration_closes_at > registration_opens_at),
    ADD CONSTRAINT terms_grading_window_check CHECK (grading_closes_at > registration_closes_at);

-- +goose Down
ALTER TABLE terms
    DROP COLUMN grading_closes_at,
    DROP COLUMN registration_closes_at,
    DROP COLUMN registration_opens_at;
