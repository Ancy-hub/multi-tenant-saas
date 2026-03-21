package models

import "github.com/google/uuid"

type UserOrg struct {
	OrgID uuid.UUID `json:"org_id"`
	Name  string    `json:"name"`
	Role  string    `json:"role"`
}