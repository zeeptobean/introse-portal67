-- +goose Up
-- Every registration gets its own row, so cancelled rows stay as history and a
-- re-registration inserts a new row. A student still has at most one active
-- enrollment per section.
ALTER TABLE enrollments DROP CONSTRAINT enrollments_student_id_section_id_key;

CREATE UNIQUE INDEX enrollments_one_active_idx
    ON enrollments (student_id, section_id) WHERE status = 'enrolled';
CREATE INDEX enrollments_student_idx ON enrollments (student_id, section_id);

-- +goose Down
-- Fails while a student has more than one row for the same section; remove
-- the extra cancelled rows by hand before rolling back.
DROP INDEX enrollments_student_idx;
DROP INDEX enrollments_one_active_idx;
ALTER TABLE enrollments
    ADD CONSTRAINT enrollments_student_id_section_id_key UNIQUE (student_id, section_id);
