-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION registration.count_seats() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    effective_capacity integer;
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.status = 'enrolled' THEN
            UPDATE registration.section_seats SET active_count = active_count + 1 WHERE section_id = NEW.section_id;
        END IF;
    ELSIF OLD.status = 'enrolled' AND NEW.status = 'cancelled' THEN
        UPDATE registration.section_seats
        SET active_count = active_count - 1,
            capacity = CASE WHEN pending_capacity IS NOT NULL THEN greatest(pending_capacity, active_count - 1) ELSE capacity END,
            pending_capacity = CASE WHEN pending_capacity < active_count - 1 THEN pending_capacity END
        WHERE section_id = NEW.section_id RETURNING capacity INTO effective_capacity;
        UPDATE registration.sections SET capacity = effective_capacity
        WHERE id = NEW.section_id AND capacity <> effective_capacity;
    END IF;
    RETURN NEW;
END
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION registration.count_seats() RETURNS trigger LANGUAGE plpgsql AS $$
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
