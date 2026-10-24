package documents

import (
	"net/http"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/httpx"
)

// Handler exposes HTTP endpoints for document requests and staff decisions.
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Submit(w http.ResponseWriter, r *http.Request) {
	// TODO: Create a pending student request and its initial history event.
	httpx.NotImplemented(w)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: List the student's own requests or the staff work queue.
	httpx.NotImplemented(w)
}

func (h *Handler) Decide(w http.ResponseWriter, r *http.Request) {
	// TODO: Record a staff decision and append history in the same transaction.
	httpx.NotImplemented(w)
}
