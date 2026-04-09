package models

import (
	"time"

	"github.com/google/uuid"
)

// Membership represents the relationship between a user and an organization, including their role.
type Membership struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	OrgID     uuid.UUID `json:"org_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
