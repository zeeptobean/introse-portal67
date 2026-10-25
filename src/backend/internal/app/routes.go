package app

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/academics"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/auth"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/documents"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/httpx"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/tuition"
)

// NewHandler wires shared dependencies, registers routes, and applies sessions.
func NewHandler(pool *pgxpool.Pool) http.Handler {
	sessions := newSessionManager(pool)
	authStore := auth.NewStore(pool)
	authService := auth.NewService(authStore)
	authHandler := auth.NewHandler(authService, sessions)

	academicStore := academics.NewStore(pool)
	academicService := academics.NewService(academicStore)
	academicHandler := academics.NewHandler(academicService)

	tuitionStore := tuition.NewStore(pool)
	tuitionService := tuition.NewService(tuitionStore)
	tuitionHandler := tuition.NewHandler(tuitionService)

	documentStore := documents.NewStore(pool)
	documentService := documents.NewService(documentStore)
	documentHandler := documents.NewHandler(documentService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /api/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			httpx.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
			return
		}
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})
	requireRole := func(role auth.Role, handler http.HandlerFunc) http.Handler {
		return authHandler.RequireLogin(auth.RequireRole(role)(handler))
	}

	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.Handle("POST /api/auth/logout", authHandler.RequireLogin(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("GET /api/me", authHandler.RequireLogin(http.HandlerFunc(authHandler.Me)))

	mux.Handle("GET /api/terms", authHandler.RequireLogin(http.HandlerFunc(academicHandler.ListTerms)))
	mux.Handle("GET /api/sections", authHandler.RequireLogin(http.HandlerFunc(academicHandler.ListSections)))
	mux.Handle("POST /api/enrollments", requireRole(auth.RoleStudent, academicHandler.Register))
	mux.Handle("DELETE /api/enrollments/{id}", requireRole(auth.RoleStudent, academicHandler.Cancel))
	mux.Handle("GET /api/me/enrollments", requireRole(auth.RoleStudent, academicHandler.Enrollments))
	mux.Handle("GET /api/me/timetable", requireRole(auth.RoleStudent, academicHandler.Timetable))
	mux.Handle("GET /api/me/exams", requireRole(auth.RoleStudent, academicHandler.Exams))
	mux.Handle("GET /api/me/results", requireRole(auth.RoleStudent, academicHandler.Results))
	mux.Handle("GET /api/lecturer/sections", requireRole(auth.RoleLecturer, academicHandler.LecturerSections))
	mux.Handle("GET /api/sections/{id}/grades", requireRole(auth.RoleLecturer, academicHandler.Roster))
	mux.Handle("PUT /api/sections/{id}/grades", requireRole(auth.RoleLecturer, academicHandler.SaveGrades))

	mux.Handle("GET /api/me/tuition", requireRole(auth.RoleStudent, tuitionHandler.Summary))
	mux.Handle("GET /api/staff/students", requireRole(auth.RoleStaff, authHandler.ListStudents))
	mux.Handle("GET /api/staff/students/{id}/tuition", requireRole(auth.RoleStaff, tuitionHandler.StudentSummary))
	mux.Handle("POST /api/staff/charges", requireRole(auth.RoleStaff, tuitionHandler.AddCharge))
	mux.Handle("POST /api/staff/payments", requireRole(auth.RoleStaff, tuitionHandler.RecordPayment))
	mux.Handle("PATCH /api/staff/payments/{id}", requireRole(auth.RoleStaff, tuitionHandler.UpdatePayment))

	mux.Handle("POST /api/document-requests", requireRole(auth.RoleStudent, documentHandler.Submit))
	mux.Handle("GET /api/document-requests", authHandler.RequireLogin(http.HandlerFunc(documentHandler.List)))
	mux.Handle("GET /api/document-requests/{id}", authHandler.RequireLogin(http.HandlerFunc(documentHandler.Get)))
	mux.Handle("PATCH /api/document-requests/{id}", requireRole(auth.RoleStaff, documentHandler.Decide))

	return sessions.LoadAndSave(mux)
}
