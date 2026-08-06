package badges

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

var _ ports.BadgeTemplateRepo = (*Repo)(nil)
var _ ports.BadgeEmissionRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("badge_template"),
	}
}

func mapBadgeTemplate(src sqlc.BadgeTemplate) models.BadgeTemplate {
	return models.BadgeTemplate{
		ID:           src.ID,
		EditionID:    src.EditionID,
		TicketTypeID: src.TicketTypeID,
		Origin:       (*models.BadgeTemplateOrigin)(src.Origin),
		Name:         src.Name,
		DesignData:   src.DesignData,
		CreatedAt:    src.CreatedAt,
		UpdatedAt:    src.UpdatedAt,
		DeletedAt:    src.DeletedAt,
	}
}

func mapBadgeEmission(src sqlc.BadgeEmission) models.BadgeEmission {
	return models.BadgeEmission{
		ID:             src.ID,
		EditionID:      src.EditionID,
		UserID:         src.UserID,
		Origin:         models.BadgeEmissionOrigin(src.Origin),
		RegistrationID: src.RegistrationID,
		Status:         models.BadgeEmissionStatus(src.Status),
		StatusReason:   src.StatusReason,
		EmailSentAt:    src.EmailSentAt,
		EmittedAt:      src.EmittedAt,
		UpdatedAt:      src.UpdatedAt,
	}
}

func mapBadgeEmissionView(src sqlc.ListBadgeEmissionViewsByUserRow) models.BadgeEmissionView {
	return models.BadgeEmissionView{
		BadgeEmission: models.BadgeEmission{
			ID:             src.ID,
			EditionID:      src.EditionID,
			UserID:         src.UserID,
			Origin:         models.BadgeEmissionOrigin(src.Origin),
			RegistrationID: src.RegistrationID,
			Status:         models.BadgeEmissionStatus(src.Status),
			StatusReason:   src.StatusReason,
			EmailSentAt:    src.EmailSentAt,
			EmittedAt:      src.EmittedAt,
			UpdatedAt:      src.UpdatedAt,
		},
		EditionName:  src.EditionName,
		EndsAt:       src.EndsAt,
		EventName:    src.EventName,
		TicketTypeID: src.TicketTypeID,
		TicketName:   src.TicketName,
	}
}

func mapBadgeEmissionViewFromEdition(src sqlc.ListBadgeEmissionViewsByEditionRow) models.BadgeEmissionView {
	return models.BadgeEmissionView{
		BadgeEmission: models.BadgeEmission{
			ID:             src.ID,
			EditionID:      src.EditionID,
			UserID:         src.UserID,
			Origin:         models.BadgeEmissionOrigin(src.Origin),
			RegistrationID: src.RegistrationID,
			Status:         models.BadgeEmissionStatus(src.Status),
			StatusReason:   src.StatusReason,
			EmailSentAt:    src.EmailSentAt,
			EmittedAt:      src.EmittedAt,
			UpdatedAt:      src.UpdatedAt,
		},
		EditionName:  src.EditionName,
		EndsAt:       src.EndsAt,
		EventName:    src.EventName,
		TicketTypeID: src.TicketTypeID,
		TicketName:   src.TicketName,
	}
}
