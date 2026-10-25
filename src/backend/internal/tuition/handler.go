package tuition

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/auth"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/httpx"
)

// Handler exposes HTTP endpoints for semester charges and simulated payments.
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Summary returns the current student's tuition. It accepts an optional
// term_id query parameter.
func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	termID, ok := httpx.QueryID(w, r, "term_id")
	if !ok {
		return
	}
	u, _ := auth.UserFromContext(r.Context())
	summary, err := h.service.Summary(r.Context(), u.ID, termID)
	if err != nil {
		h.fail(w, r, "tuition summary", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, summary)
}

// StudentSummary returns a student's tuition for staff.
func (h *Handler) StudentSummary(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.PathID(w, r)
	if !ok {
		return
	}
	termID, ok := httpx.QueryID(w, r, "term_id")
	if !ok {
		return
	}
	summary, err := h.service.Summary(r.Context(), id, termID)
	if err != nil {
		h.fail(w, r, "student tuition summary", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, summary)
}

func (h *Handler) AddCharge(w http.ResponseWriter, r *http.Request) {
	var in ChargeInput
	if err := httpx.DecodeJSON(w, r, &in); err != nil {
		httpx.BadRequest(w, "invalid JSON body")
		return
	}
	if in.StudentID <= 0 || in.TermID <= 0 {
		httpx.BadRequest(w, "student_id and term_id are required")
		return
	}
	u, _ := auth.UserFromContext(r.Context())
	charge, err := h.service.AddCharge(r.Context(), u.ID, in)
	if err != nil {
		h.fail(w, r, "add charge", err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, charge)
}

func (h *Handler) RecordPayment(w http.ResponseWriter, r *http.Request) {
	var in PaymentInput
	if err := httpx.DecodeJSON(w, r, &in); err != nil {
		httpx.BadRequest(w, "invalid JSON body")
		return
	}
	if in.StudentID <= 0 || in.TermID <= 0 {
		httpx.BadRequest(w, "student_id and term_id are required")
		return
	}
	u, _ := auth.UserFromContext(r.Context())
	payment, err := h.service.RecordPayment(r.Context(), u.ID, in)
	if err != nil {
		h.fail(w, r, "record payment", err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, payment)
}

func (h *Handler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.PathID(w, r)
	if !ok {
		return
	}
	var in PaymentUpdate
	if err := httpx.DecodeJSON(w, r, &in); err != nil {
		httpx.BadRequest(w, "invalid JSON body")
		return
	}
	u, _ := auth.UserFromContext(r.Context())
	payment, err := h.service.UpdatePayment(r.Context(), u.ID, id, in)
	if err != nil {
		h.fail(w, r, "update payment", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, payment)
}

// Payment returns one payment with its change history for staff.
func (h *Handler) Payment(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.PathID(w, r)
	if !ok {
		return
	}
	payment, err := h.service.Payment(r.Context(), id)
	if err != nil {
		h.fail(w, r, "get payment", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, payment)
}

// fail maps service errors to HTTP responses.
func (h *Handler) fail(w http.ResponseWriter, r *http.Request, action string, err error) {
	var validationErr *ValidationError
	switch {
	case errors.As(err, &validationErr):
		httpx.BadRequest(w, validationErr.Error())
	case errors.Is(err, ErrStudentNotFound), errors.Is(err, ErrTermNotFound), errors.Is(err, ErrPaymentNotFound):
		httpx.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrPaymentVoided):
		httpx.Error(w, http.StatusConflict, err.Error())
	default:
		slog.ErrorContext(r.Context(), action, "error", err)
		httpx.InternalError(w)
	}
}
