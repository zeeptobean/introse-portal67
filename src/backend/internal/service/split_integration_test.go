package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/academics"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/projections"
)

func TestSplitStorageAndCancellationProjection(t *testing.T) {
	ctx, pool, dsn := serviceTestDB(t)
	asRole := func(service string, apply func(pgx.Tx) error) error {
		return pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
			if _, err := tx.Exec(ctx, "SET LOCAL ROLE "+pgx.Identifier{"portal67_" + service}.Sanitize()); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "SET LOCAL search_path TO "+pgx.Identifier{service}.Sanitize()); err != nil {
				return err
			}
			return apply(tx)
		})
	}
	err := asRole("registration", func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "SELECT password_hash FROM identity.users LIMIT 1")
		return err
	})
	var permissionErr *pgconn.PgError
	if !errors.As(err, &permissionErr) || permissionErr.Code != "42501" {
		t.Fatalf("registration role could read identity credentials or returned wrong error: %v", err)
	}
	var foreignKeys int
	if err := pool.QueryRow(ctx, `
		SELECT count(*) FROM pg_constraint con
		JOIN pg_class own ON own.oid = con.conrelid
		JOIN pg_class ref ON ref.oid = con.confrelid
		WHERE con.contype = 'f' AND own.relnamespace <> 'public'::regnamespace
			AND own.relnamespace <> ref.relnamespace`).Scan(&foreignKeys); err != nil || foreignKeys != 0 {
		t.Fatalf("cross-schema foreign keys = %d, err = %v", foreignKeys, err)
	}
	if _, err := pool.Exec(ctx, `
		UPDATE registration.terms SET registration_opens_at = now() - interval '1 day',
			registration_closes_at = now() + interval '1 day', grading_closes_at = now() + interval '2 days';
		UPDATE academic.sections SET capacity = capacity + 1 WHERE id = 9`); err != nil {
		t.Fatal(err)
	}
	catalog := loadEvent(t, ctx, pool, "academic", "academic.section.upserted")
	if err := asRole("registration", func(tx pgx.Tx) error {
		return projections.For("registration")["academic"](ctx, tx, catalog)
	}); err != nil {
		t.Fatal(err)
	}
	var synced bool
	if err := pool.QueryRow(ctx, `
		SELECT a.capacity = r.capacity FROM academic.sections a
		JOIN registration.section_seats r ON r.section_id = a.id WHERE a.id = 9`).Scan(&synced); err != nil || !synced {
		t.Fatalf("catalog capacity synced = %t, err = %v", synced, err)
	}

	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal("invalid integration database configuration")
	}
	dbConfig.ConnConfig.RuntimeParams["role"] = "portal67_registration"
	dbConfig.ConnConfig.RuntimeParams["search_path"] = "registration"
	registration, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(registration.Close)
	api := academics.NewService(academics.NewRegistrationStore(registration))
	enrollment, err := api.Register(ctx, 1, 9)
	if err != nil {
		t.Fatal(err)
	}
	applyRoster := func(kind string) {
		t.Helper()
		event := loadEvent(t, ctx, pool, "registration", kind)
		if err := asRole("academic", func(tx pgx.Tx) error {
			return projections.For("academic")["registration"](ctx, tx, event)
		}); err != nil {
			t.Fatal(err)
		}
	}
	applyRoster("registration.enrollment.created")
	if err := asRole("academic", func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "INSERT INTO grades (enrollment_id, score, graded_by) SELECT $1, 8, lecturer_id FROM sections WHERE id = 9", enrollment.ID)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	// Delay the grade event in transit, then deliver cancellation to academic.
	if _, err := api.Cancel(ctx, 1, enrollment.ID); err != nil {
		t.Fatal(err)
	}
	applyRoster("registration.enrollment.cancelled")
	orphan := loadEvent(t, ctx, pool, "academic", "academic.grade.orphaned")
	if orphan.Aggregate.ID != strconv.FormatInt(enrollment.ID, 10) {
		t.Fatalf("orphan event refers to %s", orphan.Aggregate.ID)
	}
	var active int
	if err := pool.QueryRow(ctx, "SELECT active_count FROM registration.section_seats WHERE section_id = 9").Scan(&active); err != nil || active != 0 {
		t.Fatalf("held seats = %d, err = %v", active, err)
	}

	if _, err := pool.Exec(ctx, `
		UPDATE registration.section_seats SET capacity = 1 WHERE section_id = 9;
		UPDATE registration.sections SET capacity = 1 WHERE id = 9;
		INSERT INTO registration.users (id, full_name, email, role, is_active, version)
		SELECT id, 'Parallel Student', 'parallel' || id || '@example.test', 'student', true, 1
		FROM generate_series(100, 111) id`); err != nil {
		t.Fatal(err)
	}
	type registrationResult struct {
		student    int64
		enrollment academics.Enrollment
		err        error
	}
	start := make(chan struct{})
	completed := make(chan registrationResult, 12)
	for student := int64(100); student <= 111; student++ {
		go func() {
			<-start
			enrollment, err := api.Register(ctx, student, 9)
			completed <- registrationResult{student, enrollment, err}
		}()
	}
	close(start)
	var winner registrationResult
	for range 12 {
		result := <-completed
		if result.err == nil {
			if winner.student != 0 {
				t.Error("more than one registration won the final seat")
			}
			winner = result
		} else if !errors.Is(result.err, academics.ErrSectionFull) {
			t.Errorf("registration returned %v, want ErrSectionFull", result.err)
		}
	}
	if winner.student == 0 {
		t.Fatal("nobody won the final seat")
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO academic.users SELECT * FROM registration.users WHERE id = $1`, winner.student); err != nil {
		t.Fatal(err)
	}
	applyRoster("registration.enrollment.created")
	if err := asRole("academic", func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "INSERT INTO grades (enrollment_id, score, graded_by) SELECT $1, 8, lecturer_id FROM sections WHERE id = 9", winner.enrollment.ID)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	grade := loadEvent(t, ctx, pool, "academic", "academic.grade.recorded")
	if err := asRole("registration", func(tx pgx.Tx) error {
		return projections.For("registration")["academic"](ctx, tx, grade)
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := api.Cancel(ctx, winner.student, winner.enrollment.ID); !errors.Is(err, academics.ErrAlreadyGraded) {
		t.Fatalf("cancelling a projected graded enrollment = %v, want ErrAlreadyGraded", err)
	}
}

func serviceTestDB(t *testing.T) (context.Context, *pgxpool.Pool, string) {
	t.Helper()
	raw := os.Getenv("PORTAL67_SPLIT_TEST_ADMIN_URL")
	if raw == "" {
		t.Skip("set PORTAL67_SPLIT_TEST_ADMIN_URL to test the split on an isolated PostgreSQL cluster")
	}
	dsn, err := url.Parse(raw)
	if err != nil || dsn.Host == "" {
		t.Fatal("invalid split test database configuration")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	t.Cleanup(cancel)
	admin, err := pgx.Connect(ctx, raw)
	if err != nil {
		t.Fatal("could not connect to the split test cluster")
	}
	t.Cleanup(func() { _ = admin.Close(context.Background()) })
	name := "p67_test_" + uuid.New().String()
	quoted := pgx.Identifier{name}.Sanitize()
	if _, err := admin.Exec(ctx, "CREATE DATABASE "+quoted); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cleanup, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()
		if _, err := admin.Exec(cleanup, "DROP DATABASE "+quoted); err != nil {
			t.Errorf("drop test database: %v", err)
		}
	})
	dsn.Path = "/" + name
	if err := Migrate(ctx, dsn.String(), "../../db", true, nil); err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.New(ctx, dsn.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	return ctx, pool, dsn.String()
}

func loadEvent(t *testing.T, ctx context.Context, pool *pgxpool.Pool, schema, kind string) events.Event {
	t.Helper()
	var data []byte
	query := "SELECT payload FROM " + pgx.Identifier{schema, "outbox"}.Sanitize() + " WHERE type = $1 ORDER BY seq DESC LIMIT 1"
	if err := pool.QueryRow(ctx, query, kind).Scan(&data); err != nil {
		t.Fatal(err)
	}
	var event events.Event
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatal(err)
	}
	return event
}
