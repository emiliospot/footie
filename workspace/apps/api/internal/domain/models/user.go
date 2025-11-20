package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system.
type User struct {
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	FirstName     string         `gorm:"not null" json:"first_name"`
	Email         string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash  string         `gorm:"not null" json:"-"`
	LastName      string         `gorm:"not null" json:"last_name"`
	Role          string         `gorm:"not null;default:'user'" json:"role"`
	Avatar        string         `json:"avatar,omitempty"`
	Organization  string         `json:"organization,omitempty"`
	ID            uint           `gorm:"primarykey" json:"id"`
	IsActive      bool           `gorm:"not null;default:true" json:"is_active"`
	EmailVerified bool           `gorm:"not null;default:false" json:"email_verified"`
}

// TableName specifies the table name for User model.
func (User) TableName() string {
	return "users"
}

// IsAdmin returns true if user is an admin.
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsAnalyst returns true if user is an analyst or admin.
func (u *User) IsAnalyst() bool {
	return u.Role == "analyst" || u.Role == "admin"
}
