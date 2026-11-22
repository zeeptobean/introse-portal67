package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/auth"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/httpx"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events"
)

type Item struct {
	EventID    string    `json:"event_id"`
	Kind       string    `json:"kind"`
	SubjectID  string    `json:"subject_id"`
	Message    string    `json:"message"`
	OccurredAt time.Time `json:"occurred_at"`
}

func Handler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := auth.UserFromContext(r.Context())
		ctx := r.Context()
		rows, err := pool.Query(ctx, `
			SELECT event_id::text, kind, subject_id, message, occurred_at
			FROM notifications WHERE recipient_id = $1
			ORDER BY occurred_at DESC, event_id DESC LIMIT 50`, user.ID)
		if err != nil {
			slog.ErrorContext(ctx, "list notifications", "error", err)
			httpx.InternalError(w)
			return
		}
		items, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Item])
		if err != nil {
			slog.ErrorContext(ctx, "decode notifications", "error", err)
			httpx.InternalError(w)
			return
		}
		if items == nil {
			items = []Item{}
		}
		httpx.WriteJSON(w, http.StatusOK, items)
	}
}

type template struct {
	message string
	staff   bool
}

var templates = map[string]template{
	"registration.enrollment.created":    {"You registered for a section.", false},
	"registration.enrollment.cancelled":  {"You cancelled a section registration.", false},
	"academic.grade.recorded":            {"A lecturer recorded or updated a grade.", false},
	"academic.grade.orphaned":            {"A grade needs review after an enrollment was cancelled.", true},
	"affairs.document_request.submitted": {"A document request is ready for review.", true},
	"affairs.document_request.decided":   {"Staff reviewed your document request.", false},
	"affairs.tuition.charge.added":       {"Staff added a tuition charge.", false},
	"affairs.tuition.payment.recorded":   {"Staff recorded a simulated tuition payment.", false},
	"affairs.tuition.payment.updated":    {"Staff updated a simulated tuition payment.", false},
	"affairs.tuition.payment.voided":     {"Staff voided a simulated tuition payment.", false},
}

func Project(ctx context.Context, tx pgx.Tx, event events.Event) error {
	template, ok := templates[event.Type]
	if !ok {
		return nil
	}
	if event.SchemaVersion != 1 {
		return fmt.Errorf("unsupported schema version for %s", event.Type)
	}
	if template.staff {
		_, err := tx.Exec(ctx, `
			INSERT INTO notifications (event_id, recipient_id, kind, subject_id, message, occurred_at)
			SELECT $1, id, $2, $3, $4, $5 FROM users WHERE role = 'staff' AND is_active
			ON CONFLICT DO NOTHING`, event.ID.String(), event.Type, event.Aggregate.ID, template.message, event.OccurredAt)
		return err
	}
	var data struct {
		StudentID int64 `json:"student_id"`
	}
	if err := json.Unmarshal(event.Data, &data); err != nil {
		return err
	}
	if data.StudentID <= 0 {
		return fmt.Errorf("%s is missing student_id", event.Type)
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO notifications (event_id, recipient_id, kind, subject_id, message, occurred_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT DO NOTHING`, event.ID.String(), data.StudentID, event.Type, event.Aggregate.ID, template.message, event.OccurredAt)
	return err
}
