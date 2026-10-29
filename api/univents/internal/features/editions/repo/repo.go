package repo

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

var _ ports.EditionRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrEditionNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapEdition(src sqlc.Edition) models.Edition {
	return models.Edition{
		ID:                  src.ID,
		EventID:             src.EventID,
		Name:                src.EditionName,
		Slug:                src.Slug,
		Tagline:             src.Tagline,
		Description:         src.Description,
		IsDraft:             src.IsDraft,
		RegistrationOpensAt: src.RegistrationOpensAt,
		StartsAt:            src.StartsAt,
		EndsAt:              src.EndsAt,
		LocationName:        src.LocationName,
		LocationAddress:     src.LocationAddress,
		LogoURL:             src.LogoUrl,
		BannerURL:           src.BannerUrl,
		ContactEmail:        src.ContactEmail,
		CreatedBy:           src.CreatedBy,
		CreatedAt:           src.CreatedAt,
		UpdatedAt:           src.UpdatedAt,
		DeletedAt:           src.DeletedAt,
	}
}

func (repo *Repo) Create(ctx context.Context, toCreate *models.Edition) (*models.Edition, error) {
	edition, err := repo.db(ctx).CreateEdition(ctx, sqlc.CreateEditionParams{
		EventID:     toCreate.EventID,
		EditionName: toCreate.Name,
		Slug:        toCreate.Slug,
		StartsAt:    toCreate.StartsAt,
		EndsAt:      toCreate.EndsAt,
		CreatedBy:   toCreate.CreatedBy,
	})
	return database.One(edition, err, repo.dbe, mapEdition)
}

func (repo *Repo) GetActive(ctx context.Context, eventID uuid.UUID) (*models.Edition, error) {
	edition, err := repo.db(ctx).GetActiveEdition(ctx, eventID)
	return database.One(edition, err, repo.dbe, mapEdition)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Edition, error) {
	edition, err := repo.db(ctx).GetEditionByID(ctx, id)
	return database.One(edition, err, repo.dbe, mapEdition)
}

func (repo *Repo) GetBySlug(ctx context.Context, eventID uuid.UUID, slug string) (*models.Edition, error) {
	edition, err := repo.db(ctx).GetEditionBySlug(ctx, sqlc.GetEditionBySlugParams{
		Slug:    slug,
		EventID: eventID,
	})
	return database.One(edition, err, repo.dbe, mapEdition)
}

func (repo *Repo) GetPast(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	editions, err := repo.db(ctx).GetPastEditions(ctx, eventID)
	return database.Many(editions, err, repo.dbe, mapEdition)
}

func (repo *Repo) GetUpcoming(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	editions, err := repo.db(ctx).GetUpcomingEditions(ctx, eventID)
	return database.Many(editions, err, repo.dbe, mapEdition)
}

func (repo *Repo) ListDraft(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	editions, err := repo.db(ctx).ListDraftEditions(ctx, eventID)
	return database.Many(editions, err, repo.dbe, mapEdition)
}

func (repo *Repo) ListPublic(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	editions, err := repo.db(ctx).ListPublicEditions(ctx, eventID)
	return database.Many(editions, err, repo.dbe, mapEdition)
}

func (repo *Repo) Patch(ctx context.Context, id uuid.UUID, edition *models.Edition) (*models.Edition, error) {
	result, err := repo.db(ctx).PatchEdition(ctx, sqlc.PatchEditionParams{
		EditionName:         edition.Name,
		Slug:                edition.Slug,
		Tagline:             edition.Tagline,
		Description:         edition.Description,
		RegistrationOpensAt: edition.RegistrationOpensAt,
		StartsAt:            edition.StartsAt,
		EndsAt:              edition.EndsAt,
		LocationName:        edition.LocationName,
		LocationAddress:     edition.LocationAddress,
		LogoUrl:             edition.LogoURL,
		BannerUrl:           edition.BannerURL,
		ContactEmail:        edition.ContactEmail,
		ID:                  id,
	})
	return database.One(result, err, repo.dbe, mapEdition)
}

func (repo *Repo) Publish(ctx context.Context, id uuid.UUID) error {
	err := repo.db(ctx).PublishEdition(ctx, id)
	return repo.dbe(err)
}
