package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/routing"
)

// Republish emits current state for consumers that joined after stream
// history was cleared. Snapshot-only events do not create student alerts.
func Republish(ctx context.Context, pool *pgxpool.Pool, name string) error {
	var queries []string
	switch name {
	case routing.Identity:
		queries = []string{"SELECT identity.emit_user(u) FROM users u ORDER BY id"}
	case routing.Academic:
		queries = []string{
			"SELECT academic.emit_catalog('term', to_jsonb(t)) FROM terms t ORDER BY id",
			"SELECT academic.emit_catalog('course', to_jsonb(c)) FROM courses c ORDER BY id",
			"SELECT academic.emit_catalog('section', to_jsonb(s)) FROM sections s ORDER BY id",
			"SELECT academic.emit_catalog('meeting', to_jsonb(m)) FROM class_meetings m ORDER BY id",
			"SELECT academic.emit_catalog('exam', to_jsonb(e)) FROM exams e ORDER BY id",
			`SELECT portal67_events.emit('academic', 'academic.grade.snapshot', 'grade',
				enrollment_id, version, graded_by, jsonb_build_object('enrollment_id', enrollment_id))
			  FROM grades WHERE orphaned_at IS NULL ORDER BY enrollment_id`,
		}
	case routing.Registration:
		queries = []string{`
			SELECT portal67_events.emit('registration', 'registration.enrollment.snapshot', 'enrollment',
				id, version, student_id, (to_jsonb(e) - 'id') || jsonb_build_object('enrollment_id', id))
			FROM enrollments e ORDER BY id`}
	default:
		return fmt.Errorf("%s does not own snapshot projections", name)
	}
	return pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		for _, query := range queries {
			if _, err := tx.Exec(ctx, query); err != nil {
				return err
			}
		}
		return nil
	})
}
