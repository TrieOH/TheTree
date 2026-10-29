package ticket_types

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

var _ ports.TicketTypeRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrTicketTypeNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapTicketType(src sqlc.TicketType) models.TicketType {
	return models.TicketType{
		ID:          src.ID,
		EditionID:   src.EditionID,
		Name:        src.Name,
		Description: src.Description,
		AccessLevel: src.AccessLevel,
		PriceCents:  src.Price,
		MaxQuantity: src.MaxQuantity,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
		DeletedAt:   src.DeletedAt,
	}
}

func (repo *Repo) Create(ctx context.Context, toCreate *models.TicketType) (*models.TicketType, error) {
	result, err := repo.db(ctx).CreateTicketType(ctx, sqlc.CreateTicketTypeParams{
		EditionID:   toCreate.EditionID,
		Name:        toCreate.Name,
		Description: toCreate.Description,
		AccessLevel: toCreate.AccessLevel,
		Price:       toCreate.PriceCents,
		MaxQuantity: toCreate.MaxQuantity,
	})
	return database.One(result, err, repo.dbe, mapTicketType)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.TicketType, error) {
	result, err := repo.db(ctx).GetTicketTypeByID(ctx, id)
	return database.One(result, err, repo.dbe, mapTicketType)
}

// GetByIDForUpdate is the row-lock variant used inside the checkout tx
// (split 7): serializes concurrent checkouts on the same ticket type before
// availability is checked — no oversell on the last unit. The lock is held
// until the enclosing tx commits/rolls back.
func (repo *Repo) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*models.TicketType, error) {
	result, err := repo.db(ctx).GetTicketTypeByIDForUpdate(ctx, id)
	return database.One(result, err, repo.dbe, mapTicketType)
}

func (repo *Repo) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.TicketType, error) {
	results, err := repo.db(ctx).ListTicketTypesByEdition(ctx, editionID)
	return database.Many(results, err, repo.dbe, mapTicketType)
}

func (repo *Repo) Patch(ctx context.Context, id uuid.UUID, ticketType *models.TicketType) (*models.TicketType, error) {
	result, err := repo.db(ctx).PatchTicketType(ctx, sqlc.PatchTicketTypeParams{
		Name:        ticketType.Name,
		Description: ticketType.Description,
		AccessLevel: ticketType.AccessLevel,
		Price:       ticketType.PriceCents,
		MaxQuantity: ticketType.MaxQuantity,
		ID:          id,
	})
	return database.One(result, err, repo.dbe, mapTicketType)
}
