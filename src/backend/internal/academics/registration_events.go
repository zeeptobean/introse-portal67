package academics

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events"
)

const enrollmentStateColumns = "id, student_id, section_id, status, registered_at, cancelled_at, version"

type enrollmentState struct {
	ID           int64      `json:"enrollment_id"`
	StudentID    int64      `json:"student_id"`
	SectionID    int64      `json:"section_id"`
	Status       string     `json:"status"`
	RegisteredAt time.Time  `json:"registered_at"`
	CancelledAt  *time.Time `json:"cancelled_at,omitempty"`
	Version      int64      `json:"version"`
}

func readEnrollmentState(row pgx.Row) (enrollmentState, error) {
	var state enrollmentState
	err := row.Scan(&state.ID, &state.StudentID, &state.SectionID, &state.Status,
		&state.RegisteredAt, &state.CancelledAt, &state.Version)
	return state, err
}

func (s enrollmentState) appendEvent(ctx context.Context, tx pgx.Tx, kind string, occurredAt time.Time) error {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("encode enrollment state: %w", err)
	}
	return events.Add(ctx, tx, "registration", events.Event{
		Type:          kind,
		SchemaVersion: 1,
		OccurredAt:    occurredAt.UTC(),
		Producer:      "registration-service",
		Aggregate:     events.Aggregate{Type: "enrollment", ID: strconv.FormatInt(s.ID, 10), Version: s.Version},
		ActorID:       s.StudentID,
		Data:          data,
	})
}
