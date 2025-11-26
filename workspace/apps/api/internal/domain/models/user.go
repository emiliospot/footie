package models

import (
	"time"
)

// User represents a user in the system.
// This is a domain model - database-agnostic, contains business logic.
type User struct {
	ID            int32     `json:"id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"` // Never expose password hash in JSON
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Role          string    `json:"role"`
	Avatar        *string   `json:"avatar,omitempty"`
	Organization  *string   `json:"organization,omitempty"`
	IsActive      bool      `json:"is_active"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"-"` // Soft delete timestamp
}

// IsAdmin returns true if user is an admin.
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsAnalyst returns true if user is an analyst or admin.
func (u *User) IsAnalyst() bool {
	return u.Role == "analyst" || u.Role == "admin"
}
