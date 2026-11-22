package academics

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/auth"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/httpx"
)

// Handler exposes HTTP endpoints for course registration, scheduling, and academic results.
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListTerms(w http.ResponseWriter, r *http.Request) {
	terms, err := h.service.Terms(r.Context())
	if err != nil {
		h.fail(w, r, "list terms", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, terms)
}

// ListSections accepts the optional query parameters term_id, kind
// (course or specialized_topic), and q (course code or name).
func (h *Handler) ListSections(w http.ResponseWriter, r *http.Request) {
	u, _ := auth.UserFromContext(r.Context())
	termID, ok := httpx.QueryID(w, r, "term_id")
	if !ok {
		return
	}
	f := SectionFilter{TermID: termID, Query: r.URL.Query().Get("q")}
	if k := CourseKind(r.URL.Query().Get("kind")); k != "" {
		if k != KindCourse && k != KindSpecializedTopic {
			httpx.BadRequest(w, "kind must be course or specialized_topic")
			return
		}
		f.Kind = &k
	}
	var studentID int64
	if u.Role == auth.RoleStudent {
		studentID = u.ID
	}
	sections, err := h.service.Sections(r.Context(), f, studentID)
	if err != nil {
		h.fail(w, r, "list sections", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, sections)
}

type registerRequest struct {
	SectionID int64 `json:"section_id,string"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var in registerRequest
	if err := httpx.DecodeJSON(w, r, &in); err != nil {
		httpx.BadRequest(w, "invalid JSON body")
		return
	}
	if in.SectionID <= 0 {
		httpx.BadRequest(w, "section_id is required")
		return
	}
	u, _ := auth.UserFromContext(r.Context())
	e, err := h.service.Register(r.Context(), u.ID, in.SectionID)
	if err != nil {
		h.fail(w, r, "register", err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, e)
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.PathID(w, r)
	if !ok {
		return
	}
	u, _ := auth.UserFromContext(r.Context())
	e, err := h.service.Cancel(r.Context(), u.ID, id)
	if err != nil {
		h.fail(w, r, "cancel enrollment", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, e)
}

func (h *Handler) Enrollments(w http.ResponseWriter, r *http.Request) {
	termID, ok := httpx.QueryID(w, r, "term_id")
	if !ok {
		return
	}
	u, _ := auth.UserFromContext(r.Context())
	enrollments, err := h.service.Enrollments(r.Context(), u.ID, termID)
	if err != nil {
		h.fail(w, r, "list enrollments", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, enrollments)
}

func (h *Handler) Timetable(w http.ResponseWriter, r *http.Request) {
	termID, ok := httpx.QueryID(w, r, "term_id")
	if !ok {
		return
	}
	u, _ := auth.UserFromContext(r.Context())
	entries, err := h.service.Timetable(r.Context(), u.ID, termID)
	if err != nil {
		h.fail(w, r, "list timetable", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, entries)
}

func (h *Handler) Exams(w http.ResponseWriter, r *http.Request) {
	termID, ok := httpx.QueryID(w, r, "term_id")
	if !ok {
		return
	}
	u, _ := auth.UserFromContext(r.Context())
	exams, err := h.service.Exams(r.Context(), u.ID, termID)
	if err != nil {
		h.fail(w, r, "list exams", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, exams)
}

func (h *Handler) Results(w http.ResponseWriter, r *http.Request) {
	u, _ := auth.UserFromContext(r.Context())
	results, err := h.service.Results(r.Context(), u.ID)
	if err != nil {
		h.fail(w, r, "list results", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, results)
}

func (h *Handler) LecturerSections(w http.ResponseWriter, r *http.Request) {
	u, _ := auth.UserFromContext(r.Context())
	sections, err := h.service.LecturerSections(r.Context(), u.ID)
	if err != nil {
		h.fail(w, r, "list lecturer sections", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, sections)
}

func (h *Handler) Roster(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.PathID(w, r)
	if !ok {
		return
	}
	u, _ := auth.UserFromContext(r.Context())
	roster, err := h.service.Roster(r.Context(), u.ID, id)
	if err != nil {
		h.fail(w, r, "list roster", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, roster)
}

func (h *Handler) GradeHistory(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.PathID(w, r)
	if !ok {
		return
	}
	u, _ := auth.UserFromContext(r.Context())
	history, err := h.service.GradeHistory(r.Context(), u.ID, id)
	if err != nil {
		h.fail(w, r, "list grade history", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, history)
}

type saveGradesRequest struct {
	Grades []GradeInput `json:"grades"`
}

func (h *Handler) SaveGrades(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.PathID(w, r)
	if !ok {
		return
	}
	var in saveGradesRequest
	if err := httpx.DecodeJSON(w, r, &in); err != nil {
		httpx.BadRequest(w, "invalid JSON body")
		return
	}
	u, _ := auth.UserFromContext(r.Context())
	roster, err := h.service.SaveGrades(r.Context(), u.ID, id, in.Grades)
	if err != nil {
		h.fail(w, r, "save grades", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, roster)
}

// fail maps service errors to HTTP responses.
func (h *Handler) fail(w http.ResponseWriter, r *http.Request, action string, err error) {
	var conflictErr *ScheduleConflictError
	var gradeErr *GradeInputError
	switch {
	case errors.As(err, &conflictErr):
		httpx.WriteJSON(w, http.StatusConflict, map[string]any{
			"error":     conflictErr.Error(),
			"conflicts": conflictErr.Conflicts,
		})
	case errors.As(err, &gradeErr):
		httpx.BadRequest(w, gradeErr.Error())
	case errors.Is(err, ErrSectionNotFound), errors.Is(err, ErrEnrollmentNotFound):
		httpx.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrNotStudent), errors.Is(err, ErrNotAssigned):
		httpx.Error(w, http.StatusForbidden, err.Error())
	case errors.Is(err, ErrAlreadyEnrolled), errors.Is(err, ErrCourseTaken), errors.Is(err, ErrSectionFull),
		errors.Is(err, ErrRegistrationNotOpen), errors.Is(err, ErrRegistrationClosed), errors.Is(err, ErrGradingClosed),
		errors.Is(err, ErrAlreadyCancelled), errors.Is(err, ErrAlreadyGraded):
		httpx.Error(w, http.StatusConflict, err.Error())
	default:
		slog.ErrorContext(r.Context(), action, "error", err)
		httpx.InternalError(w)
	}
}
