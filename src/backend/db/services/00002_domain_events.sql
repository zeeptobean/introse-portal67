-- +goose Up
ALTER TABLE registration.catalog_tombstones RENAME TO catalog_versions;
-- +goose StatementBegin
CREATE FUNCTION portal67_events.emit(
    owner text, kind text, entity text, entity_id bigint, version bigint, actor bigint, data jsonb
) RETURNS void LANGUAGE plpgsql AS $$
DECLARE
    event_id uuid := uuidv7();
BEGIN
    EXECUTE format('INSERT INTO %I.outbox (event_id, stream, type, payload) VALUES ($1, $2, $3, $4)', owner)
    USING event_id, owner, kind, jsonb_build_object(
        'id', event_id, 'type', kind, 'schema_version', 1, 'occurred_at', clock_timestamp(),
        'producer', owner || '-service', 'actor_id', actor,
        'trace_id', COALESCE(NULLIF(current_setting('portal67.trace_id', true), ''), event_id::text),
        'aggregate', jsonb_build_object('type', entity, 'id', entity_id::text, 'version', version),
        'data', data
    );
END
$$;

CREATE FUNCTION portal67_events.bump_version() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    excluded text[] := string_to_array(TG_ARGV[0], ',');
BEGIN
    NEW.version := OLD.version;
    IF (to_jsonb(NEW) - excluded) IS DISTINCT FROM (to_jsonb(OLD) - excluded) THEN
        NEW.version := OLD.version + 1;
    END IF;
    RETURN NEW;
END
$$;

CREATE FUNCTION identity.emit_user(value identity.users) RETURNS void LANGUAGE sql AS $$
    SELECT portal67_events.emit('identity',
        CASE WHEN value.is_active THEN 'identity.user.upserted' ELSE 'identity.user.deactivated' END,
        'user', value.id, value.version, 0, jsonb_build_object(
            'id', value.id, 'institutional_id', value.institutional_id, 'full_name', value.full_name,
            'email', value.email, 'role', value.role, 'is_active', value.is_active, 'version', value.version
        ))
$$;

CREATE FUNCTION identity.user_event() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF TG_OP = 'INSERT' OR NEW.version <> OLD.version THEN
        PERFORM identity.emit_user(NEW);
    END IF;
    RETURN NULL;
END
$$;

CREATE FUNCTION academic.emit_catalog(entity text, data jsonb, deleted boolean DEFAULT false)
RETURNS void LANGUAGE plpgsql AS $$
DECLARE
    row_id bigint := (data->>'id')::bigint;
    version bigint := (data->>'version')::bigint;
BEGIN
    IF entity = 'section' AND NOT deleted THEN
        data := data || (
            SELECT jsonb_build_object(
                'course_code', c.code, 'course_name', c.name, 'credits', c.credits,
                'lecturer_name', u.full_name,
                'meetings', COALESCE((SELECT jsonb_agg(to_jsonb(m) ORDER BY m.id)
                    FROM academic.class_meetings m WHERE m.section_id = row_id), '[]'::jsonb)
            )
            FROM academic.courses c JOIN academic.users u ON u.id = (data->>'lecturer_id')::bigint
            WHERE c.id = (data->>'course_id')::bigint
        );
    END IF;
    PERFORM portal67_events.emit('academic',
        'academic.' || entity || CASE WHEN deleted THEN '.deleted' ELSE '.upserted' END,
        entity, row_id, version + CASE WHEN deleted THEN 1 ELSE 0 END, 0, data);
END
$$;

CREATE FUNCTION academic.catalog_event() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        PERFORM academic.emit_catalog(TG_ARGV[0], to_jsonb(OLD), true);
    ELSIF TG_OP = 'INSERT' OR NEW.version <> OLD.version THEN
        PERFORM academic.emit_catalog(TG_ARGV[0], to_jsonb(NEW));
    END IF;
    RETURN NULL;
END
$$;

