package infrastructure

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/plataform/database"
	"univents/internal/plataform/database/sqlc"
	"univents/internal/plataform/telemetry"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type editionsRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ domain.EditionsRepository = (*editionsRepo)(nil)

func NewEditionRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) domain.EditionsRepository {
	return &editionsRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *editionsRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapEditionFromDB(src *sqlc.Edition) *domain.Edition {
	return &domain.Edition{
		ID:                            src.ID,
		GoauthScopeID:                 src.GoauthScopeID,
		EventID:                       src.EventID,
		Type:                          domain.EditionType(src.Type),
		EditionName:                   src.EditionName,
		Tagline:                       src.Tagline,
		Description:                   src.Description,
		Status:                        domain.EditionStatus(src.Status),
		MonetaryType:                  domain.EditionMonetaryType(src.MonetaryType),
		RegistrationOpensAt:           src.RegistrationOpensAt,
		RegistrationClosesAt:          src.RegistrationClosesAt,
		StartsAt:                      src.StartsAt,
		EndsAt:                        src.EndsAt,
		Timezone:                      src.Timezone,
		LocationName:                  src.LocationName,
		LocationAddress:               src.LocationAddress,
		LogoUrl:                       src.LogoUrl,
		BannerUrl:                     src.BannerUrl,
		ContactEmail:                  src.ContactEmail,
		ContactPhone:                  src.ContactPhone,
		OrganizerName:                 src.OrganizerName,
		TriePaymentsCredentialID:      src.TriePaymentsCredentialID,
		TriePaymentsProvider:          src.TriePaymentsProvider,
		TriePaymentsProviderPublicKey: src.TriePaymentsProviderPublicKey,
		CreatedBy:                     src.CreatedBy,
		CreatedAt:                     src.CreatedAt,
		UpdatedAt:                     src.UpdatedAt,
		DeletedAt:                     src.DeletedAt,
	}
}

func EditionStatusPtrToString(p *domain.EditionStatus) string {
	if p == nil {
		return ""
	}
	return string(*p)
}

func (repo *editionsRepo) Create(ctx context.Context, toCreate *domain.Edition) (*domain.Edition, error) {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.Create")
	defer span.End()

	edition, err := repo.queries(ctx).CreateEdition(ctx, sqlc.CreateEditionParams{
		ID:                   toCreate.ID,
		EventID:              toCreate.EventID,
		GoauthScopeID:        toCreate.GoauthScopeID,
		Type:                 sqlc.EditionType(toCreate.Type),
		EditionName:          toCreate.EditionName,
		Tagline:              toCreate.Tagline,
		Description:          toCreate.Description,
		Status:               sqlc.EditionStatus(domain.EditionStatusDraft),
		RegistrationOpensAt:  toCreate.RegistrationOpensAt,
		RegistrationClosesAt: toCreate.RegistrationClosesAt,
		MonetaryType:         sqlc.EditionMonetaryType(toCreate.MonetaryType),
		StartsAt:             toCreate.StartsAt,
		EndsAt:               toCreate.EndsAt,
		Timezone:             toCreate.Timezone,
		LocationName:         toCreate.LocationName,
		LocationAddress:      toCreate.LocationAddress,
		LogoUrl:              toCreate.LogoUrl,
		BannerUrl:            toCreate.BannerUrl,
		ContactEmail:         toCreate.ContactEmail,
		ContactPhone:         toCreate.ContactPhone,
		OrganizerName:        toCreate.OrganizerName,
		CreatedBy:            toCreate.CreatedBy,
	})
	if err != nil {
		return nil, errx.FromDB(err, "edition")
	}

	return mapEditionFromDB(&edition), nil
}

func (repo *editionsRepo) GetByID(ctx context.Context, editionID uuid.UUID) (*domain.Edition, error) {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.GetByID")
	defer span.End()

	sqlcEdition, err := repo.queries(ctx).GetEditionByID(ctx, editionID)
	if err != nil {
		return nil, errx.FromDB(err, "edition")
	}

	return mapEditionFromDB(&sqlcEdition), nil
}

func (repo *editionsRepo) List(ctx context.Context, eventID uuid.UUID) ([]domain.Edition, error) {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.List")
	defer span.End()

	sqlcEditions, err := repo.queries(ctx).ListEditions(ctx, eventID)
	if err != nil {
		return nil, errx.FromDB(err, "edition")
	}

	outEditions := make([]domain.Edition, 0, len(sqlcEditions))
	for _, sqlcEdition := range sqlcEditions {
		outEditions = append(outEditions, *mapEditionFromDB(&sqlcEdition))
	}
	return outEditions, nil
}

func (repo *editionsRepo) ListAdmin(ctx context.Context, eventID uuid.UUID) ([]domain.Edition, error) {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.ListAdmin")
	defer span.End()

	sqlcEditions, err := repo.queries(ctx).ListEditionsAdmin(ctx, eventID)
	if err != nil {
		return nil, errx.FromDB(err, "edition")
	}

	outEditions := make([]domain.Edition, 0, len(sqlcEditions))
	for _, sqlcEdition := range sqlcEditions {
		outEditions = append(outEditions, *mapEditionFromDB(&sqlcEdition))
	}
	return outEditions, nil
}

func (repo *editionsRepo) Announce(ctx context.Context, editionID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.Announce")
	defer span.End()

	err := repo.queries(ctx).AnnounceEdition(ctx, editionID)
	if err != nil {
		return errx.FromDB(err, "edition")
	}

	return nil
}

func (repo *editionsRepo) Open(ctx context.Context, editionID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.Open")
	defer span.End()

	err := repo.queries(ctx).OpenEditionRegistrations(ctx, editionID)
	if err != nil {
		return errx.FromDB(err, "edition")
	}

	return nil
}

func (repo *editionsRepo) Start(ctx context.Context, editionID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.Start")
	defer span.End()

	err := repo.queries(ctx).StartEdition(ctx, editionID)
	if err != nil {
		return errx.FromDB(err, "edition")
	}

	return nil
}

func (repo *editionsRepo) Finish(ctx context.Context, editionID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.Finish")
	defer span.End()

	err := repo.queries(ctx).FinishEdition(ctx, editionID)
	if err != nil {
		return errx.FromDB(err, "edition")
	}

	return nil
}

func (repo *editionsRepo) ConnectPaymentsAccount(ctx context.Context, editionID, triePaymentsCredentialID uuid.UUID, triePaymentsProvider, publicKey string) error {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.ConnectPaymentsAccount")
	defer span.End()

	telemetry.Log().Info("ConnectPaymentsAccount DATA", zap.String("edition_id", editionID.String()), zap.String("trie_payments_credential_id", triePaymentsCredentialID.String()), zap.String("trie_payments_provider", triePaymentsProvider))

	err := repo.queries(ctx).ConnectEditionPaymentAccount(ctx, sqlc.ConnectEditionPaymentAccountParams{
		TriePaymentsCredentialID:      &triePaymentsCredentialID,
		TriePaymentsProvider:          &triePaymentsProvider,
		TriePaymentsProviderPublicKey: &publicKey,
		ID:                            editionID,
	})
	if err != nil {
		return errx.FromDB(err, "edition")
	}

	return nil
}

func (repo *editionsRepo) DisconnectPaymentsAccount(ctx context.Context, editionID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.DisconnectPaymentsAccount")
	defer span.End()

	err := repo.queries(ctx).DisconnectEditionPaymentAccount(ctx, editionID)
	if err != nil {
		return errx.FromDB(err, "edition")
	}

	return nil
}
