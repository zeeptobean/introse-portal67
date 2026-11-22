package academics

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events"
)

func TestParallelRegistrationDoesNotOverfill(t *testing.T) {
	ctx, pool := registrationTestDB(t)
	const students = 16
	sectionID, studentIDs := seedRegistrationSection(t, ctx, pool, students)
	service := NewService(NewStore(pool))
	start := make(chan struct{})
	results := make(chan error, students)
	for _, studentID := range studentIDs {
		go func() {
			<-start
			_, err := service.Register(ctx, studentID, sectionID)
			results <- err
		}()
	}
	close(start)

	var successes, full int
	for range students {
		err := <-results
		switch {
		case err == nil:
			successes++
		case errors.Is(err, ErrSectionFull):
			full++
		default:
			t.Errorf("unexpected registration error: %v", err)
		}
	}
	if successes != 1 || full != students-1 {
		t.Fatalf("registrations succeeded = %d, section full = %d", successes, full)
	}
	var active int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM enrollments WHERE section_id = $1 AND status = 'enrolled'", sectionID).Scan(&active); err != nil {
		t.Fatal(err)
	}
	if active != 1 {
		t.Fatalf("active enrollments = %d, want 1", active)
	}
}

func TestRegistrationOutbox(t *testing.T) {
	ctx, pool := registrationTestDB(t)
	ctx = events.WithTraceID(ctx, "test-registration")
	sectionID, students := seedRegistrationSection(t, ctx, pool, 1)
	studentID := students[0]
	service := NewService(NewStore(pool))

	// Reject the outbox insert to prove the enrollment cannot commit on its own.
	if _, err := pool.Exec(ctx, "ALTER TABLE outbox ADD CONSTRAINT test_publish_failure CHECK (false)"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Register(ctx, studentID, sectionID); err == nil || !strings.Contains(err.Error(), "test_publish_failure") {
		t.Fatalf("expected outbox insert to fail, got %v", err)
	}
	if _, err := pool.Exec(ctx, "ALTER TABLE outbox DROP CONSTRAINT test_publish_failure"); err != nil {
		t.Fatal(err)
	}
	enrollment, err := service.Register(ctx, studentID, sectionID)
	if err != nil {
		t.Fatalf("seat was not available after rollback: %v", err)
	}
	if _, err := service.Cancel(ctx, studentID, enrollment.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Cancel(ctx, studentID, enrollment.ID); !errors.Is(err, ErrAlreadyCancelled) {
		t.Fatalf("duplicate cancellation: %v", err)
	}

	publishErr := errors.New("simulated Redis outage during batch")
	var firstAttempt, secondAttempt []events.Message
	relay := events.NewRelay(pool, publishFunc(func(_ context.Context, messages []events.Message) error {
		if firstAttempt == nil {
			firstAttempt = slices.Clone(messages)
			return publishErr
		}
		secondAttempt = slices.Clone(messages)
		return nil
	}))
	if count, err := relay.Flush(ctx); count != 0 || !errors.Is(err, publishErr) {
		t.Fatalf("failed flush = (%d, %v)", count, err)
	}
	if count, err := relay.Flush(ctx); count != 2 || err != nil {
		t.Fatalf("retry flush = (%d, %v)", count, err)
	}
	if count, err := relay.Flush(ctx); count != 0 || err != nil {
		t.Fatalf("empty flush = (%d, %v)", count, err)
	}
	if !slices.Equal(firstAttempt, secondAttempt) {
		t.Fatal("retry changed the event payload")
	}
	var previousID uuid.UUID
	want := []struct {
		kind   string
		status string
	}{
		{"registration.enrollment.created", StatusEnrolled},
		{"registration.enrollment.cancelled", StatusCancelled},
	}
	if len(secondAttempt) != len(want) {
		t.Fatalf("published events = %d, want %d", len(secondAttempt), len(want))
	}
	for i, message := range secondAttempt {
		var event events.Event
		if err := json.Unmarshal([]byte(message.Payload), &event); err != nil {
			t.Fatal(err)
		}
		var state enrollmentState
		if err := json.Unmarshal(event.Data, &state); err != nil {
			t.Fatal(err)
		}
		version := int64(i + 1)
		occurredAt := state.RegisteredAt
		if state.CancelledAt != nil {
			occurredAt = *state.CancelledAt
		}
		if message.Stream != "registration" || message.Type != want[i].kind ||
			event.Type != message.Type || event.SchemaVersion != 1 || event.Producer != "registration-service" ||
			event.ID.Version() != 7 || event.ID == previousID || event.TraceID != "test-registration" ||
			event.Aggregate.Type != "enrollment" || event.Aggregate.ID != strconv.FormatInt(enrollment.ID, 10) ||
			event.Aggregate.Version != version || !event.OccurredAt.Equal(occurredAt) ||
			event.ActorID != studentID || state.ID != enrollment.ID || state.StudentID != studentID ||
			state.SectionID != sectionID || state.Version != version || state.Status != want[i].status ||
			(state.CancelledAt != nil) != (version == 2) || !state.RegisteredAt.Equal(enrollment.RegisteredAt) {
			t.Errorf("unexpected event %s: %s", message.Type, message.Payload)
		}
		previousID = event.ID
	}
}

type publishFunc func(context.Context, []events.Message) error

func (f publishFunc) Publish(ctx context.Context, messages []events.Message) error {
	return f(ctx, messages)
}

func registrationTestDB(t *testing.T) (context.Context, *pgxpool.Pool) {
	t.Helper()
	dsn := os.Getenv("PORTAL67_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set PORTAL67_TEST_DATABASE_URL to run PostgreSQL integration tests")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	t.Cleanup(cancel)
	admin, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal("invalid integration test database configuration")
	}
	t.Cleanup(admin.Close)
	schema := "portal67_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	quoted := pgx.Identifier{schema}.Sanitize()
	if _, err := admin.Exec(ctx, "CREATE SCHEMA "+quoted); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := admin.Exec(cleanup, "DROP SCHEMA "+quoted+" CASCADE"); err != nil {
			t.Errorf("drop integration schema: %v", err)
		}
	})
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal("invalid integration test database configuration")
	}
	config.ConnConfig.RuntimeParams["search_path"] = schema
	config.MaxConns = 20
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	migrations, err := filepath.Glob("../../db/migrations/*.sql")
	if err != nil || len(migrations) == 0 {
		t.Fatalf("find migrations: %v", err)
	}
	for _, path := range migrations {
		sql, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		up, _, _ := strings.Cut(string(sql), "-- +goose Down")
		if _, err := pool.Exec(ctx, up); err != nil {
			t.Fatalf("%s: %v", path, err)
		}
	}
	return ctx, pool
}