CREATE FUNCTION academic.emit_grade(value academic.grades) RETURNS void LANGUAGE plpgsql AS $$
DECLARE
    roster_row record;
    kind text := 'academic.grade.recorded';
    version bigint := value.version;
BEGIN
    SELECT e.student_id, e.section_id, (
        SELECT band.letter_grade FROM academic.grade_bands band
        WHERE band.grading_scale_id = term.grading_scale_id AND band.lower_bound <= value.score
        ORDER BY band.lower_bound DESC LIMIT 1
    ) AS letter_grade INTO STRICT roster_row
    FROM academic.enrollments e
    JOIN academic.sections section ON section.id = e.section_id
    JOIN academic.terms term ON term.id = section.term_id
    WHERE e.id = value.enrollment_id;
    IF value.orphaned_at IS NOT NULL THEN
        kind := 'academic.grade.orphaned';
    END IF;
    PERFORM portal67_events.emit('academic', kind, 'grade', value.enrollment_id, version, value.graded_by,
        to_jsonb(value) || to_jsonb(roster_row));
END
$$;

CREATE FUNCTION academic.grade_event() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF TG_OP = 'INSERT' OR NEW.version <> OLD.version OR NEW.orphaned_at IS DISTINCT FROM OLD.orphaned_at THEN
        PERFORM academic.emit_grade(NEW);
    END IF;
    RETURN NULL;
END
$$;

CREATE FUNCTION affairs.document_event() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    request affairs.document_requests;
    kind text := CASE WHEN NEW.to_status = 'pending' THEN 'submitted' ELSE 'decided' END;
BEGIN
    SELECT * INTO STRICT request FROM affairs.document_requests WHERE id = NEW.document_request_id;
    PERFORM portal67_events.emit('affairs', 'affairs.document_request.' || kind, 'document_request',
        request.id, NEW.id, NEW.actor_id, jsonb_build_object(
            'id', request.id, 'student_id', request.student_id, 'request_type', request.request_type,
            'status', NEW.to_status, 'submitted_at', request.submitted_at, 'comment', NEW.comment,
            'actor_id', NEW.actor_id, 'decided_at', CASE WHEN kind = 'decided' THEN NEW.created_at END
        ));
    RETURN NULL;
END
$$;

CREATE FUNCTION affairs.payment_event() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    payment affairs.tuition_payments;
BEGIN
    IF NEW.action = 'snapshot' THEN
        RETURN NULL;
    END IF;
    SELECT * INTO STRICT payment FROM affairs.tuition_payments WHERE id = NEW.payment_id;
    PERFORM portal67_events.emit('affairs', 'affairs.tuition.payment.' || NEW.action, 'payment',
        payment.id, NEW.id, NEW.actor_id, to_jsonb(payment));
    RETURN NULL;
END
$$;

CREATE FUNCTION affairs.charge_event() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    PERFORM portal67_events.emit('affairs', 'affairs.tuition.charge.added', 'charge', NEW.id, 1, NEW.created_by,
        jsonb_build_object('id', NEW.id, 'student_id', NEW.student_id, 'term_id', NEW.term_id, 'amount_vnd', NEW.amount_vnd));
    RETURN NULL;
END
$$;
-- +goose StatementEnd

CREATE TRIGGER user_version BEFORE UPDATE ON identity.users FOR EACH ROW
    EXECUTE FUNCTION portal67_events.bump_version('version,password_hash,date_of_birth');
CREATE TRIGGER user_outbox AFTER INSERT OR UPDATE ON identity.users FOR EACH ROW
    EXECUTE FUNCTION identity.user_event();
CREATE TRIGGER term_version BEFORE UPDATE ON academic.terms FOR EACH ROW
    EXECUTE FUNCTION portal67_events.bump_version('version');
CREATE TRIGGER term_outbox AFTER INSERT OR UPDATE ON academic.terms FOR EACH ROW
    EXECUTE FUNCTION academic.catalog_event('term');
