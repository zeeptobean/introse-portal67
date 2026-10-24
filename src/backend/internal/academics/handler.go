package academics

import (
	"net/http"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/httpx"
)

// Handler exposes HTTP endpoints for course registration, scheduling, and academic results.
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListSections(w http.ResponseWriter, r *http.Request) {
	// TODO: List offered courses and sections with schedules and available seats.
	httpx.NotImplemented(w)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	// TODO: Decode the section ID and call the registration transaction.
	httpx.NotImplemented(w)
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	// TODO: Cancel the current student's enrollment.
	httpx.NotImplemented(w)
}

func (h *Handler) Timetable(w http.ResponseWriter, r *http.Request) {
	// TODO: Return meetings from the current student's active enrollments.
	httpx.NotImplemented(w)
}

func (h *Handler) Exams(w http.ResponseWriter, r *http.Request) {
	// TODO: Return the current student's exams and overlap information.
	httpx.NotImplemented(w)
}

func (h *Handler) Results(w http.ResponseWriter, r *http.Request) {
	// TODO: Return grades and the sample semester/cumulative GPA.
	httpx.NotImplemented(w)
}

func (h *Handler) SaveGrades(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate grade input and the lecturer's section assignment.
	httpx.NotImplemented(w)
}