func seedRegistrationSection(t *testing.T, ctx context.Context, pool *pgxpool.Pool, students int) (int64, []int64) {
	t.Helper()
	var sectionID int64
	err := pool.QueryRow(ctx, `
		WITH lecturer AS (
			INSERT INTO users (email, password_hash, full_name, date_of_birth, role)
			VALUES ('lecturer@example.test', 'unused', 'Test Lecturer', '1990-01-01', 'lecturer') RETURNING id
		), scale AS (
			INSERT INTO grading_scales (name) VALUES ('Test scale') RETURNING id
		), term AS (
			INSERT INTO terms (code, name, starts_on, ends_on, grading_scale_id,
				registration_opens_at, registration_closes_at, grading_closes_at)
			SELECT 'T1', 'Test term', current_date, current_date + 60, id,
				now() - interval '1 day', now() + interval '1 day', now() + interval '60 days' FROM scale RETURNING id
		), course AS (
			INSERT INTO courses (code, name, credits) VALUES ('C1', 'Test course', 3) RETURNING id
		)
		INSERT INTO sections (course_id, term_id, lecturer_id, section_code, capacity)
		SELECT course.id, term.id, lecturer.id, 'S1', 1 FROM course, term, lecturer RETURNING id`).Scan(&sectionID)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := pool.Query(ctx, `
		INSERT INTO users (email, password_hash, full_name, date_of_birth, role)
		SELECT 'student' || n || '@example.test', 'unused', 'Test Student', '2005-01-01', 'student'
		FROM generate_series(1, $1::integer) AS n RETURNING id`, students)
	if err != nil {
		t.Fatal(err)
	}
	studentIDs, err := pgx.CollectRows(rows, pgx.RowTo[int64])
	if err != nil {
		t.Fatal(err)
	}
	return sectionID, studentIDs
}
