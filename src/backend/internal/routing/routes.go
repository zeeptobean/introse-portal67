package routing

const (
	Identity     = "identity"
	Academic     = "academic"
	Registration = "registration"
	Affairs      = "affairs"
	Notification = "notification"
)

// Routes declares the public API ownership used by the gateway.
var Routes = map[string][]string{
	Identity: {
		"POST /api/auth/login", "POST /api/auth/logout", "GET /api/me", "GET /api/staff/students",
	},
	Academic: {
		"GET /api/terms", "GET /api/me/results", "GET /api/lecturer/sections",
		"GET /api/sections/{id}/grades", "PUT /api/sections/{id}/grades", "GET /api/sections/{id}/grade-history",
	},
	Registration: {
		"GET /api/sections", "POST /api/enrollments", "DELETE /api/enrollments/{id}",
		"GET /api/me/enrollments", "GET /api/me/timetable", "GET /api/me/exams",
	},
	Affairs: {
		"GET /api/me/tuition", "GET /api/staff/students/{id}/tuition", "POST /api/staff/charges",
		"POST /api/staff/payments", "GET /api/staff/payments/{id}", "PATCH /api/staff/payments/{id}",
		"GET /api/document-requests", "POST /api/document-requests", "GET /api/document-requests/{id}", "PATCH /api/document-requests/{id}",
	},
	Notification: {"GET /api/me/notifications"},
}
