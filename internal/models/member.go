package models

import "github.com/google/uuid"

type Member struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
}
//this is for api respone, while membership.go is for db table