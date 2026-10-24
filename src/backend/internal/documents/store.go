package documents

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store holds the database dependency for document requests and staff decisions.
type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

// withTx runs fn in a transaction, committing when fn returns nil.
func (s *Store) withTx(ctx context.Context, fn func(pgx.Tx) error) error {
	return pgx.BeginFunc(ctx, s.db, fn)
}

// requestQuery selects requests with their student and the decision event,
// if any.
const requestQuery = `
	SELECT r.id, u.id, u.institutional_id, u.full_name, u.email,
		r.request_type, r.purpose, r.status, r.submitted_at, d.created_at, d.comment
	FROM document_requests r
	JOIN users u ON u.id = r.student_id
	LEFT JOIN document_request_events d
		ON d.document_request_id = r.id AND d.from_status = 'pending'`

func scanRequest(row pgx.CollectableRow) (Request, error) {
	var r Request
	err := row.Scan(&r.ID, &r.Student.ID, &r.Student.InstitutionalID, &r.Student.FullName, &r.Student.Email,
		&r.RequestType, &r.Purpose, &r.Status, &r.SubmittedAt, &r.DecidedAt, &r.DecisionComment)
	return r, err
}

// listRequests returns requests, newest first. studentID restricts the list to
// one student when non-nil.
func (s *Store) listRequests(ctx context.Context, studentID *int64, status *Status) ([]Request, error) {
	rows, err := s.db.Query(ctx, requestQuery+`
		WHERE ($1::bigint IS NULL OR r.student_id = $1) AND ($2::text IS NULL OR r.status = $2)
		ORDER BY r.submitted_at DESC, r.id DESC`, studentID, status)
	if err != nil {
		return nil, fmt.Errorf("list document requests: %w", err)
	}
	requests, err := pgx.CollectRows(rows, scanRequest)
	if err != nil {
		return nil, fmt.Errorf("list document requests: %w", err)
	}
	return requests, nil
}

// requestByID returns pgx.ErrNoRows (wrapped) when the request does not exist.
func (s *Store) requestByID(ctx context.Context, id int64) (Request, error) {
	rows, err := s.db.Query(ctx, requestQuery+" WHERE r.id = $1", id)
	if err != nil {
		return Request{}, fmt.Errorf("find document request %d: %w", id, err)
	}
	r, err := pgx.CollectExactlyOneRow(rows, scanRequest)
	if err != nil {
		return Request{}, fmt.Errorf("find document request %d: %w", id, err)
	}
	return r, nil
}

func (s *Store) history(ctx context.Context, requestID int64) ([]Event, error) {
	rows, err := s.db.Query(ctx, `
		SELECT e.id, u.id, u.full_name, u.role, e.from_status, e.to_status, e.comment, e.created_at
		FROM document_request_events e
		JOIN users u ON u.id = e.actor_id
		WHERE e.document_request_id = $1
		ORDER BY e.created_at, e.id`, requestID)
	if err != nil {
		return nil, fmt.Errorf("list history of document request %d: %w", requestID, err)
	}
	events, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Event, error) {
		var e Event
		err := row.Scan(&e.ID, &e.Actor.ID, &e.Actor.FullName, &e.Actor.Role,
			&e.FromStatus, &e.ToStatus, &e.Comment, &e.CreatedAt)
		return e, err
	})
	if err != nil {
		return nil, fmt.Errorf("list history of document request %d: %w", requestID, err)
	}
	return events, nil
}

func (s *Store) insertRequest(ctx context.Context, tx pgx.Tx, studentID int64, in SubmitInput) (int64, error) {
	var id int64
	err := tx.QueryRow(ctx, `
		INSERT INTO document_requests (student_id, request_type, purpose)
		VALUES ($1, $2, $3)
		RETURNING id`, studentID, in.RequestType, in.Purpose).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert document request: %w", err)
	}
	return id, nil
}

func (s *Store) insertEvent(ctx context.Context, tx pgx.Tx, requestID, actorID int64, from *Status, to Status, comment *string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO document_request_events (document_request_id, actor_id, from_status, to_status, comment)
		VALUES ($1, $2, $3, $4, $5)`, requestID, actorID, from, to, comment)
	if err != nil {
		return fmt.Errorf("insert event for document request %d: %w", requestID, err)
	}
	return nil
}

// lockRequestStatus locks the request row and returns its current status.
func (s *Store) lockRequestStatus(ctx context.Context, tx pgx.Tx, id int64) (Status, error) {
	var status Status
	err := tx.QueryRow(ctx, "SELECT status FROM document_requests WHERE id = $1 FOR UPDATE", id).Scan(&status)
	if err != nil {
		return "", fmt.Errorf("lock document request %d: %w", id, err)
	}
	return status, nil
}

func (s *Store) updateStatus(ctx context.Context, tx pgx.Tx, id int64, status Status) error {
	_, err := tx.Exec(ctx, "UPDATE document_requests SET status = $2 WHERE id = $1", id, status)
	if err != nil {
		return fmt.Errorf("update document request %d: %w", id, err)
	}
	return nil
}
