package app

import (
	"context"
	"crypto/ed25519"
	"net/http"
	"net/netip"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/academics"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/auth"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/documents"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/notification"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/server"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/routing"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/tuition"
)

// NewHandler wires the legacy single-process entry point.
func NewHandler(pool *pgxpool.Pool, trustedProxies []netip.Prefix, pingRedis func(context.Context) error) http.Handler {
	sessions := newSessionManager(pool)
	identity := auth.NewHandler(auth.NewService(auth.NewStore(pool)), sessions)
	mux := http.NewServeMux()
	checks := []func(context.Context) error{pool.Ping}
	if pingRedis != nil {
		checks = append(checks, pingRedis)
	}
	server.HealthRoutes(mux, checks...)
	routes := serviceRoutes{mux, identity.RequireLogin}
	routes.mountIdentity(identity)
	academic := academics.NewHandler(academics.NewService(academics.NewStore(pool)))
	routes.mountAcademic(academic)
	routes.mountRegistration(academic)
	routes.mountAffairs(pool)
	return auth.TrustClientIP(trustedProxies, sessions.LoadAndSave(requestTrace(mux)))
}

func NewDomainHandler(service string, pool *pgxpool.Pool, keys []ed25519.PublicKey, checks ...func(context.Context) error) http.Handler {
	mux := http.NewServeMux()
	server.HealthRoutes(mux, append([]func(context.Context) error{pool.Ping}, checks...)...)
	routes := serviceRoutes{mux, auth.RequireServiceIdentity(keys, service)}
	switch service {
	case routing.Academic:
		routes.mountAcademic(academics.NewHandler(academics.NewService(academics.NewStore(pool))))
	case routing.Registration:
		routes.mountRegistration(academics.NewHandler(academics.NewService(academics.NewRegistrationStore(pool))))
	case routing.Affairs:
		routes.mountAffairs(pool)
	case routing.Notification:
		routes.protected("GET /api/me/notifications", notification.Handler(pool))
	}
	return requestTrace(mux)
}

type serviceRoutes struct {
	mux   *http.ServeMux
	login func(http.Handler) http.Handler
}

func (r serviceRoutes) protected(pattern string, handler http.HandlerFunc, role ...auth.Role) {
	var endpoint http.Handler = handler
	if len(role) > 0 {
		endpoint = auth.RequireRole(role...)(endpoint)
	}
	r.mux.Handle(pattern, r.login(endpoint))
}

func (r serviceRoutes) mountIdentity(h *auth.Handler) {
	r.mux.HandleFunc("POST /api/auth/login", h.Login)
	r.protected("POST /api/auth/logout", h.Logout)
	r.protected("GET /api/me", h.Me)
	r.protected("GET /api/staff/students", h.ListStudents, auth.RoleStaff)
}

func (r serviceRoutes) mountAcademic(h *academics.Handler) {
	r.protected("GET /api/terms", h.ListTerms)
	r.protected("GET /api/me/results", h.Results, auth.RoleStudent)
	r.protected("GET /api/lecturer/sections", h.LecturerSections, auth.RoleLecturer)
	r.protected("GET /api/sections/{id}/grades", h.Roster, auth.RoleLecturer)
	r.protected("PUT /api/sections/{id}/grades", h.SaveGrades, auth.RoleLecturer)
	r.protected("GET /api/sections/{id}/grade-history", h.GradeHistory, auth.RoleLecturer)
}

func (r serviceRoutes) mountRegistration(h *academics.Handler) {
	r.protected("GET /api/sections", h.ListSections)
	r.protected("POST /api/enrollments", h.Register, auth.RoleStudent)
	r.protected("DELETE /api/enrollments/{id}", h.Cancel, auth.RoleStudent)
	r.protected("GET /api/me/enrollments", h.Enrollments, auth.RoleStudent)
	r.protected("GET /api/me/timetable", h.Timetable, auth.RoleStudent)
	r.protected("GET /api/me/exams", h.Exams, auth.RoleStudent)
}

func (r serviceRoutes) mountAffairs(pool *pgxpool.Pool) {
	fees := tuition.NewHandler(tuition.NewService(tuition.NewStore(pool)))
	r.protected("GET /api/me/tuition", fees.Summary, auth.RoleStudent)
	r.protected("GET /api/staff/students/{id}/tuition", fees.StudentSummary, auth.RoleStaff)
	r.protected("POST /api/staff/charges", fees.AddCharge, auth.RoleStaff)
	r.protected("POST /api/staff/payments", fees.RecordPayment, auth.RoleStaff)
	r.protected("GET /api/staff/payments/{id}", fees.Payment, auth.RoleStaff)
	r.protected("PATCH /api/staff/payments/{id}", fees.UpdatePayment, auth.RoleStaff)

	docs := documents.NewHandler(documents.NewService(documents.NewStore(pool)))
	r.protected("POST /api/document-requests", docs.Submit, auth.RoleStudent)
	r.protected("GET /api/document-requests", docs.List)
	r.protected("GET /api/document-requests/{id}", docs.Get)
	r.protected("PATCH /api/document-requests/{id}", docs.Decide, auth.RoleStaff)
}

func requestTrace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := events.WithTraceID(r.Context(), r.Header.Get("X-Request-ID"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
