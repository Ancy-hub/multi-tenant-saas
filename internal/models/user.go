package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` //hide from API
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
}