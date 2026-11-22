-- +goose Up
-- Take a consistent snapshot after stopping legacy API writers. The public
-- schema remains available as the pre-cutover checkpoint.
LOCK TABLE public.users, public.sessions, public.terms, public.courses, public.sections,
    public.class_meetings, public.exams, public.enrollments, public.grades, public.grading_scales,
    public.grade_bands, public.grade_revisions, public.document_requests, public.document_request_events,
    public.tuition_charges, public.tuition_payments, public.tuition_payment_events, public.outbox IN SHARE MODE;

CREATE SCHEMA identity;
CREATE SCHEMA academic;
CREATE SCHEMA registration;
CREATE SCHEMA affairs;
CREATE SCHEMA notification;
CREATE SCHEMA portal67_events;

-- +goose StatementBegin
DO $storage$
DECLARE
    item record;
    rel text;
    seq text;
    constraint_row record;
BEGIN
    FOR item IN SELECT key AS service, value AS tables FROM jsonb_each('{
        "identity": ["users", "sessions"],
        "academic": ["grading_scales", "grade_bands", "terms", "courses", "sections", "class_meetings",
                     "exams", "enrollments", "grades", "grade_revisions"],
        "registration": ["terms", "courses", "sections", "class_meetings", "exams", "enrollments"],
        "affairs": ["terms", "document_requests", "document_request_events", "tuition_charges",
                     "tuition_payments", "tuition_payment_events"],
        "notification": []
    }') LOOP
        IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'portal67_' || item.service) THEN
            EXECUTE format('CREATE ROLE %I NOLOGIN', 'portal67_' || item.service);
        END IF;
        EXECUTE format('CREATE TABLE %I.outbox (LIKE public.outbox INCLUDING ALL)', item.service);
        EXECUTE format('CREATE TABLE %I.inbox (LIKE public.inbox INCLUDING ALL)', item.service);
        FOR rel IN SELECT jsonb_array_elements_text(item.tables) LOOP
            EXECUTE format('CREATE TABLE %I.%I (LIKE public.%I INCLUDING ALL)', item.service, rel, rel);
            EXECUTE format('INSERT INTO %I.%I OVERRIDING SYSTEM VALUE SELECT * FROM public.%I', item.service, rel, rel);
            SELECT CASE WHEN EXISTS (
                SELECT 1 FROM pg_attribute
                WHERE attrelid = to_regclass(format('%I.%I', item.service, rel)) AND attname = 'id'
            ) THEN pg_get_serial_sequence(format('%I.%I', item.service, rel), 'id') END INTO seq;
            IF seq IS NOT NULL THEN
                EXECUTE format('SELECT setval(%L, coalesce(max(id), 1), count(*) > 0) FROM %I.%I',
                    seq, item.service, rel);
            END IF;
        END LOOP;
    END LOOP;
END
$storage$;
-- +goose StatementEnd

-- Runtime directories keep only the fields other services require.
CREATE TABLE academic.users (
    id bigint PRIMARY KEY, institutional_id text, full_name text NOT NULL, email text NOT NULL,
    role text NOT NULL, is_active boolean NOT NULL, version bigint NOT NULL
);
CREATE TABLE registration.users (LIKE academic.users INCLUDING ALL);
CREATE TABLE affairs.users (LIKE academic.users INCLUDING ALL);
CREATE TABLE notification.users (LIKE academic.users INCLUDING ALL);
INSERT INTO academic.users SELECT id, institutional_id, full_name, email, role, is_active, 1 FROM public.users;
INSERT INTO registration.users SELECT * FROM academic.users;
INSERT INTO affairs.users SELECT * FROM academic.users;
INSERT INTO notification.users SELECT * FROM academic.users;

ALTER TABLE identity.users ADD COLUMN version bigint NOT NULL DEFAULT 1 CHECK (version > 0);
ALTER TABLE academic.terms ADD COLUMN version bigint NOT NULL DEFAULT 1;
ALTER TABLE academic.courses ADD COLUMN version bigint NOT NULL DEFAULT 1;
ALTER TABLE academic.sections ADD COLUMN version bigint NOT NULL DEFAULT 1;
ALTER TABLE academic.class_meetings ADD COLUMN version bigint NOT NULL DEFAULT 1;
ALTER TABLE academic.exams ADD COLUMN version bigint NOT NULL DEFAULT 1;
ALTER TABLE academic.grades ADD COLUMN version bigint NOT NULL DEFAULT 1, ADD COLUMN orphaned_at timestamptz;
ALTER TABLE registration.terms ADD COLUMN version bigint NOT NULL DEFAULT 1;
ALTER TABLE registration.courses ADD COLUMN version bigint NOT NULL DEFAULT 1;
ALTER TABLE registration.sections ADD COLUMN version bigint NOT NULL DEFAULT 1;
ALTER TABLE registration.class_meetings ADD COLUMN version bigint NOT NULL DEFAULT 1;
ALTER TABLE registration.exams ADD COLUMN version bigint NOT NULL DEFAULT 1;
ALTER TABLE affairs.terms ADD COLUMN version bigint NOT NULL DEFAULT 1;

