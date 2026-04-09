package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the multi-tenant SaaS system.
type User struct {
	// ID is the unique identifier for the user.
	ID uuid.UUID `json:"id"`
	// Email is the user's email address.
	Email string `json:"email"`
	// PasswordHash is the hashed password for the user.
	PasswordHash string `json:"-"` //hide from API
	// Name is the user's full name.
	Name string `json:"name"`
	// CreatedAt is the timestamp when the user was created.
	CreatedAt time.Time `json:"created_at"`
}
