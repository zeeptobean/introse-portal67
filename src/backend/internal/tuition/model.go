package tuition

import (
	"errors"
	"time"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/term"
)

type PaymentStatus string

const (
	PaymentRecorded PaymentStatus = "recorded"
	PaymentVoided   PaymentStatus = "voided"
)

// BalanceStatus is derived from the charge and payment totals; it is never stored.
type BalanceStatus string

const (
	BalanceNoCharges BalanceStatus = "no_charges"
	BalanceUnpaid    BalanceStatus = "unpaid"
	BalancePartial   BalanceStatus = "partial"
	BalancePaid      BalanceStatus = "paid"
	BalanceCredit    BalanceStatus = "credit"
)

// maxAmountVND bounds a single charge or payment to catch typing mistakes.
const maxAmountVND = 1_000_000_000

const maxTextLength = 500

var (
	ErrStudentNotFound = errors.New("student not found")
	ErrTermNotFound    = errors.New("term not found")
	ErrPaymentNotFound = errors.New("payment not found")
	ErrPaymentVoided   = errors.New("a voided payment cannot be changed")
)

// ValidationError reports invalid charge or payment input.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

type Student struct {
	ID              int64   `json:"id,string"`
	InstitutionalID *string `json:"institutional_id"`
	FullName        string  `json:"full_name"`
	Email           string  `json:"email"`
}

type Staff struct {
	ID       int64  `json:"id,string"`
	FullName string `json:"full_name"`
}

type Charge struct {
	ID          int64     `json:"id,string"`
	StudentID   int64     `json:"student_id,string"`
	Term        term.Term `json:"term"`
	Description string    `json:"description"`
	AmountVND   int64     `json:"amount_vnd"`
	CreatedBy   Staff     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// Payment is a simulated payment. History is only filled when a single
// payment is loaded.
type Payment struct {
	ID         int64          `json:"id,string"`
	StudentID  int64          `json:"student_id,string"`
	Term       term.Term      `json:"term"`
	AmountVND  int64          `json:"amount_vnd"`
	Status     PaymentStatus  `json:"status"`
	PaidAt     time.Time      `json:"paid_at"`
	Note       *string        `json:"note"`
	RecordedBy Staff          `json:"recorded_by"`
	UpdatedBy  Staff          `json:"updated_by"`
	UpdatedAt  time.Time      `json:"updated_at"`
	History    []PaymentEvent `json:"history,omitempty"`
}

// PaymentAction names the change recorded by a PaymentEvent. ActionSnapshot
// is the state of a payment when history tracking began.
type PaymentAction string

const (
	ActionSnapshot PaymentAction = "snapshot"
	ActionRecorded PaymentAction = "recorded"
	ActionUpdated  PaymentAction = "updated"
	ActionVoided   PaymentAction = "voided"
)

// PaymentEvent is the payment as it was right after one change.
type PaymentEvent struct {
	ID        int64         `json:"id,string"`
	Actor     Staff         `json:"actor"`
	Action    PaymentAction `json:"action"`
	AmountVND int64         `json:"amount_vnd"`
	Status    PaymentStatus `json:"status"`
	PaidAt    time.Time     `json:"paid_at"`
	Note      *string       `json:"note"`
	CreatedAt time.Time     `json:"created_at"`
}

// Totals summarizes charges against recorded payments. A negative balance is
// a credit in the student's favour.
type Totals struct {
	ChargedVND int64         `json:"charged_vnd"`
	PaidVND    int64         `json:"paid_vnd"`
	BalanceVND int64         `json:"balance_vnd"`
	Status     BalanceStatus `json:"status"`
}

type TermSummary struct {
	Term     term.Term `json:"term"`
	Charges  []Charge  `json:"charges"`
	Payments []Payment `json:"payments"`
	Totals
}

type Summary struct {
	Student Student       `json:"student"`
	Terms   []TermSummary `json:"terms"`
	Totals
}

type ChargeInput struct {
	StudentID   int64  `json:"student_id,string"`
	TermID      int64  `json:"term_id,string"`
	Description string `json:"description"`
	AmountVND   int64  `json:"amount_vnd"`
}

type PaymentInput struct {
	StudentID int64      `json:"student_id,string"`
	TermID    int64      `json:"term_id,string"`
	AmountVND int64      `json:"amount_vnd"`
	PaidAt    *time.Time `json:"paid_at"`
	Note      *string    `json:"note"`
}

// PaymentUpdate changes only the fields that are present. Setting status to
// voided removes the payment from the paid total.
type PaymentUpdate struct {
	AmountVND *int64         `json:"amount_vnd"`
	PaidAt    *time.Time     `json:"paid_at"`
	Note      *string        `json:"note"`
	Status    *PaymentStatus `json:"status"`
}