CREATE TRIGGER course_version BEFORE UPDATE ON academic.courses FOR EACH ROW
    EXECUTE FUNCTION portal67_events.bump_version('version');
CREATE TRIGGER course_outbox AFTER INSERT OR UPDATE ON academic.courses FOR EACH ROW
    EXECUTE FUNCTION academic.catalog_event('course');
CREATE TRIGGER section_version BEFORE UPDATE ON academic.sections FOR EACH ROW
    EXECUTE FUNCTION portal67_events.bump_version('version');
CREATE TRIGGER section_outbox AFTER INSERT OR UPDATE ON academic.sections FOR EACH ROW
    EXECUTE FUNCTION academic.catalog_event('section');
CREATE TRIGGER meeting_version BEFORE UPDATE ON academic.class_meetings FOR EACH ROW
    EXECUTE FUNCTION portal67_events.bump_version('version');
CREATE TRIGGER meeting_outbox AFTER INSERT OR UPDATE OR DELETE ON academic.class_meetings FOR EACH ROW
    EXECUTE FUNCTION academic.catalog_event('meeting');
CREATE TRIGGER exam_version BEFORE UPDATE ON academic.exams FOR EACH ROW
    EXECUTE FUNCTION portal67_events.bump_version('version');
CREATE TRIGGER exam_outbox AFTER INSERT OR UPDATE OR DELETE ON academic.exams FOR EACH ROW
    EXECUTE FUNCTION academic.catalog_event('exam');
CREATE TRIGGER grade_version BEFORE UPDATE ON academic.grades FOR EACH ROW
    EXECUTE FUNCTION portal67_events.bump_version('version,orphaned_at');
CREATE TRIGGER grade_outbox AFTER INSERT OR UPDATE ON academic.grades FOR EACH ROW
    EXECUTE FUNCTION academic.grade_event();
CREATE TRIGGER document_outbox AFTER INSERT ON affairs.document_request_events FOR EACH ROW
    EXECUTE FUNCTION affairs.document_event();
CREATE TRIGGER payment_outbox AFTER INSERT ON affairs.tuition_payment_events FOR EACH ROW
    EXECUTE FUNCTION affairs.payment_event();
CREATE TRIGGER charge_outbox AFTER INSERT ON affairs.tuition_charges FOR EACH ROW
    EXECUTE FUNCTION affairs.charge_event();

-- +goose Down
DROP TRIGGER charge_outbox ON affairs.tuition_charges;
DROP TRIGGER payment_outbox ON affairs.tuition_payment_events;
DROP TRIGGER document_outbox ON affairs.document_request_events;
DROP TRIGGER grade_outbox ON academic.grades;
DROP TRIGGER grade_version ON academic.grades;
DROP TRIGGER exam_outbox ON academic.exams;
DROP TRIGGER exam_version ON academic.exams;
DROP TRIGGER meeting_outbox ON academic.class_meetings;
DROP TRIGGER meeting_version ON academic.class_meetings;
DROP TRIGGER section_outbox ON academic.sections;
DROP TRIGGER section_version ON academic.sections;
DROP TRIGGER course_outbox ON academic.courses;
DROP TRIGGER course_version ON academic.courses;
DROP TRIGGER term_outbox ON academic.terms;
DROP TRIGGER term_version ON academic.terms;
DROP TRIGGER user_outbox ON identity.users;
DROP TRIGGER user_version ON identity.users;
DROP FUNCTION affairs.charge_event(), affairs.payment_event(), affairs.document_event();
DROP FUNCTION academic.grade_event(), academic.emit_grade(academic.grades);
DROP FUNCTION academic.catalog_event(), academic.emit_catalog(text, jsonb, boolean);
DROP FUNCTION identity.user_event(), identity.emit_user(identity.users);
DROP FUNCTION portal67_events.bump_version(), portal67_events.emit(text, text, text, bigint, bigint, bigint, jsonb);
ALTER TABLE registration.catalog_versions RENAME TO catalog_tombstones;
