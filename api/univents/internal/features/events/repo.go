package events

import (
	"context"
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

var _ ports.EventRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrEventNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapEventMember(src sqlc.EventMember) models.EventMember {
	return models.EventMember{
		ID:        src.ID,
		EventID:   src.EventID,
		UserID:    src.UserID,
		Role:      models.EventMemberRole(src.Role),
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
		DeletedAt: src.DeletedAt,
	}
}

func mapEvent(src sqlc.Event) models.Event {
	return models.Event{
		ID:                src.ID,
		OwnerID:           src.OwnerID,
		FullName:          src.FullName,
		Acronym:           src.Acronym,
		Slug:              src.Slug,
		Description:       src.Description,
		Style:             src.Style,
		Status:            models.EventStatus(src.Status),
		PayssageSellerID:  src.PayssageSellerID,
		PayssagePublicKey: src.PayssagePublicKey,
		LogoURL:           src.LogoUrl,
		BannerURL:         src.BannerUrl,
		ContactEmail:      src.ContactEmail,
		CreatedAt:         src.CreatedAt,
		UpdatedAt:         src.UpdatedAt,
		DeletedAt:         src.DeletedAt,
	}
}

func (repo *Repo) AddEventMember(ctx context.Context, eventID, userID uuid.UUID, role models.EventMemberRole) (*models.EventMember, error) {
	member, err := repo.db(ctx).AddEventMember(ctx, sqlc.AddEventMemberParams{
		EventID: eventID,
		UserID:  userID,
		Role:    string(role),
	})
	return database.One(member, err, repo.dbe, mapEventMember)
}

func (repo *Repo) Create(ctx context.Context, toCreate *models.Event) (*models.Event, error) {
	event, err := repo.db(ctx).CreateEvent(ctx, sqlc.CreateEventParams{
		OwnerID:      toCreate.OwnerID,
		FullName:     toCreate.FullName,
		Acronym:      toCreate.Acronym,
		Slug:         toCreate.Slug,
		Description:  toCreate.Description,
		Status:       string(toCreate.Status),
		ContactEmail: toCreate.ContactEmail,
	})
	return database.One(event, err, repo.dbe, mapEvent)
}

func (repo *Repo) Discontinue(ctx context.Context, id uuid.UUID) error {
	err := repo.db(ctx).DiscontinueEvent(ctx, id)
	return repo.dbe(err)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	event, err := repo.db(ctx).GetEventByID(ctx, id)
	return database.One(event, err, repo.dbe, mapEvent)
}

func (repo *Repo) GetBySlug(ctx context.Context, slug string) (*models.Event, error) {
	event, err := repo.db(ctx).GetEventBySlug(ctx, slug)
	return database.One(event, err, repo.dbe, mapEvent)
}

func (repo *Repo) GetMember(ctx context.Context, eventID, userID uuid.UUID) (*models.EventMember, error) {
	member, err := repo.db(ctx).GetEventMember(ctx, sqlc.GetEventMemberParams{
		EventID: eventID,
		UserID:  userID,
	})
	return database.One(member, err, repo.dbe, mapEventMember)
}

func (repo *Repo) GetRole(ctx context.Context, actorID, eventID uuid.UUID) (models.EventMemberRole, error) {
	member, err := repo.GetMember(ctx, eventID, actorID)
	if err != nil {
		return "", err
	}
	return member.Role, nil
}

func (repo *Repo) ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Event, error) {
	sqlcEvents, err := repo.db(ctx).ListJoinedEvents(ctx, userID)
	return database.Many(sqlcEvents, err, repo.dbe, mapEvent)
}

func (repo *Repo) ListEventMembers(ctx context.Context, eventID uuid.UUID) ([]models.EventMember, error) {
	members, err := repo.db(ctx).ListEventMembers(ctx, eventID)
	return database.Many(members, err, repo.dbe, mapEventMember)
}

func (repo *Repo) ListOwned(ctx context.Context, ownerID uuid.UUID) ([]models.Event, error) {
	sqlcEvents, err := repo.db(ctx).ListOwnedEvents(ctx, ownerID)
	return database.Many(sqlcEvents, err, repo.dbe, mapEvent)
}

func (repo *Repo) ListPublic(ctx context.Context) ([]models.Event, error) {
	sqlcEvents, err := repo.db(ctx).ListPublicEvents(ctx)
	return database.Many(sqlcEvents, err, repo.dbe, mapEvent)
}

func (repo *Repo) Patch(ctx context.Context, id uuid.UUID, event *models.Event) (*models.Event, error) {
	result, err := repo.db(ctx).PatchEvent(ctx, sqlc.PatchEventParams{
		FullName:     event.FullName,
		Acronym:      event.Acronym,
		Slug:         event.Slug,
		Description:  event.Description,
		LogoUrl:      event.LogoURL,
		BannerUrl:    event.BannerURL,
		ContactEmail: event.ContactEmail,
		ID:           id,
	})
	return database.One(result, err, repo.dbe, mapEvent)
}

func (repo *Repo) SetPaymentsConfig(ctx context.Context, id uuid.UUID, sellerID *uuid.UUID, publicKey *string) (*models.Event, error) {
	result, err := repo.db(ctx).SetEventPaymentsConfig(ctx, sqlc.SetEventPaymentsConfigParams{
		PayssageSellerID:  sellerID,
		PayssagePublicKey: publicKey,
		ID:                id,
	})
	return database.One(result, err, repo.dbe, mapEvent)
}

func (repo *Repo) ClearPaymentsConfig(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	result, err := repo.db(ctx).ClearEventPaymentsConfig(ctx, id)
	return database.One(result, err, repo.dbe, mapEvent)
}

func (repo *Repo) Publish(ctx context.Context, id uuid.UUID) error {
	err := repo.db(ctx).PublishEvent(ctx, id)
	return repo.dbe(err)
}

func (repo *Repo) RemoveEventMember(ctx context.Context, eventID, userID uuid.UUID) error {
	err := repo.db(ctx).RemoveEventMember(ctx, sqlc.RemoveEventMemberParams{
		EventID: eventID,
		UserID:  userID,
	})
	return repo.dbe(err)
}
