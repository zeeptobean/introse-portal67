package auth

import "time"

type Role string

const (
	RoleStudent  Role = "student"
	RoleLecturer Role = "lecturer"
	RoleStaff    Role = "staff"
)

type User struct {
	ID              int64     `json:"id,string"`
	Email           string    `json:"email"`
	InstitutionalID *string   `json:"institutional_id,omitempty"`
	FullName        string    `json:"full_name"`
	DateOfBirth     time.Time `json:"date_of_birth"`
	Role            Role      `json:"role"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}

type UserUpdate struct {
	InstitutionalID *string    `json:"institutional_id,omitempty"`
	FullName        *string    `json:"full_name"`
	DateOfBirth     *time.Time `json:"date_of_birth"`
}
