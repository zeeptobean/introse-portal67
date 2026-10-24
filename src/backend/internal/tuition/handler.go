package tuition

import (
	"net/http"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/httpx"
)

// Handler exposes HTTP endpoints for semester charges and simulated payments.
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	// TODO: Return the current student's charges, payments, and balance.
	httpx.NotImplemented(w)
}

func (h *Handler) RecordPayment(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate and record a simulated payment as staff.
	httpx.NotImplemented(w)
}

func (h *Handler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	// TODO: Update or void a simulated payment as staff.
	httpx.NotImplemented(w)
}
