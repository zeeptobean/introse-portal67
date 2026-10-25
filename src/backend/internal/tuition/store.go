package tuition

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/dbx"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/term"
)

// Store holds the database dependency for semester charges and simulated payments.
type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) withTx(ctx context.Context, fn func(pgx.Tx) error) error {
	return dbx.WithTx(ctx, s.db, fn)
}

// student returns the user with the student role, or pgx.ErrNoRows (wrapped)
// when there is none.
func (s *Store) student(ctx context.Context, id int64) (Student, error) {
	var st Student
	err := s.db.QueryRow(ctx, `
		SELECT id, institutional_id, full_name, email FROM users
		WHERE id = $1 AND role = 'student'`, id).
		Scan(&st.ID, &st.InstitutionalID, &st.FullName, &st.Email)
	if err != nil {
		return Student{}, fmt.Errorf("find student %d: %w", id, err)
	}
	return st, nil
}

func (s *Store) termExists(ctx context.Context, tx pgx.Tx, id int64) (bool, error) {
	var exists bool
	if err := tx.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM terms WHERE id = $1)", id).Scan(&exists); err != nil {
		return false, fmt.Errorf("check term %d: %w", id, err)
	}
	return exists, nil
}

func (s *Store) isStudent(ctx context.Context, tx pgx.Tx, id int64) (bool, error) {
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE id = $1 AND role = 'student')", id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check student %d: %w", id, err)
	}
	return exists, nil
}

const chargeQuery = `
	SELECT ch.id, ch.student_id, ` + term.Columns + `, ch.description, ch.amount_vnd,
		u.id, u.full_name, ch.created_at
	FROM tuition_charges ch
	JOIN terms t ON t.id = ch.term_id
	JOIN users u ON u.id = ch.created_by`

func scanCharge(row pgx.CollectableRow) (Charge, error) {
	var c Charge
	dest := append([]any{&c.ID, &c.StudentID}, term.Dest(&c.Term)...)
	dest = append(dest, &c.Description, &c.AmountVND, &c.CreatedBy.ID, &c.CreatedBy.FullName, &c.CreatedAt)
	err := row.Scan(dest...)
	return c, err
}

func (s *Store) charges(ctx context.Context, studentID int64, termID *int64) ([]Charge, error) {
	rows, err := s.db.Query(ctx, chargeQuery+`
		WHERE ch.student_id = $1 AND ($2::bigint IS NULL OR ch.term_id = $2)
		ORDER BY t.starts_on DESC, ch.created_at, ch.id`, studentID, termID)
	if err != nil {
		return nil, fmt.Errorf("list charges: %w", err)
	}
	charges, err := pgx.CollectRows(rows, scanCharge)
	if err != nil {
		return nil, fmt.Errorf("list charges: %w", err)
	}
	return charges, nil
}

func (s *Store) chargeByID(ctx context.Context, id int64) (Charge, error) {
	rows, err := s.db.Query(ctx, chargeQuery+" WHERE ch.id = $1", id)
	if err != nil {
		return Charge{}, fmt.Errorf("find charge %d: %w", id, err)
	}
	c, err := pgx.CollectExactlyOneRow(rows, scanCharge)
	if err != nil {
		return Charge{}, fmt.Errorf("find charge %d: %w", id, err)
	}
	return c, nil
}

func (s *Store) insertCharge(ctx context.Context, tx pgx.Tx, staffID int64, in ChargeInput) (int64, error) {
	var id int64
	err := tx.QueryRow(ctx, `
		INSERT INTO tuition_charges (student_id, term_id, description, amount_vnd, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`, in.StudentID, in.TermID, in.Description, in.AmountVND, staffID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert charge: %w", err)
	}
	return id, nil
}

const paymentQuery = `
	SELECT p.id, p.student_id, ` + term.Columns + `, p.amount_vnd, p.status, p.paid_at, p.note,
		rb.id, rb.full_name, ub.id, ub.full_name, p.updated_at
	FROM tuition_payments p
	JOIN terms t ON t.id = p.term_id
	JOIN users rb ON rb.id = p.recorded_by
	JOIN users ub ON ub.id = p.updated_by`

func scanPayment(row pgx.CollectableRow) (Payment, error) {
	var p Payment
	dest := append([]any{&p.ID, &p.StudentID}, term.Dest(&p.Term)...)
	dest = append(dest, &p.AmountVND, &p.Status, &p.PaidAt, &p.Note,
		&p.RecordedBy.ID, &p.RecordedBy.FullName, &p.UpdatedBy.ID, &p.UpdatedBy.FullName, &p.UpdatedAt)
	err := row.Scan(dest...)
	return p, err
}

func (s *Store) payments(ctx context.Context, studentID int64, termID *int64) ([]Payment, error) {
	rows, err := s.db.Query(ctx, paymentQuery+`
		WHERE p.student_id = $1 AND ($2::bigint IS NULL OR p.term_id = $2)
		ORDER BY t.starts_on DESC, p.paid_at, p.id`, studentID, termID)
	if err != nil {
		return nil, fmt.Errorf("list payments: %w", err)
	}
	payments, err := pgx.CollectRows(rows, scanPayment)
	if err != nil {
		return nil, fmt.Errorf("list payments: %w", err)
	}
	return payments, nil
}

