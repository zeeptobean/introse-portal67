package documents

import (
	"errors"
	"time"
)

type RequestType string

const (
	TypeEnrollmentConfirmation RequestType = "enrollment_confirmation"
	TypeTranscript             RequestType = "transcript"
	TypeOther                  RequestType = "other"
)

func (t RequestType) valid() bool {
	return t == TypeEnrollmentConfirmation || t == TypeTranscript || t == TypeOther
}

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

func (s Status) valid() bool {
	return s == StatusPending || s == StatusApproved || s == StatusRejected
}

const maxTextLength = 1000

var (
	ErrRequestNotFound = errors.New("document request not found")
	ErrAlreadyDecided  = errors.New("document request has already been decided")
	ErrForbidden       = errors.New("permission denied")
)

// ValidationError reports invalid request input.
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

// Request is a document request with its decision, if one has been made.
// History is only filled when a single request is loaded.
type Request struct {
	ID              int64       `json:"id,string"`
	Student         Student     `json:"student"`
	RequestType     RequestType `json:"request_type"`
	Purpose         string      `json:"purpose"`
	Status          Status      `json:"status"`
	SubmittedAt     time.Time   `json:"submitted_at"`
	DecidedAt       *time.Time  `json:"decided_at"`
	DecisionComment *string     `json:"decision_comment"`
	History         []Event     `json:"history,omitempty"`
}

type Actor struct {
	ID       int64  `json:"id,string"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// Event is one status change in a request's history.
type Event struct {
	ID         int64     `json:"id,string"`
	Actor      Actor     `json:"actor"`
	FromStatus *Status   `json:"from_status"`
	ToStatus   Status    `json:"to_status"`
	Comment    *string   `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
}

type SubmitInput struct {
	RequestType RequestType `json:"request_type"`
	Purpose     string      `json:"purpose"`
}

type DecisionInput struct {
	Status  Status `json:"status"`
	Comment string `json:"comment"`
}
