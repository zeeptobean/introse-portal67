package tuition

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/term"
)

// Service owns business rules and transaction boundaries for semester charges and simulated payments.
type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

// Summary returns the student's charges and payments grouped by term, newest
// term first. Charges and recorded payments are summed separately, so a
// balance never depends on how the two lists would join.
func (s *Service) Summary(ctx context.Context, studentID int64, termID *int64) (Summary, error) {
	student, err := s.store.student(ctx, studentID)
	if errors.Is(err, pgx.ErrNoRows) {
		return Summary{}, ErrStudentNotFound
	}
	if err != nil {
		return Summary{}, err
	}
	charges, err := s.store.charges(ctx, studentID, termID)
	if err != nil {
		return Summary{}, err
	}
	payments, err := s.store.payments(ctx, studentID, termID)
	if err != nil {
		return Summary{}, err
	}

	byTerm := make(map[int64]*TermSummary)
	termFor := func(t term.Term) *TermSummary {
		ts, ok := byTerm[t.ID]
		if !ok {
			ts = &TermSummary{Term: t, Charges: []Charge{}, Payments: []Payment{}}
			byTerm[t.ID] = ts
		}
		return ts
	}
	for _, c := range charges {
		ts := termFor(c.Term)
		ts.Charges = append(ts.Charges, c)
		ts.ChargedVND += c.AmountVND
	}
	for _, p := range payments {
		ts := termFor(p.Term)
		ts.Payments = append(ts.Payments, p)
		if p.Status == PaymentRecorded {
			ts.PaidVND += p.AmountVND
		}
	}

	out := Summary{Student: student, Terms: make([]TermSummary, 0, len(byTerm))}
	for _, ts := range byTerm {
		ts.Totals = totals(ts.ChargedVND, ts.PaidVND)
		out.Terms = append(out.Terms, *ts)
		out.ChargedVND += ts.ChargedVND
		out.PaidVND += ts.PaidVND
	}
	slices.SortFunc(out.Terms, func(a, b TermSummary) int {
		return cmp.Or(cmp.Compare(b.Term.StartsOn, a.Term.StartsOn), cmp.Compare(b.Term.ID, a.Term.ID))
	})
	out.Totals = totals(out.ChargedVND, out.PaidVND)
	return out, nil
}

func totals(charged, paid int64) Totals {
	t := Totals{ChargedVND: charged, PaidVND: paid, BalanceVND: charged - paid}
	switch {
	case charged == 0 && paid == 0:
		t.Status = BalanceNoCharges
	case t.BalanceVND < 0:
		t.Status = BalanceCredit
	case t.BalanceVND == 0:
		t.Status = BalancePaid
	case paid == 0:
		t.Status = BalanceUnpaid
	default:
		t.Status = BalancePartial
	}
	return t
}

// AddCharge records a charge line for a student's term.
func (s *Service) AddCharge(ctx context.Context, staffID int64, in ChargeInput) (Charge, error) {
	in.Description = strings.TrimSpace(in.Description)
	if in.Description == "" {
		return Charge{}, &ValidationError{Message: "description is required"}
	}
	if err := validateText("description", in.Description); err != nil {
		return Charge{}, err
	}
	if err := validateAmount(in.AmountVND); err != nil {
		return Charge{}, err
	}

	var id int64
	err := s.store.withTx(ctx, func(tx pgx.Tx) error {
		if err := s.checkStudentAndTerm(ctx, tx, in.StudentID, in.TermID); err != nil {
			return err
		}
		var err error
		id, err = s.store.insertCharge(ctx, tx, staffID, in)
		return err
	})
	if err != nil {
		return Charge{}, err
	}
	return s.store.chargeByID(ctx, id)
}

