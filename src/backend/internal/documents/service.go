package documents

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/auth"
)

// Service owns business rules and transaction boundaries for document requests and staff decisions.
type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

// Submit creates a pending request and its initial history event together.
func (s *Service) Submit(ctx context.Context, student auth.User, in SubmitInput) (Request, error) {
	in.Purpose = strings.TrimSpace(in.Purpose)
	if !in.RequestType.valid() {
		return Request{}, &ValidationError{Message: "request_type must be enrollment_confirmation, transcript, or other"}
	}
	if err := validateText("purpose", in.Purpose); err != nil {
		return Request{}, err
	}

	var id int64
	err := s.store.withTx(ctx, func(tx pgx.Tx) error {
		var err error
		id, err = s.store.insertRequest(ctx, tx, student.ID, in)
		if err != nil {
			return err
		}
		return s.store.insertEvent(ctx, tx, id, student.ID, nil, StatusPending, nil)
	})
	if err != nil {
		return Request{}, err
	}
	return s.withHistory(ctx, id)
}

// List returns the student's own requests, or every request for staff.
func (s *Service) List(ctx context.Context, viewer auth.User, status *Status) ([]Request, error) {
	if status != nil && !status.valid() {
		return nil, &ValidationError{Message: "status must be pending, approved, or rejected"}
	}
	switch viewer.Role {
	case auth.RoleStudent:
		return s.store.listRequests(ctx, &viewer.ID, status)
	case auth.RoleStaff:
		return s.store.listRequests(ctx, nil, status)
	default:
		return nil, ErrForbidden
	}
}

// Get returns a request with its history. Students can only read their own
// requests; other students' requests are reported as not found.
func (s *Service) Get(ctx context.Context, viewer auth.User, id int64) (Request, error) {
	if viewer.Role != auth.RoleStudent && viewer.Role != auth.RoleStaff {
		return Request{}, ErrForbidden
	}
	r, err := s.withHistory(ctx, id)
	if err != nil {
		return Request{}, err
	}
	if viewer.Role == auth.RoleStudent && r.Student.ID != viewer.ID {
		return Request{}, ErrRequestNotFound
	}
	return r, nil
}

// Decide approves or rejects a pending request. The status change and its
// history event are written in the same transaction while the request row is
// locked, so a request is decided at most once.
func (s *Service) Decide(ctx context.Context, staff auth.User, id int64, in DecisionInput) (Request, error) {
	in.Comment = strings.TrimSpace(in.Comment)
	if in.Status != StatusApproved && in.Status != StatusRejected {
		return Request{}, &ValidationError{Message: "status must be approved or rejected"}
	}
	if err := validateText("comment", in.Comment); err != nil {
		return Request{}, err
	}

	err := s.store.withTx(ctx, func(tx pgx.Tx) error {
		current, err := s.store.lockRequestStatus(ctx, tx, id)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRequestNotFound
		}
		if err != nil {
			return err
		}
		if current != StatusPending {
			return ErrAlreadyDecided
		}
		if err := s.store.updateStatus(ctx, tx, id, in.Status); err != nil {
			return err
		}
		return s.store.insertEvent(ctx, tx, id, staff.ID, &current, in.Status, &in.Comment)
	})
	if err != nil {
		return Request{}, err
	}
	return s.withHistory(ctx, id)
}

func (s *Service) withHistory(ctx context.Context, id int64) (Request, error) {
	r, err := s.store.requestByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Request{}, ErrRequestNotFound
	}
	if err != nil {
		return Request{}, err
	}
	r.History, err = s.store.history(ctx, id)
	if err != nil {
		return Request{}, err
	}
	return r, nil
}

func validateText(field, value string) error {
	if value == "" {
		return &ValidationError{Message: field + " is required"}
	}
	if utf8.RuneCountInString(value) > maxTextLength {
		return &ValidationError{Message: fmt.Sprintf("%s must be at most %d characters", field, maxTextLength)}
	}
	return nil
}
