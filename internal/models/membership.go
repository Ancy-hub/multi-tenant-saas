package models

import "time"

type Membership struct{
	ID string
	UserID string
	OrganizationID string
	Role string
	CreatedAt time.Time
}