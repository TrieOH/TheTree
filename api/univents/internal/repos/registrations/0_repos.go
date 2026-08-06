package registrations

import (
	"univents/internal/sqlc"
	"univents/models"
	"univents/ports"

	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.RegistrationRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("registration"),
	}
}

func mapRegistration(src sqlc.Registration) models.Registration {
	return models.Registration{
		ID:               src.ID,
		EditionID:        src.EditionID,
		TicketTypeID:     src.TicketTypeID,
		PurchaserID:      src.PurchaserID,
		AttendeeUserID:   src.AttendeeUserID,
		AttendeeEmail:    src.AttendeeEmail,
		AttendeeName:     src.AttendeeName,
		Status:           models.RegistrationStatus(src.Status),
		StatusReason:     src.StatusReason,
		PayssageIntentID: src.PayssageIntentID,
		CreatedAt:        src.CreatedAt,
		UpdatedAt:        src.UpdatedAt,
		DeletedAt:        src.DeletedAt,
	}
}
