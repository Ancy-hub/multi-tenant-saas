package models

import (
	"time"

	"github.com/google/uuid"
)

type Membership struct {
	ID        uuid.UUID    `json:"id"`
	UserID    uuid.UUID    `json:"user_id"`
	OrgID     uuid.UUID    `json:"org_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}