package projections

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/routing"
)

func For(service string) map[string]events.Handler {
	switch service {
	case routing.Academic:
		return map[string]events.Handler{"identity": User, "registration": roster}
	case routing.Registration:
		return map[string]events.Handler{"identity": User, "academic": registrationCatalog}
	case routing.Affairs:
		return map[string]events.Handler{"identity": User, "academic": affairsTerm}
	default:
		return nil
	}
}

// User projects only directory fields. Identity retains credential ownership.
func User(ctx context.Context, tx pgx.Tx, event events.Event) error {
	if event.Type != "identity.user.upserted" && event.Type != "identity.user.deactivated" {
		return nil
	}
	if err := versionOne(event); err != nil {
		return err
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO users (id, institutional_id, full_name, email, role, is_active, version)
		SELECT id, institutional_id, full_name, email, role, is_active, version
		FROM jsonb_populate_record(NULL::users, $1)
		ON CONFLICT (id) DO UPDATE SET institutional_id = excluded.institutional_id,
			full_name = excluded.full_name, email = excluded.email, role = excluded.role,
			is_active = excluded.is_active, version = excluded.version
		WHERE users.version < excluded.version`, event.Data)
	return err
}

func roster(ctx context.Context, tx pgx.Tx, event events.Event) error {
	if event.Type != "registration.enrollment.created" && event.Type != "registration.enrollment.cancelled" &&
		event.Type != "registration.enrollment.snapshot" {
		return nil
	}
	if err := versionOne(event); err != nil {
		return err
	}
	var row struct {
		ID           int64      `json:"enrollment_id"`
		StudentID    int64      `json:"student_id"`
		SectionID    int64      `json:"section_id"`
		Status       string     `json:"status"`
		RegisteredAt time.Time  `json:"registered_at"`
		CancelledAt  *time.Time `json:"cancelled_at"`
		Version      int64      `json:"version"`
	}
	if err := json.Unmarshal(event.Data, &row); err != nil {
		return err
	}
	tag, err := tx.Exec(ctx, `
		INSERT INTO enrollments (id, student_id, section_id, status, registered_at, cancelled_at, version)
		OVERRIDING SYSTEM VALUE VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET status = excluded.status,
			cancelled_at = excluded.cancelled_at, version = excluded.version
		WHERE enrollments.version < excluded.version`, row.ID, row.StudentID, row.SectionID,
		row.Status, row.RegisteredAt, row.CancelledAt, row.Version)
	if err != nil || tag.RowsAffected() == 0 || row.Status != "cancelled" {
		return err
	}
	if _, err := tx.Exec(ctx, "SELECT set_config('portal67.trace_id', $1, true)", event.TraceID); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "UPDATE grades SET orphaned_at = now() WHERE enrollment_id = $1 AND orphaned_at IS NULL", row.ID)
	return err
}

var term = catalogRow{"terms", "code name starts_on ends_on grading_scale_id registration_opens_at registration_closes_at grading_closes_at version"}

func affairsTerm(ctx context.Context, tx pgx.Tx, event events.Event) error {
	if event.Type != "academic.term.upserted" {
		return nil
	}
	if err := versionOne(event); err != nil {
		return err
	}
	return term.upsert(ctx, tx, event.Data, "")
}

type catalogRow struct{ table, fields string }

var catalog = map[string]catalogRow{
	"term":    term,
	"course":  {"courses", "code name credits kind version"},
	"section": {"sections", "course_id term_id lecturer_id section_code capacity version"},
	"meeting": {"class_meetings", "section_id weekday starts_at ends_at room version"},
	"exam":    {"exams", "section_id label starts_at ends_at room version"},
}

func registrationCatalog(ctx context.Context, tx pgx.Tx, event events.Event) error {
	if event.Type == "academic.grade.recorded" || event.Type == "academic.grade.snapshot" {
		if err := versionOne(event); err != nil {
			return err
		}
		id, err := strconv.ParseInt(event.Aggregate.ID, 10, 64)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO grades (enrollment_id, version) VALUES ($1, $2)
			ON CONFLICT (enrollment_id) DO UPDATE SET version = excluded.version
			WHERE grades.version < excluded.version`, id, event.Aggregate.Version)
		return err
	}
	entity := event.Aggregate.Type
	row, ok := catalog[entity]
	deleted := event.Type == "academic."+entity+".deleted" && (entity == "exam" || entity == "meeting")
	if !ok || (!deleted && event.Type != "academic."+entity+".upserted") {
		return nil
	}
	if err := versionOne(event); err != nil {
		return err
	}
	id, err := strconv.ParseInt(event.Aggregate.ID, 10, 64)
	if err != nil {
		return err
	}
	tag, err := tx.Exec(ctx, `
		INSERT INTO catalog_versions (entity, id, version) VALUES ($1, $2, $3)
		ON CONFLICT (entity, id) DO UPDATE SET version = excluded.version
		WHERE catalog_versions.version < excluded.version`, entity, id, event.Aggregate.Version)
	if err != nil || tag.RowsAffected() == 0 {
		return err
	}
	if deleted {
		_, err := tx.Exec(ctx, "DELETE FROM "+row.table+" WHERE id = $1 AND version < $2", id, event.Aggregate.Version)
		return err
	}
	if entity == "section" {
		return projectSection(ctx, tx, row, id, event.Data)
	}
	return row.upsert(ctx, tx, event.Data, "")
}

func projectSection(ctx context.Context, tx pgx.Tx, row catalogRow, id int64, data json.RawMessage) error {
	// Match the registration lock order when the projected section already exists.
	if _, err := tx.Exec(ctx, "SELECT 1 FROM section_seats WHERE section_id = $1 FOR UPDATE", id); err != nil {
		return err
	}
	if err := row.upsert(ctx, tx, data, "capacity"); err != nil {
		return err
	}
	var value struct {
		Capacity int32 `json:"capacity"`
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	var effective int32
	err := tx.QueryRow(ctx, `
		INSERT INTO section_seats (section_id, capacity) VALUES ($1, $2)
		ON CONFLICT (section_id) DO UPDATE SET
			capacity = CASE WHEN excluded.capacity >= section_seats.active_count THEN excluded.capacity ELSE section_seats.capacity END,
			pending_capacity = CASE WHEN excluded.capacity < section_seats.active_count THEN excluded.capacity END
		RETURNING capacity`, id, value.Capacity).Scan(&effective)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "UPDATE sections SET capacity = $2 WHERE id = $1 AND capacity <> $2", id, effective)
	return err
}

func (r catalogRow) upsert(ctx context.Context, tx pgx.Tx, data json.RawMessage, preserve string) error {
	fields := strings.Fields(r.fields)
	updates := make([]string, 0, len(fields))
	for _, field := range fields {
		if field != preserve {
			updates = append(updates, field+" = excluded."+field)
		}
	}
	columns := "id, " + strings.Join(fields, ", ")
	// Table and column names are declared in this package.
	query := "INSERT INTO " + r.table + " (" + columns + ") OVERRIDING SYSTEM VALUE " +
		"SELECT " + columns + " FROM jsonb_populate_record(NULL::" + r.table + ", $1) " +
		"ON CONFLICT (id) DO UPDATE SET " + strings.Join(updates, ", ") + " WHERE " + r.table + ".version < excluded.version"
	_, err := tx.Exec(ctx, query, data)
	return err
}

func versionOne(event events.Event) error {
	if event.SchemaVersion != 1 {
		return fmt.Errorf("unsupported schema version %d for %s", event.SchemaVersion, event.Type)
	}
	return nil
}
