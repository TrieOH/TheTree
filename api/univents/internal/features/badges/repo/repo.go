package repo

import (
	"context"
	"encoding/json"
	"lib/database"
	"univents/errs"
	"univents/internal/sqlc"
	"univents/models"
	"univents/ports"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.BadgeTemplateRepo = (*Repo)(nil)
var _ ports.BadgeEmissionRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrBadgeTemplateNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

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
		EditionName:    src.EditionName,
		EndsAt:         src.EndsAt,
		EventName:      src.EventName,
		TicketTypeID:   src.TicketTypeID,
		TicketName:     src.TicketName,
	}
}

func mapBadgeEmissionViewFromEdition(src sqlc.ListBadgeEmissionViewsByEditionRow) models.BadgeEmissionView {
	return models.BadgeEmissionView{
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
		EditionName:    src.EditionName,
		EndsAt:         src.EndsAt,
		EventName:      src.EventName,
		TicketTypeID:   src.TicketTypeID,
		TicketName:     src.TicketName,
	}
}

func (repo *Repo) Create(ctx context.Context, template *models.BadgeTemplate) (*models.BadgeTemplate, error) {
	row, err := repo.db(ctx).CreateBadgeTemplate(ctx, sqlc.CreateBadgeTemplateParams{
		EditionID:    template.EditionID,
		TicketTypeID: template.TicketTypeID,
		Origin:       (*string)(template.Origin),
		Name:         template.Name,
		DesignData:   template.DesignData,
	})
	return database.One(row, err, repo.dbe, mapBadgeTemplate)
}

func (repo *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := repo.db(ctx).DeleteBadgeTemplate(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}

	return nil
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.BadgeTemplate, error) {
	row, err := repo.db(ctx).GetBadgeTemplateByID(ctx, id)
	return database.One(row, err, repo.dbe, mapBadgeTemplate)
}

func (repo *Repo) ListByEditionIDs(ctx context.Context, editionIDs []uuid.UUID) ([]models.BadgeTemplate, error) {
	rows, err := repo.db(ctx).ListBadgeTemplatesByEditionIDs(ctx, editionIDs)
	return database.Many(rows, err, repo.dbe, mapBadgeTemplate)
}

func (repo *Repo) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.BadgeTemplate, error) {
	rows, err := repo.db(ctx).ListBadgeTemplatesByEdition(ctx, editionID)
	return database.Many(rows, err, repo.dbe, mapBadgeTemplate)
}

func (repo *Repo) ListViewsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.BadgeEmissionView, error) {
	rows, err := repo.db(ctx).ListBadgeEmissionViewsByEdition(ctx, editionID)
	return database.Many(rows, err, repo.dbe, mapBadgeEmissionViewFromEdition)
}

func (repo *Repo) ListViewsByUser(ctx context.Context, userID uuid.UUID) ([]models.BadgeEmissionView, error) {
	rows, err := repo.db(ctx).ListBadgeEmissionViewsByUser(ctx, userID)
	return database.Many(rows, err, repo.dbe, mapBadgeEmissionView)
}

func (repo *Repo) MarkEmailSent(ctx context.Context, id uuid.UUID) error {
	err := repo.db(ctx).MarkBadgeEmissionEmailSent(ctx, id)
	return repo.dbe(err)
}

func (repo *Repo) Revoke(ctx context.Context, editionID, userID uuid.UUID, origin models.BadgeEmissionOrigin, reason string) error {
	err := repo.db(ctx).RevokeBadgeEmission(ctx, sqlc.RevokeBadgeEmissionParams{
		EditionID: editionID,
		UserID:    userID,
		Origin:    string(origin),
		Reason:    new(reason),
	})
	return repo.dbe(err)
}

func (repo *Repo) RevokeByRegistration(ctx context.Context, registrationID uuid.UUID, reason string) error {
	err := repo.db(ctx).RevokeBadgeEmissionByRegistration(ctx, sqlc.RevokeBadgeEmissionByRegistrationParams{
		RegistrationID: new(registrationID),
		Reason:         new(reason),
	})
	return repo.dbe(err)
}

func (repo *Repo) Update(ctx context.Context, id uuid.UUID, name string, designData json.RawMessage) (*models.BadgeTemplate, error) {
	row, err := repo.db(ctx).UpdateBadgeTemplate(ctx, sqlc.UpdateBadgeTemplateParams{
		ID:         id,
		Name:       name,
		DesignData: designData,
	})
	return database.One(row, err, repo.dbe, mapBadgeTemplate)
}

func (repo *Repo) Upsert(ctx context.Context, emission *models.BadgeEmission) (*models.BadgeEmission, error) {
	row, err := repo.db(ctx).UpsertBadgeEmission(ctx, sqlc.UpsertBadgeEmissionParams{
		EditionID:      emission.EditionID,
		UserID:         emission.UserID,
		Origin:         string(emission.Origin),
		RegistrationID: emission.RegistrationID,
	})
	return database.One(row, err, repo.dbe, mapBadgeEmission)
}
