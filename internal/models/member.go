package models

import "github.com/google/uuid"

// Member represents a member of an organization for API responses.
// This is used for API responses, while membership.go is for the database table.
type Member struct {
	// UserID is the ID of the user.
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
}