// paymentByID returns pgx.ErrNoRows (wrapped) when the payment does not exist.
func (s *Store) paymentByID(ctx context.Context, id int64) (Payment, error) {
	rows, err := s.db.Query(ctx, paymentQuery+" WHERE p.id = $1", id)
	if err != nil {
		return Payment{}, fmt.Errorf("find payment %d: %w", id, err)
	}
	p, err := pgx.CollectExactlyOneRow(rows, scanPayment)
	if err != nil {
		return Payment{}, fmt.Errorf("find payment %d: %w", id, err)
	}
	return p, nil
}

func (s *Store) paymentHistory(ctx context.Context, paymentID int64) ([]PaymentEvent, error) {
	rows, err := s.db.Query(ctx, `
		SELECT e.id, u.id, u.full_name, e.action, e.amount_vnd, e.status, e.paid_at, e.note, e.created_at
		FROM tuition_payment_events e
		JOIN users u ON u.id = e.actor_id
		WHERE e.payment_id = $1
		ORDER BY e.created_at, e.id`, paymentID)
	if err != nil {
		return nil, fmt.Errorf("list history of payment %d: %w", paymentID, err)
	}
	events, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (PaymentEvent, error) {
		var e PaymentEvent
		err := row.Scan(&e.ID, &e.Actor.ID, &e.Actor.FullName, &e.Action, &e.AmountVND, &e.Status,
			&e.PaidAt, &e.Note, &e.CreatedAt)
		return e, err
	})
	if err != nil {
		return nil, fmt.Errorf("list history of payment %d: %w", paymentID, err)
	}
	return events, nil
}

func (s *Store) insertPayment(ctx context.Context, tx pgx.Tx, staffID int64, in PaymentInput) (int64, error) {
	var id int64
	err := tx.QueryRow(ctx, `
		INSERT INTO tuition_payments (student_id, term_id, amount_vnd, recorded_by, updated_by, paid_at, note)
		VALUES ($1, $2, $3, $4, $4, COALESCE($5::timestamptz, now()), $6)
		RETURNING id`, in.StudentID, in.TermID, in.AmountVND, staffID, in.PaidAt, in.Note).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert payment: %w", err)
	}
	return id, nil
}

// insertPaymentEvent records the payment's current state as changed by actorID.
func (s *Store) insertPaymentEvent(ctx context.Context, tx pgx.Tx, paymentID, actorID int64, action PaymentAction) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO tuition_payment_events (payment_id, actor_id, action, amount_vnd, status, paid_at, note)
		SELECT id, $2, $3, amount_vnd, status, paid_at, note
		FROM tuition_payments
		WHERE id = $1`, paymentID, actorID, action)
	if err != nil {
		return fmt.Errorf("insert event for payment %d: %w", paymentID, err)
	}
	return nil
}

// lockPaymentStatus locks the payment row and returns its current status.
func (s *Store) lockPaymentStatus(ctx context.Context, tx pgx.Tx, id int64) (PaymentStatus, error) {
	var status PaymentStatus
	err := tx.QueryRow(ctx, "SELECT status FROM tuition_payments WHERE id = $1 FOR UPDATE", id).Scan(&status)
	if err != nil {
		return "", fmt.Errorf("lock payment %d: %w", id, err)
	}
	return status, nil
}

// updatePayment applies the present fields on behalf of staffID. clearNote
// removes the note. It reports whether any value actually changed; an update
// that changes nothing leaves the row, including updated_by, untouched.
func (s *Store) updatePayment(ctx context.Context, tx pgx.Tx, staffID, id int64, in PaymentUpdate, clearNote bool) (bool, error) {
	tag, err := tx.Exec(ctx, `
		WITH next AS (
			SELECT id,
				COALESCE($2::bigint, amount_vnd) AS amount_vnd,
				COALESCE($3::timestamptz, paid_at) AS paid_at,
				CASE WHEN $4::boolean THEN NULL ELSE COALESCE($5::text, note) END AS note,
				COALESCE($6::text, status) AS status
			FROM tuition_payments
			WHERE id = $1
		)
		UPDATE tuition_payments p
		SET amount_vnd = n.amount_vnd, paid_at = n.paid_at, note = n.note, status = n.status,
			updated_by = $7, updated_at = now()
		FROM next n
		WHERE p.id = n.id
			AND (p.amount_vnd, p.paid_at, p.note, p.status) IS DISTINCT FROM (n.amount_vnd, n.paid_at, n.note, n.status)`,
		id, in.AmountVND, in.PaidAt, clearNote, in.Note, in.Status, staffID)
	if err != nil {
		return false, fmt.Errorf("update payment %d: %w", id, err)
	}
	return tag.RowsAffected() > 0, nil
}
