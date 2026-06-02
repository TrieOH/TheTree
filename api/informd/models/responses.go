package models

import (
	"time"

	"github.com/google/uuid"
)

type Response struct {
	ID          uuid.UUID  `json:"id"`
	FormID      uuid.UUID  `json:"form_id"`
	InviteID    *uuid.UUID `json:"invite_id"`
	ResponderID *uuid.UUID `json:"responder_id"`
	Email       *string    `json:"email"`
	StartedAt   time.Time  `json:"started_at"`
	FinishedAt  *time.Time `json:"finished_at"`
}
