package documents

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/auth"
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
	var in SubmitInput
	if err := httpx.DecodeJSON(w, r, &in); err != nil {
		httpx.BadRequest(w, "invalid JSON body")
		return
	}
	u, _ := auth.UserFromContext(r.Context())
	req, err := h.service.Submit(r.Context(), u, in)
	if err != nil {
		h.fail(w, r, "submit document request", err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, req)
}

// List accepts an optional status query parameter.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	var status *Status
	if raw := r.URL.Query().Get("status"); raw != "" {
		st := Status(raw)
		status = &st
	}
	u, _ := auth.UserFromContext(r.Context())
	requests, err := h.service.List(r.Context(), u, status)
	if err != nil {
		h.fail(w, r, "list document requests", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, requests)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.PathID(w, r)
	if !ok {
		return
	}
	u, _ := auth.UserFromContext(r.Context())
	req, err := h.service.Get(r.Context(), u, id)
	if err != nil {
		h.fail(w, r, "get document request", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, req)
}

func (h *Handler) Decide(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.PathID(w, r)
	if !ok {
		return
	}
	var in DecisionInput
	if err := httpx.DecodeJSON(w, r, &in); err != nil {
		httpx.BadRequest(w, "invalid JSON body")
		return
	}
	u, _ := auth.UserFromContext(r.Context())
	req, err := h.service.Decide(r.Context(), u, id, in)
	if err != nil {
		h.fail(w, r, "decide document request", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, req)
}

// fail maps service errors to HTTP responses.
func (h *Handler) fail(w http.ResponseWriter, r *http.Request, action string, err error) {
	var validationErr *ValidationError
	switch {
	case errors.As(err, &validationErr):
		httpx.BadRequest(w, validationErr.Error())
	case errors.Is(err, ErrRequestNotFound):
		httpx.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrAlreadyDecided):
		httpx.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, ErrForbidden):
		httpx.PermissionDenied(w)
	default:
		slog.ErrorContext(r.Context(), action, "error", err)
		httpx.InternalError(w)
	}
}
