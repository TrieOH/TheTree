package models

import "github.com/google/uuid"

type Responder struct {
	ID     uuid.UUID  `json:"id"`
	UserID *uuid.UUID `json:"user_id"`
	Email  string     `json:"email"`
}
