package models

import "github.com/google/uuid"

// UserOrg represents the organizations a user belongs to, used in API responses.
type UserOrg struct {
	// OrgID is the ID of the organization.
	OrgID uuid.UUID `json:"org_id"`
	Name  string    `json:"name"`
	Role  string    `json:"role"`
}