CREATE TABLE registration.grades (enrollment_id bigint PRIMARY KEY REFERENCES registration.enrollments(id), version bigint NOT NULL);
INSERT INTO registration.grades SELECT enrollment_id, 1 FROM public.grades;
CREATE TABLE registration.catalog_tombstones (
    entity text NOT NULL, id bigint NOT NULL, version bigint NOT NULL, PRIMARY KEY (entity, id)
);
CREATE TABLE registration.section_seats (
    section_id bigint PRIMARY KEY REFERENCES registration.sections(id),
    capacity integer NOT NULL CHECK (capacity > 0),
    active_count integer NOT NULL DEFAULT 0 CHECK (active_count >= 0),
    pending_capacity integer CHECK (pending_capacity > 0),
    CHECK (active_count <= capacity)
);
INSERT INTO registration.section_seats (section_id, capacity, active_count)
SELECT s.id, s.capacity, count(e.id) FROM registration.sections s
LEFT JOIN registration.enrollments e ON e.section_id = s.id AND e.status = 'enrolled'
GROUP BY s.id;

INSERT INTO registration.outbox (event_id, stream, type, payload, created_at, published_at)
SELECT event_id, stream, type, payload, created_at, published_at FROM public.outbox
WHERE stream = 'registration' ORDER BY seq;

-- +goose StatementBegin
DO $constraints$
DECLARE
    service text;
    row record;
BEGIN
    SET LOCAL search_path = public;
    FOREACH service IN ARRAY ARRAY['identity', 'academic', 'registration', 'affairs', 'notification'] LOOP
        FOR row IN
            SELECT con.conname AS name, src.relname AS source, dst.relname AS target, pg_get_constraintdef(con.oid) AS definition
            FROM pg_constraint con
            JOIN pg_class src ON src.oid = con.conrelid
            JOIN pg_class dst ON dst.oid = con.confrelid
            WHERE con.contype = 'f' AND con.connamespace = 'public'::regnamespace
              AND to_regclass(format('%I.%I', service, src.relname)) IS NOT NULL
              AND to_regclass(format('%I.%I', service, dst.relname)) IS NOT NULL
              AND NOT (service = 'registration' AND src.relname = 'grades')
        LOOP
            EXECUTE format('ALTER TABLE %I.%I ADD CONSTRAINT %I %s', service, row.source, row.name,
                replace(row.definition, 'REFERENCES ' || quote_ident(row.target) || '(',
                    'REFERENCES ' || format('%I.%I', service, row.target) || '('));
        END LOOP;
        EXECUTE format('GRANT USAGE ON SCHEMA %I, portal67_events TO %I', service, 'portal67_' || service);
        EXECUTE format('GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA %I TO %I',
            service, 'portal67_' || service);
        EXECUTE format('GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA %I TO %I', service, 'portal67_' || service);
    END LOOP;
END
$constraints$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION registration.count_seats() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.status = 'enrolled' THEN
            UPDATE registration.section_seats SET active_count = active_count + 1 WHERE section_id = NEW.section_id;
        END IF;
    ELSIF OLD.status = 'enrolled' AND NEW.status = 'cancelled' THEN
        UPDATE registration.section_seats SET active_count = active_count - 1 WHERE section_id = NEW.section_id;
    END IF;
    RETURN NEW;
END
$$;
-- +goose StatementEnd
CREATE TRIGGER enrollment_seat_count AFTER INSERT OR UPDATE OF status ON registration.enrollments
FOR EACH ROW EXECUTE FUNCTION registration.count_seats();

-- +goose Down
-- This discards writes made in the split schemas. Restore service backups or
-- recover those writes before choosing the public checkpoint.
DROP SCHEMA notification CASCADE;
DROP SCHEMA affairs CASCADE;
DROP SCHEMA registration CASCADE;
DROP SCHEMA academic CASCADE;
DROP SCHEMA identity CASCADE;
DROP SCHEMA portal67_events CASCADE;