// RecordPayment records a simulated payment. paid_at defaults to now and may
// not be in the future.
func (s *Service) RecordPayment(ctx context.Context, staffID int64, in PaymentInput) (Payment, error) {
	if err := validateAmount(in.AmountVND); err != nil {
		return Payment{}, err
	}
	if err := validatePaidAt(in.PaidAt); err != nil {
		return Payment{}, err
	}
	note, err := normalizeNote(in.Note)
	if err != nil {
		return Payment{}, err
	}
	in.Note = note

	var id int64
	err = s.store.withTx(ctx, func(tx pgx.Tx) error {
		if err := s.checkStudentAndTerm(ctx, tx, in.StudentID, in.TermID); err != nil {
			return err
		}
		var err error
		id, err = s.store.insertPayment(ctx, tx, staffID, in)
		if err != nil {
			return err
		}
		return s.store.insertPaymentEvent(ctx, tx, id, staffID, ActionRecorded)
	})
	if err != nil {
		return Payment{}, err
	}
	return s.Payment(ctx, id)
}

// Payment returns a payment with its change history.
func (s *Service) Payment(ctx context.Context, id int64) (Payment, error) {
	p, err := s.store.paymentByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Payment{}, ErrPaymentNotFound
	}
	if err != nil {
		return Payment{}, err
	}
	p.History, err = s.store.paymentHistory(ctx, id)
	if err != nil {
		return Payment{}, err
	}
	return p, nil
}

// UpdatePayment corrects or voids a recorded payment on behalf of staffID.
// Every change is appended to the payment's history; a request that changes
// nothing records nothing. Voided payments are kept for history and cannot be
// changed again.
func (s *Service) UpdatePayment(ctx context.Context, staffID, id int64, in PaymentUpdate) (Payment, error) {
	if in.AmountVND == nil && in.PaidAt == nil && in.Note == nil && in.Status == nil {
		return Payment{}, &ValidationError{Message: "no fields to update"}
	}
	if in.AmountVND != nil {
		if err := validateAmount(*in.AmountVND); err != nil {
			return Payment{}, err
		}
	}
	if err := validatePaidAt(in.PaidAt); err != nil {
		return Payment{}, err
	}
	if in.Status != nil && *in.Status != PaymentRecorded && *in.Status != PaymentVoided {
		return Payment{}, &ValidationError{Message: "status must be recorded or voided"}
	}
	note, err := normalizeNote(in.Note)
	if err != nil {
		return Payment{}, err
	}
	clearNote := in.Note != nil && note == nil
	in.Note = note

	err = s.store.withTx(ctx, func(tx pgx.Tx) error {
		status, err := s.store.lockPaymentStatus(ctx, tx, id)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPaymentNotFound
		}
		if err != nil {
			return err
		}
		if status == PaymentVoided {
			return ErrPaymentVoided
		}
		changed, err := s.store.updatePayment(ctx, tx, staffID, id, in, clearNote)
		if err != nil || !changed {
			return err
		}
		action := ActionUpdated
		if in.Status != nil && *in.Status == PaymentVoided {
			action = ActionVoided
		}
		return s.store.insertPaymentEvent(ctx, tx, id, staffID, action)
	})
	if err != nil {
		return Payment{}, err
	}
	return s.Payment(ctx, id)
}

func (s *Service) checkStudentAndTerm(ctx context.Context, tx pgx.Tx, studentID, termID int64) error {
	isStudent, err := s.store.isStudent(ctx, tx, studentID)
	if err != nil {
		return err
	}
	if !isStudent {
		return ErrStudentNotFound
	}
	exists, err := s.store.termExists(ctx, tx, termID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrTermNotFound
	}
	return nil
}

func validateAmount(amount int64) error {
	if amount <= 0 || amount > maxAmountVND {
		return &ValidationError{Message: fmt.Sprintf("amount_vnd must be between 1 and %d", maxAmountVND)}
	}
	return nil
}

func validatePaidAt(paidAt *time.Time) error {
	if paidAt != nil && paidAt.After(time.Now().Add(time.Minute)) {
		return &ValidationError{Message: "paid_at cannot be in the future"}
	}
	return nil
}

func validateText(field, value string) error {
	if utf8.RuneCountInString(value) > maxTextLength {
		return &ValidationError{Message: fmt.Sprintf("%s must be at most %d characters", field, maxTextLength)}
	}
	return nil
}

// normalizeNote trims the note and turns a blank note into nil.
func normalizeNote(note *string) (*string, error) {
	if note == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(*note)
	if trimmed == "" {
		return nil, nil
	}
	if err := validateText("note", trimmed); err != nil {
		return nil, err
	}
	return &trimmed, nil
}
